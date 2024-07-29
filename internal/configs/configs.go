package configs

import (
	"cloud_backuper/internal/filesystem"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/danieljoos/wincred"

	"gopkg.in/yaml.v3"
)

// DirectoryItem содержит информацию о префиксе локального файла и каталоге в облаке, в который его нужно загружать.
// - PrefixFile: строка для хранения префикса имени локального файла.
// - CloudDir: дирректория для хранения этого файла в S3.
type DirectoryItem struct {
	PrefixFile string `yaml:"prefix_file"`
	CloudDir   string `yaml:"cloud_dir"`
}

// S3Config содержит конфигурацию для подключения к S3.
// - Endpoint: строка, содержащая URL-адрес конечной точки S3.
// - Backet: строка, содержащая имя бакета в S3.
// - AccessKeyID: строка, содержащая идентификатор ключа доступа к S3.
// - SecretAccessKey: строка, содержащая секретный ключ доступа к S3.
// - UseSSL: логическое значение, указывающее, следует ли использовать SSL для подключения к S3.
type S3Config struct {
	Enabled         bool   `yaml:"enabled"`
	Endpoint        string `yaml:"endpoint"`
	Backet          string `yaml:"backet"`
	CredentialName  string `yaml:"credential_name"`
	AccessKeyID     string
	SecretAccessKey string
	UseSSL          bool
}

// YDConfig содержит конфигурацию для подключения к Яндекс Диску.
// - Token: строка, содержащая токен доступа к Яндекс Диску.
// - PathToBackup: строка, содержащая путь к каталогу для бэкапов на Яндекс Диске.
type YDConfig struct {
	Enabled        bool   `yaml:"enabled"`
	PathToBackup   string `yaml:"path_to_backup"`
	CredentialName string `yaml:"credential_name"`
	Token          string
}

// Config содержит конфигурацию для загрузки файлов в S3 и локального каталога.
// - S3: структура S3Config, содержащая конфигурацию для подключения к S3.
// - YD: содержит конфигурацию для подключения к Яндекс Диску.
// - LocalDirectoryPath: строка, содержащая путь к локальному каталогу для загрузки файлов.
// - WindowsCredential: строка, содержащая имя учетных данных в Windows Credential для доступа к S3.
// - DirectoryStruct: массив структур DirectoryItem, содержащих информацию о префиксах локальных файлов
// и каталогах в облаке, в которые их нужно загружать.
type Config struct {
	S3                 S3Config        `yaml:"s3"`
	YD                 YDConfig        `yaml:"yandex_disk"`
	LocalDirectoryPath string          `yaml:"local_directory_path"`
	DirectoryStruct    []DirectoryItem `yaml:"directory_struct"`
}

// LoadConfigs загружает конфигурацию из файла config.yml и учетных данных Windows Credential.
//
// Возвращает:
// - *Config: указатель на структуру Config, содержащую загруженную конфигурацию.
//
// Ошибки:
// - В случае ошибки чтения файла конфигурации или декодирования YAML, функция завершает выполнение с логированием фатальной ошибки.
// - В случае ошибки получения учетных данных из Windows Credential, функция завершает выполнение с логированием фатальной ошибки.
func LoadConfigs() *Config {
	configPath := filepath.Join(filesystem.GetExecutablePath(), "config.yml")
	data, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("не удалось прочитать файл конфигурации: %v", err)
	}

	config := Config{}

	if err := yaml.Unmarshal(data, &config); err != nil {
		log.Fatalf("не удалось декодировать YAML: %v", err)
	}

	if !config.S3.Enabled && !config.YD.Enabled {
		log.Fatalln("ни один из сервисов (S3 или Яндекс Диск) не включен!")
	}

	if config.S3.Enabled {
		cred, err := wincred.GetGenericCredential(config.S3.CredentialName)
		if err != nil && cred == nil {
			log.Fatalf("ошибка получения учетных данных S3 из Windows Credential: %v", err)
		}

		config.S3.AccessKeyID = cred.UserName
		config.S3.SecretAccessKey = strings.ReplaceAll(string(cred.CredentialBlob), "\x00", "")
		config.S3.UseSSL = true
	}

	if config.YD.Enabled {
		cred, err := wincred.GetGenericCredential(config.YD.CredentialName)
		if err != nil && cred == nil {
			log.Fatalf("ошибка получения учетных данных Яндекс Диска из Windows Credential: %v", err)
		}

		config.YD.Token = strings.ReplaceAll(string(cred.CredentialBlob), "\x00", "")
	}

	return &config
}
