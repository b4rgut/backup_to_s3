package main

import (
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"cloud_backuper/internal/client/s3"
	"cloud_backuper/internal/client/yadisk"
	"cloud_backuper/internal/configs"
	"cloud_backuper/internal/crypto"
	"cloud_backuper/internal/filesystem"
)

func main() {
	startTime := time.Now()

	logPath := filepath.Join(filesystem.GetExecutablePath(), "cloud_upload.log")

	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("не удалось открыть файл логов: %v", err)
	}
	defer logFile.Close()

	log.SetOutput(logFile)

	log.Println("начало загрузки бэкапов...")

	config := configs.LoadConfigs()

	var s3Client *s3.S3Client
	if config.S3.Enabled {
		s3Client = s3.New(&config.S3)
	}

	var ydClient *yadisk.YDClient
	if config.YD.Enabled {
		ydClient = yadisk.New(config.YD.Token)
	}

	files := filesystem.LS(config.LocalDirectoryPath)

	for _, file := range files {
		for _, dir := range config.DirectoryStruct {
			if strings.HasPrefix(file.Name(), dir.PrefixFile) {
				fullPath := filepath.Join(config.LocalDirectoryPath, file.Name())
				hash := crypto.ComputeFileMD5(fullPath)

				S3ObjectName := path.Join(dir.CloudDir, file.Name())
				if config.S3.Enabled && !s3Client.FileExists(hash, config.S3.Backet, S3ObjectName) {
					s3Client.UploadFile(config.S3.Backet, S3ObjectName, fullPath)
				}

				YDPath := path.Join(config.YD.PathToBackup, dir.CloudDir, file.Name())
				if config.YD.Enabled && !ydClient.FileExists(hash, filepath.Join("Backup1C", dir.CloudDir, file.Name())) {
					ydClient.UploadFile(fullPath, YDPath)
				}

				break
			}
		}

	}

	endTime := time.Now()
	log.Printf("загрузка бэкапов завершена! Время выполнения %v\n\n", endTime.Sub(startTime))
}
