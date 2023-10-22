package main

import (
	"fmt"
	"github.com/EdmondTabaku/gompressor/compressor"
	"os"
	"path/filepath"
)

const ext = ".gozip"

func main() {

	if len(os.Args) < 3 {
		fmt.Println("Usage: gompressor [zip/unzip] [file/dir] [new file name]")
		os.Exit(1)
	}

	operation := os.Args[1]
	path := os.Args[2]

	var exportName = ""
	if len(os.Args) > 3 {
		exportName = os.Args[3]
	}

	fileBase := filepath.Base(path)
	fileExt := filepath.Ext(fileBase)
	fileName := fileBase[0 : len(fileBase)-len(fileExt)]
	dir := filepath.Dir(path)

	if path == "" {
		fmt.Println("Path cannot be empty")
		os.Exit(1)
	}

	comp := compressor.NewCompressorBase(fileExt)

	switch operation {
	case "zip":
		file, err := os.Open(path)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			os.Exit(1)
		}

		compressed, err := comp.Compress(file)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			os.Exit(1)
		}

		if exportName != "" {
			fileName = exportName
		}

		err = os.WriteFile(dir+`\`+fileName+ext, []byte(compressed), 0644)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			os.Exit(1)
		}

	case "unzip":
		file, err := os.ReadFile(path)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			os.Exit(1)
		}

		decompressed, ex, err := comp.Decompress(string(file))
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			os.Exit(1)
		}

		if exportName != "" {
			fileName = exportName
		}

		err = os.WriteFile(dir+`\`+fileName+ex, []byte(decompressed), 0644)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			os.Exit(1)
		}
	}

}
