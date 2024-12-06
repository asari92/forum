package utils

import (
	"bufio"
	"os"
	"strings"
)

// LoadEnv читает переменные окружения из файла и устанавливает их.
func LoadEnv(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 || line[0] == '#' { // Игнорируем пустые строки и комментарии
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue // Игнорируем строки, которые не могут быть разделены на ключ и значение
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		if err := os.Setenv(key, value); err != nil {
			return err
		}
	}
	return scanner.Err()
}
