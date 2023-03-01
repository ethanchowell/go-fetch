package util

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func ValidateChecksum(checksum string, data []byte) (bool, error) {
	if checksum == "" {
		return true, nil
	}

	hash := sha256.New()
	hash.Write(data)

	assetSum := hex.EncodeToString(hash.Sum(nil))
	if assetSum != checksum {
		return false, fmt.Errorf("checksum mismatch: %s %s", assetSum, checksum)
	}

	return true, nil
}
