package main

import (
	"log"
	"path"
	"path/filepath"
	"s3backuper/internal/client/s3"
	"s3backuper/internal/configs"
	"s3backuper/internal/crypto"
	"s3backuper/internal/filesystem"
	"strings"

	"github.com/danieljoos/wincred"
)

func main() {
	log.Println("начало загрузки бэкапов...")

	cred, err := wincred.GetGenericCredential("s3backup")
	if err != nil {
		log.Fatalf("ошибка получения учетных данных: %v", err)
	}
	if cred == nil {
		log.Fatalf("учетные данные не найдены")
	}

	secretAccessKey := strings.ReplaceAll(string(cred.CredentialBlob), "\x00", "")

	config := configs.NewS3Config(
		"ru-msk-dr3-1.store.cloud.mts.ru",
		cred.UserName,
		secretAccessKey,
		true,
	)

	s3Client := s3.New(config)

	backetName := "backup1c"

	localDirPath := "D:/Test/sql"

	files := filesystem.LS(localDirPath)

	for _, file := range files {
		subDir := ""

		switch {
		case strings.HasPrefix(file.Name(), "HRM"):
			subDir = "HRM"
		case strings.HasPrefix(file.Name(), "Accounting"):
			subDir = "Accounting"
		case strings.HasPrefix(file.Name(), "unf_30"):
			subDir = "unf_30"
		default:
			continue
		}

		objectName := path.Join(subDir, file.Name())
		localFilePath := filepath.Join(localDirPath, file.Name())
		hash := crypto.ComputeFileMD5(localFilePath)

		if !s3Client.FileExists(hash, backetName, objectName) {
			log.Printf("загрузка файла %s в %s", localFilePath, objectName)
			s3Client.UploadFile(backetName, objectName, localFilePath)
		}
	}

	log.Println("конец загрузки бэкапов!")
}
