package main

import (
	"fmt"
	"os"
)

const (
	SIZE_KiB = 1024
	SIZE_MiB = 1048576
	SIZE_GiB = 1073741824
	SIZE_TiB = 1099511627776
)

func GetDir(path string) ([]os.FileInfo, error) {

	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		return nil, err
	}

	fileInfo, err := f.Readdir(-1)
	if err != nil {
		return nil, err
	}

	return fileInfo, nil
}

func main() {
	LoadConfig()
	ProcessCli()
}

///Helpers
func FormatSize(size int64) string {

	if size < SIZE_KiB {
		return fmt.Sprintf("%d Bytes", size)

	} else if size < SIZE_MiB {
		value := float64(size) / float64(SIZE_KiB)
		return fmt.Sprintf("%.2f KiB", value)

	} else if size < SIZE_GiB {
		value := float64(size) / float64(SIZE_MiB)
		return fmt.Sprintf("%.2f MiB", value)

	} else if size < SIZE_TiB {
		value := float64(size) / float64(SIZE_GiB)
		return fmt.Sprintf("%.2f GiB", value)

	} else {
		value := float64(size) / float64(SIZE_TiB)
		return fmt.Sprintf("%.2f TiB", value)
	}

}
