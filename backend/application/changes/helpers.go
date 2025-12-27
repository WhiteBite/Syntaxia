package changes

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sync/atomic"
	"time"
)

// idCounter is used for generating unique IDs
var idCounter uint64

// computeHash computes SHA256 hash of content
func computeHash(content []byte) string {
	hash := sha256.Sum256(content)
	return hex.EncodeToString(hash[:])
}

// generateID generates a unique ID with prefix
func generateID(prefix string) string {
	counter := atomic.AddUint64(&idCounter, 1)
	return fmt.Sprintf("%s-%d-%d", prefix, time.Now().UnixNano(), counter)
}
