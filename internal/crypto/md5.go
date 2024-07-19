package crypto

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"log"
	"os"
)

func ComputeFileMD5(filePath string) string {
	log.Printf("вычисление хеша файла %s", filePath)
	
	file, err := os.Open(filePath)
	if err != nil {
		return ""
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		log.Printf("ошибка вычисления хэша файла %s: %v", filePath, err)
		return ""
	}

	return hex.EncodeToString(hash.Sum(nil))
}