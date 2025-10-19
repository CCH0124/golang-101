package main

import (
	"fmt"
	"log"
	"os"
)

func fileLen(fileName string) (int, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return 0, err
	}
	defer file.Close()
	fileInfo, err := file.Stat()
	if err != nil {
		return 0, err
	}
	return int(fileInfo.Size()), nil
}
func main() {
	if len(os.Args) < 2 {
		return
	}
	size, err := fileLen(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(size)
}
