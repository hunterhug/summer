package stool

import (
	"crypto/md5"
	"errors"
	"fmt"
	"hash/crc64"
	"io"
	"os"
)

func Md5Simple(raw []byte) (string, error) {
	h := md5.New()
	num, err := h.Write(raw)
	if err != nil {
		return "", err
	}
	if num == 0 {
		return "", errors.New("num 0")
	}
	data := h.Sum([]byte(""))
	return fmt.Sprintf("%x", data), nil
}

func Crc64File(filePath string) (uint64, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
		}
	}(file)

	table := crc64.MakeTable(crc64.ECMA)
	hash := crc64.New(table)

	if _, err := io.Copy(hash, file); err != nil {
		return 0, err
	}

	return hash.Sum64(), nil
}

func Crc64Binary(data []byte) uint64 {
	table := crc64.MakeTable(crc64.ECMA)
	hash := crc64.New(table)
	_, _ = hash.Write(data)
	return hash.Sum64()
}
