package git

import (
	"os/exec"
	"path/filepath"
	"syntaxia/internal/executil"
	"sort"
	"strings"
	"time"
)

// ContextBuilder builds context from git history
type ContextBuilder struct {
	projectRoot string
}

// NewContextBuilder creates a new git context builder
func NewContextBuilder(projectRoot string) *ContextBuilder {
	return &ContextBuilder{projectRoot: projectRoot}
}

// RecentChange represents a recently changed file
type RecentChange struct {
	FilePath    string
	ChangeCount int
	LastChanged time.Time
	Authors     []string
}

// GetRecentChanges returns files changed recently, sorted by relevance
func (b *ContextBuilder) GetRecentChanges(since string, pathFilter string) ([]RecentChange, error) {
	if since == "" {
		since = "1 week ago"
	}

	output, err := b.runGitLog(since, pathFilter)
	if err != nil {
		return nil, err
	}

	changes := b.parseGitLogOutput(string(output))
	return b.sortRecentChanges(changes), nil
}

// runGitLog executes git log command
func (b *ContextBuilder) runGitLog(since, pathFilter string) ([]byte, error) {
	cmdArgs := []string{"log", "--since=" + since, "--name-only", "--format=%H|%an|%at"}
	if pathFilter != "" {
		cmdArgs = append(cmdArgs, "--", pathFilter)
	}
	cmd := exec.Command("git", cmdArgs...)
	executil.HideWindow(cmd)
	cmd.Dir = b.projectRoot
	return cmd.Output()
}

// parseGitLogOutput parses git log output into changes map
func (b *ContextBuilder) parseGitLogOutput(output string) map[string]*RecentChange {
	changes := make(map[string]*RecentChange)
	var currentAuthor string
	var currentTime time.Time

	for _, line := range strings.Split(output, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if strings.Contains(line, "|") {
			currentAuthor, currentTime = b.parseCommitLine(line)
			continue
		}

		b.updateFileChange(changes, line, currentAuthor, currentTime)
	}
	return changes
}

// parseCommitLine parses a commit metadata line
func (b *ContextBuilder) parseCommitLine(line string) (string, time.Time) {
	parts := strings.Split(line, "|")
	if len(parts) < 3 {
		return "", time.Time{}
	}
	var ts int64
	_, _ = parseUnixTime(parts[2], &ts)
	return parts[1], time.Unix(ts, 0)
}

// updateFileChange updates or creates a file change entry
func (b *ContextBuilder) updateFileChange(changes map[string]*RecentChange, filePath, author string, changeTime time.Time) {
	if change, exists := changes[filePath]; exists {
		change.ChangeCount++
		if changeTime.After(change.LastChanged) {
			change.LastChanged = changeTime
		}
		if !containsString(change.Authors, author) {
			change.Authors = append(change.Authors, author)
		}
	} else {
		changes[filePath] = &RecentChange{
			FilePath: filePath, ChangeCount: 1, LastChanged: changeTime, Authors: []string{author},
		}
	}
}

// sortRecentChanges converts map to sorted slice
func (b *ContextBuilder) sortRecentChanges(changes map[string]*RecentChange) []RecentChange {
	result := make([]RecentChange, 0, len(changes))
	for _, c := range changes {
		result = append(result, *c)
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].ChangeCount != result[j].ChangeCount {
			return result[i].ChangeCount > result[j].ChangeCount
		}
		return result[i].LastChanged.After(result[j].LastChanged)
	})
	return result
}

// GetRelatedByAuthor returns files frequently changed by the same author
func (b *ContextBuilder) GetRelatedByAuthor(filePath string, limit int) ([]string, error) {
	if limit <= 0 {
		limit = 10
	}

	// Get authors who changed this file
	cmd := exec.Command("git", "log", "--format=%an", "--", filePath)
	executil.HideWindow(cmd)
	cmd.Dir = b.projectRoot
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	authors := make(map[string]int)
	for _, line := range strings.Split(string(output), "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			authors[line]++
		}
	}

	// Find top author
	var topAuthor string
	var maxCount int
	for author, count := range authors {
		if count > maxCount {
			topAuthor = author
			maxCount = count
		}
	}

	if topAuthor == "" {
		return nil, nil
	}

	// Get files changed by this author
	cmd = exec.Command("git", "log", "--author="+topAuthor, "--name-only", "--format=", "-n", "100") //nolint:gosec // Git command with validated input
	executil.HideWindow(cmd)
	cmd.Dir = b.projectRoot
	output, err = cmd.Output()
	if err != nil {
		return nil, err
	}

	fileCounts := make(map[string]int)
	for _, line := range strings.Split(string(output), "\n") {
		line = strings.TrimSpace(line)
		if line != "" && line != filePath {
			fileCounts[line]++
		}
	}

	// Sort by count
	type fileCount struct {
		path  string
		count int
	}
	sorted := make([]fileCount, 0, len(fileCounts))
	for path, count := range fileCounts {
		sorted = append(sorted, fileCount{path, count})
	}
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].count > sorted[j].count
	})

	result := make([]string, 0, limit)
	for i, fc := range sorted {
		if i >= limit {
			break
		}
		result = append(result, fc.path)
	}

	return result, nil
}

