package contextbuilder

import (
	"encoding/json"
	"regexp"
	"sort"
	"strings"

	"syntaxia/domain"
)

type BuildOptions struct {
	StripComments   bool
	IncludeManifest bool
}

// entry — одна запись контекста
type entry struct {
	Path    string `json:"path"`
	Content string `json:"content"`
}

// --- File: path --- парсер
var headerRe = regexp.MustCompile(`(?m)^--- File:\s+(.*?)\s+---\s*`)

func parseContext(ctx string) []entry {
	res := []entry{}
	idxs := headerRe.FindAllStringIndex(ctx, -1)
	if len(idxs) == 0 {
		return res
	}
	for i := 0; i < len(idxs); i++ {
		start := idxs[i][0]
		end := len(ctx)
		if i+1 < len(idxs) {
			end = idxs[i+1][0]
		}
		block := ctx[start:end]
		m := headerRe.FindStringSubmatch(block)
		if len(m) < 2 {
			continue
		}
		path := strings.TrimSpace(m[1])
		content := strings.TrimSpace(block[len(m[0]):])
		res = append(res, entry{Path: path, Content: content})
	}

	// Обеспечиваем детерминированный порядок
	sort.Slice(res, func(i, j int) bool {
		return res[i].Path < res[j].Path
	})

	return res
}

// Грубая очистка комментариев
func stripComments(text string) string {
	lineRE := regexp.MustCompile(`(?m)^\s*(//|#).*$`)
	blockRE := regexp.MustCompile(`/\*[\s\S]*?\*/`)
	out := lineRE.ReplaceAllString(text, "")
	out = blockRE.ReplaceAllString(out, "")
	return out
}

