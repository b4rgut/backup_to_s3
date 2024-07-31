package crypto

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
)

// ComputeFileETag вычисляет ETag для указанного файла.
//
// Параметры:
// - filePath: строка, содержащая путь к файлу, для которого нужно вычислить ETag.
// - partSize: размер каждой части файла в байтах.
//
// Возвращает:
// - строка, содержащая ETag файла в формате MD5SUM-N. Если произошла ошибка, возвращает пустую строку.
//
// Ошибки:
// - В случае ошибки открытия файла или вычисления хэша, функция логирует ошибку и возвращает пустую строку.
func ComputeFileETag(filePath string, partSize uint64) string {
	log.Printf("вычисление хэша файла: %s", filePath)

	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("не удалось открыть файл: %v", err)
		return ""
	}
	defer file.Close()

	var hashes []byte
	buffer := make([]byte, partSize)

	parts := 0
	for {
		n, err := file.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("ошибка чтения файла: %v", err)
			return ""
		}

		if n > 0 {
			hash := md5.Sum(buffer[:n])
			hashes = append(hashes, hash[:]...)
			parts++
		}
	}

	// Объединяем все хэши в одну строку

	if parts == 1 {
		return hex.EncodeToString(hashes)
	}

	hash := md5.New()
	_, err = hash.Write(hashes)
	if err != nil {
		return ""
	}

	return fmt.Sprintf("%s-%d", hex.EncodeToString(hash.Sum(nil)), parts)
}
