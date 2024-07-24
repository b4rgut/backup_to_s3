package filesystem

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

// LS возвращает список файлов и директорий в указанной директории.
//
// Параметры:
// - path: строка, содержащая путь к директории, которую нужно просканировать.
//
// Возвращает:
// - []fs.DirEntry: срез, содержащий информацию о файлах и директориях в указанной директории.
//
// Ошибки:
// - В случае ошибки открытия директории или чтения её содержимого, функция завершает выполнение с логированием фатальной ошибки.
func LS(path string) []fs.DirEntry {
	log.Printf("сканирование дирректории %s\n", path)

	dir, err := os.Open(path)
	if err != nil {
		log.Fatalf("ошибка открытия дирректории %s: %v", path, err)
	}
	defer dir.Close()

	files, err := dir.ReadDir(-1)
	if err != nil {
		log.Fatalf("ошибка получения файлов дирректории: %v", err)
	}

	if len(files) == 0 {
		log.Fatalln("дирректория пуста.")
	}

	return files
}

func GetExecutablePath() string {
	exePath, err := os.Executable()
	if err != nil {
		log.Fatalf("не удалось получить путь к исполняемому файлу: %v", err)
	}

	return filepath.Dir(exePath)
}
