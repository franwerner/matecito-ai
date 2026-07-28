package store

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

// sha256Prefix is the algorithm tag every content hash in this store
// carries inline, so a value produced by a future algorithm never sits
// ambiguous next to one from this one (R7).
const sha256Prefix = "sha256:"

// contentHash is the single place the store computes a content hash: SHA-256
// of b exactly as given, hex-encoded lowercase, tagged with its algorithm.
// content_objects reuses it for content-addressing (R7); event idempotency
// reuses it for the payload hash (R6, Batch 3 decision) — both are the same
// operation, identify bytes by their hash, so both go through here instead
// of each hashing on its own.
func contentHash(b []byte) string {
	sum := sha256.Sum256(b)
	return sha256Prefix + hex.EncodeToString(sum[:])
}

// validateContentHash is the single place the store checks a hash string is
// shaped like one it produced: the right algorithm tag, followed by exactly
// a SHA-256 digest's worth of lowercase hex.
func validateContentHash(hash string) error {
	rest, ok := strings.CutPrefix(hash, sha256Prefix)
	if !ok || len(rest) != hex.EncodedLen(sha256.Size) {
		return fmt.Errorf("%w: %q", ErrInvalidContentHash, hash)
	}
	if _, err := hex.DecodeString(rest); err != nil {
		return fmt.Errorf("%w: %q", ErrInvalidContentHash, hash)
	}
	return nil
}
