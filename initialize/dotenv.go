package initialize

import (
	"bufio"
	"log"
	"os"
	"strings"
)

func LoadDotEnvFile() {
	file, err := os.Open(".env")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "=")

		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])

			os.Setenv(key, value)
		}

	}
}