// GetCoChangedFiles returns files that are often changed together with the given file
func (b *ContextBuilder) GetCoChangedFiles(filePath string, limit int) ([]string, error) {
	if limit <= 0 {
		limit = 10
	}

	// Get commits that changed this file
	cmd := exec.Command("git", "log", "--format=%H", "-n", "50", "--", filePath)
	executil.HideWindow(cmd)
	cmd.Dir = b.projectRoot
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	commits := strings.Split(strings.TrimSpace(string(output)), "\n")
	fileCounts := make(map[string]int)

	for _, commit := range commits {
		if commit == "" {
			continue
		}
		// Get files changed in this commit
		cmd = exec.Command("git", "show", "--name-only", "--format=", commit)
		executil.HideWindow(cmd)
		cmd.Dir = b.projectRoot
		output, err = cmd.Output()
		if err != nil {
			continue
		}

		for _, file := range strings.Split(string(output), "\n") {
			file = strings.TrimSpace(file)
			if file != "" && file != filePath {
				fileCounts[file]++
			}
		}
	}

	// Sort by count
	type fileCount struct {
		path  string
		count int
	}
	sorted := make([]fileCount, 0, len(fileCounts))
	for path, count := range fileCounts {
		sorted = append(sorted, fileCount{path, count})
	}
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].count > sorted[j].count
	})

	result := make([]string, 0, limit)
	for i, fc := range sorted {
		if i >= limit {
			break
		}
		result = append(result, fc.path)
	}

	return result, nil
}

// SuggestContextFiles suggests files to include in context based on git history
func (b *ContextBuilder) SuggestContextFiles(taskDescription string, currentFiles []string, limit int) ([]string, error) {
	if limit <= 0 {
		limit = 10
	}

	suggestions := make(map[string]int)
	b.addCoChangedSuggestions(currentFiles, suggestions)
	b.addRecentChangeSuggestions(currentFiles, suggestions)
	b.addKeywordSuggestions(taskDescription, suggestions)

	for _, f := range currentFiles {
		delete(suggestions, f)
	}

	return b.topSuggestions(suggestions, limit), nil
}

// addCoChangedSuggestions adds co-changed files to suggestions
func (b *ContextBuilder) addCoChangedSuggestions(currentFiles []string, suggestions map[string]int) {
	for _, file := range currentFiles {
		if coChanged, err := b.GetCoChangedFiles(file, 5); err == nil {
			for _, f := range coChanged {
				suggestions[f] += 3
			}
		}
	}
}

// addRecentChangeSuggestions adds recently changed files in same directories
func (b *ContextBuilder) addRecentChangeSuggestions(currentFiles []string, suggestions map[string]int) {
	for _, file := range currentFiles {
		if recent, err := b.GetRecentChanges("2 weeks ago", filepath.Dir(file)); err == nil {
			for _, r := range recent {
				if r.FilePath != file {
					suggestions[r.FilePath] += r.ChangeCount
				}
			}
		}
	}
}

// addKeywordSuggestions searches commits by task keywords
func (b *ContextBuilder) addKeywordSuggestions(taskDescription string, suggestions map[string]int) {
	for _, kw := range extractKeywords(taskDescription) {
		cmd := exec.Command("git", "log", "--grep="+kw, "-i", "--name-only", "--format=", "-n", "20") //nolint:gosec // Git command with validated input
		executil.HideWindow(cmd)
		cmd.Dir = b.projectRoot
		if output, err := cmd.Output(); err == nil {
			for _, file := range strings.Split(string(output), "\n") {
				if file = strings.TrimSpace(file); file != "" {
					suggestions[file] += 2
				}
			}
		}
	}
}

// topSuggestions returns top N suggestions sorted by score
func (b *ContextBuilder) topSuggestions(suggestions map[string]int, limit int) []string {
	type suggestion struct {
		path  string
		score int
	}
	sorted := make([]suggestion, 0, len(suggestions))
	for path, score := range suggestions {
		sorted = append(sorted, suggestion{path, score})
	}
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].score > sorted[j].score
	})

	result := make([]string, 0, limit)
	for i, s := range sorted {
		if i >= limit {
			break
		}
		result = append(result, s.path)
	}
	return result
}

func containsString(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

func parseUnixTime(s string, ts *int64) (int, error) {
	*ts = 0 // Reset to avoid accumulation from previous calls
	var n int
	for _, c := range s {
		if c >= '0' && c <= '9' {
			*ts = *ts*10 + int64(c-'0')
			n++
		}
	}
	return n, nil
}

func extractKeywords(text string) []string {
	// Simple keyword extraction - split by spaces and filter
	words := strings.Fields(strings.ToLower(text))
	var keywords []string
	stopWords := map[string]bool{
		"the": true, "a": true, "an": true, "and": true, "or": true,
		"is": true, "are": true, "was": true, "were": true, "be": true,
		"to": true, "of": true, "in": true, "for": true, "on": true,
		"with": true, "as": true, "at": true, "by": true, "from": true,
		"this": true, "that": true, "it": true, "i": true, "we": true,
		"add": true, "fix": true, "update": true, "change": true,
	}

	for _, w := range words {
		if len(w) > 3 && !stopWords[w] {
			keywords = append(keywords, w)
		}
	}

	if len(keywords) > 5 {
		keywords = keywords[:5]
	}
	return keywords
}
