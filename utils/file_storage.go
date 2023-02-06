package utils

import (
	"crypto/md5"
	"encoding/hex"
	"io"
)

func FileMD5(reader io.Reader) (string, error) {
	hash := md5.New()
	_, _ = io.Copy(hash, reader)
	return hex.EncodeToString(hash.Sum(nil)), nil
}

func GetStoragePath(md5 string) string {
	return md5[:2] + "/" + md5[2:]
}
