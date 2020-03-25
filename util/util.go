package util

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
)

func GetTempDir(tempDir, suffix string) string {
	dir := filepath.Join(tempDir, suffix)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, os.ModePerm)
	} else {
		os.RemoveAll(dir)
	}

	return dir
}

// GetHash receives one or more strings and returns
// a sha1 hash of the concatenation of all of them
func GetHash(values ...string) string {
	str := values[0]
	for i := 1; i < len(values); i++ {
		str = fmt.Sprintf("%s%s", str, values[i])
	}
	hashString := sha1.New()
	hashString.Write([]byte(str))
	return hex.EncodeToString(hashString.Sum(nil))
}

func GetEntropy(data string) (entropy float64) {
	if data == "" {
		return 0
	}

	for i := 0; i < 256; i++ {
		px := float64(strings.Count(data, string(byte(i)))) / float64(len(data))
		if px > 0 {
			entropy += -px * math.Log2(px)
		}
	}

	return entropy
}

// Obfuscate changes the last 1/3 string characters by *'s
func Obfuscate(text string) string {
	size := len(text)
	return text[0:2*size/3] + strings.Repeat("*", len(text)-2*size/3)
}

func IsX9Action(action string) bool {
	listActions := []string{"opened", "edited", "reopened"}
	for _, v := range listActions {
		if action == v {
			return true
		}
	}
	return false
}
