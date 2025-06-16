package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

func GenerateSourceHash(sourceCode string) string {
	hash := sha256.Sum256([]byte(sourceCode))
	return hex.EncodeToString(hash[:])
}
