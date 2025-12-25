package application

// SimpleTokenCounter is a basic token counter implementation.
// It approximates tokens by dividing the number of characters by 4.
func SimpleTokenCounter(text string) int {
	return len(text) / 4
}
