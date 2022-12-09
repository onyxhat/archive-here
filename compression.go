package main

import (
	"bytes"
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/itrepablik/kopy"
)

func compress(srcDir string, dstDir string, exclude []string) {
	defer wg.Done()

	// Compose paths using os-agnostic path delimiters
	src := filepath.FromSlash(srcDir)
	dst := filepath.FromSlash(dstDir)

	// Composing the archive filename
	srcBaseName := kopy.FileNameWOExt(filepath.Base(src))
	archiveDir := srcBaseName + kopy.ComFileFormat
	archiveDest := filepath.FromSlash(path.Join(dst, archiveDir))

	var buf bytes.Buffer
	if err := kopy.CompressDIR(src, &buf, exclude); err != nil {
		panic(err)
	}

	fileToWrite, err := os.OpenFile(archiveDest, os.O_CREATE|os.O_RDWR, os.FileMode(600))
	if err != nil {
		panic(err)
	}

	if _, err := io.Copy(fileToWrite, &buf); err != nil {
		panic(err)
	}
	defer fileToWrite.Close()

	if err := os.RemoveAll(src); err != nil {
		panic(err)
	}
}
