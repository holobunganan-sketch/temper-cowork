package store

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
)

// newID 生成带前缀的随机 ID(如 prj-<16 hex>)。
func newID(prefix string) string {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		panic(fmt.Sprintf("temper store: entropy unavailable: %v", err))
	}
	return prefix + "-" + hex.EncodeToString(b)
}

// isUniqueViolation 判断 SQLite 唯一约束冲突。
func isUniqueViolation(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "UNIQUE constraint failed") ||
		strings.Contains(msg, "constraint failed: UNIQUE")
}
