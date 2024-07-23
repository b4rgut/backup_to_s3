package yadisk

import (
	"context"
	"log"

	"github.com/ilyabrin/disk"
)

type YDClient struct {
	*disk.Client
}

func New(token string) *YDClient {
	client := disk.New(token)

	if _, err := client.Disk.Info(context.Background(), nil); err != nil {
		log.Fatal(err)
	}

	return &YDClient{client}
}

func (client *YDClient) FileExists(localFileMD5, cloudPath string) bool {
	res, err := client.Resources.Meta(context.Background(), cloudPath, nil)
	if err != nil {
		log.Printf("ошибка получения мета информации файла %v: %v", cloudPath, err)
		return false
	}

	if res.Md5 != localFileMD5 {
		log.Printf("загрузка файла на Яндекс Диск не требуется, так как он уже загружен %s", cloudPath)
		return false
	}

	return true
}

func (client *YDClient) UploadFile(localPath, cloudPath string) {
	link, err := client.Resources.GetUploadLink(context.Background(), cloudPath)
	if err != nil {
		log.Printf("ошибка получения ссылки для загрузки файла %v на Яндекс Диск: %v", cloudPath, err)
		return
	}

	if err = client.Resources.Upload(context.Background(), localPath, link.Href, nil); err != nil {
		log.Printf("ошибка загрузки файла %v на Яндекс Диск: %v", localPath, err)
	}
}
