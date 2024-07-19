package s3

import (
	"context"
	"log"
	"s3backuper/internal/configs"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type S3Client struct {
	*minio.Client
}

// New создает новый экземпляр S3Client с заданной конфигурацией S3Config.
//
// Параметры:
// - config: указатель на структуру S3Config, содержащую конфигурацию для подключения к S3.
//
// Возвращает:
// - *S3Client: указатель на новый экземпляр S3Client.
func New(config *configs.S3Config) *S3Client {
	minioClient, err := minio.New(config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccessKeyID, config.SecretAccessKey, ""),
		Secure: config.UseSSL,
	})

	if err != nil {
		log.Fatalf("Ошибка инициализации S3 клиента: %v", err)
	}

	return &S3Client{minioClient}
}

// FileExists проверяет, существует ли файл в S3 и совпадает ли его MD5-хеш с локальным файлом.
//
// Параметры:
// - localFileMD5: строка, содержащая MD5-хеш локального файла.
// - bucket: строка, содержащая имя бакета в S3.
// - objectName: строка, содержащая имя объекта в S3.
//
// Возвращает:
// - bool: true, если файл существует и его MD5-хеш совпадает с локальным файлом, иначе false.
func (client *S3Client) FileExists(localFileMD5, bucket, objectName string) bool {
	objectInfo, err := client.StatObject(context.Background(), bucket, objectName, minio.StatObjectOptions{})
	if err != nil {
		if minio.ToErrorResponse(err).Code == "NoSuchKey" {
			return false
		}
		log.Fatalf("Ошибка поиска файла в S3 хранилище: %v", err)
	}

	actualMD5 := strings.Trim(objectInfo.ETag, "\"")
	
	if actualMD5 != localFileMD5 {
		return false
	}
	
	log.Printf("загрузка файла не требуется, так как он уже загружен %s", objectName)

	return true
}


// UploadFile загружает файл в указанный бакет S3.
//
// Параметры:
// - bucket: строка, содержащая имя бакета в S3.
// - objectName: строка, содержащая имя объекта в S3.
// - path: строка, содержащая путь к локальному файлу.
//
// Ошибки:
// - В случае ошибки загрузки файла, ошибка будет залогирована.
func (client *S3Client) UploadFile(bucket, objectName, path string) {
	_, err := client.FPutObject(context.Background(), bucket, objectName, path, minio.PutObjectOptions{
		ContentType: "application/octet-stream",
	})

	if err != nil {
		log.Printf("ошибка загрузки файла %s: %v", objectName, err)
	}
}
