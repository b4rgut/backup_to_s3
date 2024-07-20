package main

import (
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"s3backuper/internal/client/s3"
	"s3backuper/internal/configs"
	"s3backuper/internal/crypto"
	"s3backuper/internal/filesystem"

	"github.com/danieljoos/wincred"
)

func init() {
	logFile, err := os.OpenFile("cloud_upload.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Не удалось открыть файл логов: %v", err)
	}

	log.SetOutput(logFile)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

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

	localDirPath := "D:/Test/"

	files := filesystem.LS(localDirPath)

	for _, file := range files {
		subDir := ""

		switch {
		case strings.HasPrefix(file.Name(), "HRM_backup"):
			subDir = "HRM"
		case strings.HasPrefix(file.Name(), "Accounting_backup"):
			subDir = "Accounting"
		case strings.HasPrefix(file.Name(), "unf_30_backup"):
			subDir = "UNF"
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
