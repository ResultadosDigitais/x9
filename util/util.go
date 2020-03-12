package util

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/ResultadosDigitais/x9/log"
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

func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}

	if os.IsNotExist(err) {
		return false
	}

	return false
}

func LogIfError(text string, err error) {
	if err != nil {
		log.Error(fmt.Sprintf("%s (%s", text, err.Error()), nil)
	}
}

func GetHash(values ...string) string {
	str := values[0]
	for i := 1; i < len(values); i++ {
		str = fmt.Sprintf("%s%s", str, values[i])
	}
	hashString := sha1.New()
	hashString.Write([]byte(str))
	return hex.EncodeToString(hashString.Sum(nil))
}

func Pluralize(count int, singular string, plural string) string {
	if count == 1 {
		return singular
	}

	return plural
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

func Obfuscate(text string) string {
	size := len(text)
	return text[0:2*size/3] + strings.Repeat("*", size/3)
}
