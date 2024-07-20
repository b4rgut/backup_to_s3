package configs

import (
	"os"
	"log"
	"strings"

    "gopkg.in/yaml.v3"
    "github.com/danieljoos/wincred"
)

type DirectoryItem struct {
	PrefixFile string    `yaml:"prefix_file"`
	CloudDir      string `yaml:"cloud_dir"`
}

type S3Config struct {
	Endpoint        string `yaml:"endpoint"`
	Backet          string `yaml:"backet"`
	AccessKeyID     string
	SecretAccessKey string
	UseSSL          bool
}

type Config struct {
	S3                  S3Config 		`yaml:"s3"`
	LocalDirectoryPath  string          `yaml:"local_directory_path"`
	WindowsCredential   string          `yaml:"windows_credential"`
	DirectoryStruct     []DirectoryItem `yaml:"directory_struct"`
}

func LoadConfigs() *Config {
	data, err := os.ReadFile("config.yml")
	if err != nil {
		log.Fatalf("Не удалось прочитать файл конфигурации: %v", err)
	}
	
	config := Config{}
	
	if err := yaml.Unmarshal(data, &config); err != nil {
		log.Fatalf("Не удалось декодировать YAML: %v", err)
	}
	
	cred, err := wincred.GetGenericCredential(config.WindowsCredential)
	if err != nil && cred == nil {
		log.Fatalf("ошибка получения учетных данных из Windows Credential: %v", err)
	}
	
	config.S3.AccessKeyID = cred.UserName
	config.S3.SecretAccessKey = strings.ReplaceAll(string(cred.CredentialBlob), "\x00", "")
	config.S3.UseSSL = true
	
	return &config
}

