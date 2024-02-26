package pkg

import (
	"crypto/rand"
	"encoding/hex"
)

func RandomHex(n int) (string, error) {
	bytesSize := n / 2
	if n%2 != 0 {
		bytesSize++
	}

	bytes := make([]byte, bytesSize)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	hexString := hex.EncodeToString(bytes)
	if len(hexString) > n {
		hexString = hexString[:n]
	}

	return hexString, nil
}
