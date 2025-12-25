package domain

// DiffSplitter describes a service for splitting git diffs into chunks.
type DiffSplitter interface {
	// Split takes a string containing a git diff and splits it into multiple strings,
	// each not exceeding the approximate line limit.
	Split(diffText string, approxLineLimit int) ([]string, error)
}