// Построение ASCII-дерева из путей
func buildTree(paths []string) string {
	type node struct {
		name     string
		children map[string]*node
	}
	root := &node{name: ".", children: map[string]*node{}}

	// Сортируем входные пути для детерминизма
	sortedPaths := make([]string, len(paths))
	copy(sortedPaths, paths)
	sort.Strings(sortedPaths)

	for _, p := range sortedPaths {
		// нормализуем слеши
		pp := strings.ReplaceAll(p, "\\", "/")
		parts := strings.Split(pp, "/")
		cur := root
		for _, part := range parts {
			if part == "" || part == "." {
				continue
			}
			if cur.children[part] == nil {
				cur.children[part] = &node{name: part, children: map[string]*node{}}
			}
			cur = cur.children[part]
		}
	}

	var b strings.Builder
	var walk func(n *node, prefix string, isLast bool)
	walk = func(n *node, prefix string, isLast bool) {
		if n != root {
			if isLast {
				b.WriteString(prefix + "└─ " + n.name + "\n")
				prefix += "   "
			} else {
				b.WriteString(prefix + "├─ " + n.name + "\n")
				prefix += "│  "
			}
		}
		// детерминированный порядок
		keys := make([]string, 0, len(n.children))
		for k := range n.children {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for i, k := range keys {
			walk(n.children[k], prefix, i == len(keys)-1)
		}
	}
	walk(root, "", true)
	return b.String()
}

// buildPlainFormat builds plain text format output
func buildPlainFormat(entries []entry, opts BuildOptions) string {
	var b strings.Builder
	for _, e := range entries {
		content := e.Content
		if opts.StripComments {
			content = stripComments(content)
		}
		b.WriteString("--- File: " + e.Path + " ---\n")
		b.WriteString(content)
		b.WriteString("\n\n")
	}
	return strings.TrimSpace(b.String())
}

// buildManifestFormat builds manifest format output with tree
func buildManifestFormat(entries []entry, opts BuildOptions) string {
	paths := make([]string, 0, len(entries))
	for _, e := range entries {
		paths = append(paths, e.Path)
	}
	var b strings.Builder
	b.WriteString("Manifest:\n")
	b.WriteString(buildTree(paths))
	b.WriteString("\n")
	for _, e := range entries {
		content := e.Content
		if opts.StripComments {
			content = stripComments(content)
		}
		b.WriteString("--- File: " + e.Path + " ---\n")
		b.WriteString(content)
		b.WriteString("\n\n")
	}
	return strings.TrimSpace(b.String())
}

// buildJSONFormat builds JSON format output
func buildJSONFormat(entries []entry, opts BuildOptions) (string, error) {
	j := make([]entry, 0, len(entries))
	for _, e := range entries {
		c := e.Content
		if opts.StripComments {
			c = stripComments(c)
		}
		j = append(j, entry{Path: e.Path, Content: c})
	}
	raw, err := json.MarshalIndent(j, "", "  ")
	if err != nil {
		return "", err
	}
	return string(raw), nil
}

// getLanguageFromPath returns the language identifier for syntax highlighting based on file extension
func getLanguageFromPath(path string) string {
	ext := strings.ToLower(path)
	if idx := strings.LastIndex(ext, "."); idx >= 0 {
		ext = ext[idx+1:]
	} else {
		return ""
	}

	langMap := map[string]string{
		"go":         "go",
		"js":         "javascript",
		"ts":         "typescript",
		"jsx":        "jsx",
		"tsx":        "tsx",
		"py":         "python",
		"rb":         "ruby",
		"java":       "java",
		"kt":         "kotlin",
		"cs":         "csharp",
		"cpp":        "cpp",
		"c":          "c",
		"h":          "c",
		"hpp":        "cpp",
		"rs":         "rust",
		"swift":      "swift",
		"php":        "php",
		"vue":        "vue",
		"svelte":     "svelte",
		"html":       "html",
		"css":        "css",
		"scss":       "scss",
		"sass":       "sass",
		"less":       "less",
		"json":       "json",
		"yaml":       "yaml",
		"yml":        "yaml",
		"xml":        "xml",
		"sql":        "sql",
		"sh":         "bash",
		"bash":       "bash",
		"zsh":        "bash",
		"ps1":        "powershell",
		"md":         "markdown",
		"dart":       "dart",
		"lua":        "lua",
		"r":          "r",
		"scala":      "scala",
		"groovy":     "groovy",
		"gradle":     "groovy",
		"tf":         "hcl",
		"hcl":        "hcl",
		"dockerfile": "dockerfile",
		"makefile":   "makefile",
	}

	if lang, ok := langMap[ext]; ok {
		return lang
	}
	return ext
}

// buildMarkdownFormat builds markdown format output with code blocks
func buildMarkdownFormat(entries []entry, opts BuildOptions) string {
	var b strings.Builder
	for _, e := range entries {
		content := e.Content
		if opts.StripComments {
			content = stripComments(content)
		}
		lang := getLanguageFromPath(e.Path)
		b.WriteString("## File: " + e.Path + "\n\n")
		b.WriteString("```" + lang + "\n")
		b.WriteString(content)
		if !strings.HasSuffix(content, "\n") {
			b.WriteString("\n")
		}
		b.WriteString("```\n\n")
	}
	return strings.TrimSpace(b.String())
}

// buildXMLFormat builds XML format output
func buildXMLFormat(entries []entry, opts BuildOptions) string {
	var b strings.Builder
	b.WriteString("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n")
	b.WriteString("<files>\n")
	for _, e := range entries {
		content := e.Content
		if opts.StripComments {
			content = stripComments(content)
		}
		// Escape XML special characters
		content = strings.ReplaceAll(content, "&", "&amp;")
		content = strings.ReplaceAll(content, "<", "&lt;")
		content = strings.ReplaceAll(content, ">", "&gt;")
		b.WriteString("  <file path=\"" + e.Path + "\">\n")
		b.WriteString("    <content><![CDATA[" + content + "]]></content>\n")
		b.WriteString("  </file>\n")
	}
	b.WriteString("</files>")
	return b.String()
}

// BuildFromContext — собирает строку по формату: "plain" | "manifest" | "json" | "markdown" | "xml"
func BuildFromContext(format string, ctx string, opts BuildOptions) (string, error) {
	entries := parseContext(ctx)

	switch strings.ToLower(format) {
	case "plain":
		return buildPlainFormat(entries, opts), nil
	case "manifest", "manifest+text":
		return buildManifestFormat(entries, opts), nil
	case "json":
		return buildJSONFormat(entries, opts)
	case "markdown":
		return buildMarkdownFormat(entries, opts), nil
	case "xml":
		return buildXMLFormat(entries, opts), nil
	default:
		return BuildFromContext("manifest", ctx, opts)
	}
}

// ContextFormatterImpl implements domain.ContextFormatter
type ContextFormatterImpl struct{}

// NewContextFormatter creates a new context formatter
func NewContextFormatter() *ContextFormatterImpl {
	return &ContextFormatterImpl{}
}

// Format formats context string according to specified format
func (f *ContextFormatterImpl) Format(format string, contextContent string, opts domain.ContextFormatOptions) (string, error) {
	return BuildFromContext(format, contextContent, BuildOptions{
		StripComments:   opts.StripComments,
		IncludeManifest: opts.IncludeManifest,
	})
}
