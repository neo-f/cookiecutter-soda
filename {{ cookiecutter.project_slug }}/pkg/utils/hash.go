package utils

import (
	"crypto/md5"
	"encoding/hex"
)

func MD5(d []byte) string {
	hash := md5.Sum(d)
	return hex.EncodeToString(hash[:])
}
