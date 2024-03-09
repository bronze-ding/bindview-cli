package src

import (
	"fmt"
	"os"
)

func NewFileName(oldFile string, newFile string) {
	err := os.Rename(oldFile, newFile)
	if err != nil {
		fmt.Println(err)
		return
	}
}
