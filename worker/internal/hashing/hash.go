package hashing

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"strings"
)

func Mine(input string) (string, int64) {
	var nonce int64
	for {
		data := input + strconv.FormatInt(nonce, 10)
		hash := sha256.Sum256([]byte(data))
		hexHash := hex.EncodeToString(hash[:])

		if strings.HasPrefix(hexHash, "0000") {
			return hexHash, nonce
		}
		nonce++
	}
}
