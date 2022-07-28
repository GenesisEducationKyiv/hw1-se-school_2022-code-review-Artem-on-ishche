package data

import "os"

// safelyClose loops while the file is not closed successfully.
func safelyClose(file *os.File) {
	err := file.Close()
	for err != nil {
		err = file.Close()
	}
}
