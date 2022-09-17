package repos

import (
	"log"
	"os"
)

func safelyClose(file *os.File) {
	err := file.Close()
	if err != nil {
		log.Fatal("Problem with closing a data file.")
	}
}
