package filesystem

import (
	"io/fs"
	"log"
	"os"
)


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

	return files
}
