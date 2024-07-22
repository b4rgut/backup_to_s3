package crypto

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"log"
	"os"
)

// ComputeFileMD5 вычисляет MD5-хеш для указанного файла.
//
// Параметры:
// - filePath: строка, содержащая путь к файлу, для которого нужно вычислить MD5-хеш.
//
// Возвращает:
// - строка, содержащая MD5-хеш файла в шестнадцатеричном формате. Если произошла ошибка, возвращает пустую строку.
//
// Ошибки:
// - В случае ошибки открытия файла или вычисления хэша, функция логирует ошибку и возвращает пустую строку.
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
