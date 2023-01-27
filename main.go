package main

import (
	"bytes"
	_ "embed"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func main() {
	if err := run(os.Args[1:]...); err != nil {
		log.Fatal(err)
	}
}

//go:embed template.bin
var template []byte

func run(args ...string) error {
	if len(args) < 1 {
		return errors.New("need input file as the first argument")
	}

	filename := args[0]

	// 5f 10 8e x 2 (8e = size of body)
	// 142 bytes

	b, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	dir := "/Users/bep/Music/Audio Music Apps/Plug-In Settings/Scripter"
	_, name := filepath.Split(filename)
	baseName := filepath.Base(name)
	// Remove extension.
	baseName = baseName[:len(baseName)-len(filepath.Ext(baseName))]
	outName := filepath.Join(dir, baseName+".pst")

	// Offsets
	// Size: 142 bytes
	// Content offsets:
	// 4123
	// 4317
	// hex: 0, 100c, 101a, 10c9, 10ce, 10dc, 118b
	// decimal; 0, 4108, 4122, 4297, 4302, 4316, 4491
	// values hex: 98, c2, 8e, 99, c2, 8e, 99
	// values decimal: 152, 194, 142, 153, 194, 142, 153
	// values offset from size (142): 10, 52, 0, 11, 52, 0, 11

	var (
		valueOffsetsFromSize = []int{10, 52, 0, 11, 52, 0, 11}
		offsets              = []int{0, 4108, 4122, 4297, 4302, 4316, 4491}
	)

	// Adjust offsets considering the content.
	/*for i, offset := range offsets {
		if offset > 4123 {
			//offsets[i] += len(b)
		}
		if offset > 4317 {
			//offsets[i] += len(b)
		}
	}*/

	outFile, err := os.Create(outName)
	if err != nil {
		return err
	}
	defer outFile.Close()

	//lenb := len(b)

	var contentPlaceholder = []byte("_SCRIPTER_SCRIPT_CONTENT")
	patchContent := bytes.ReplaceAll(template, contentPlaceholder, b)

	// Replace offsets.
	for i, offset := range offsets {
		// Adjust offset considering the content.
		patchContent[offset] = byte(valueOffsetsFromSize[i] + len(b))
	}

	// 98
	fmt.Printf("First byte: %x\n", patchContent[0])

	if _, err := outFile.Write(patchContent); err != nil {
		return err
	}

	return nil
}

func mustDecode(s string) []byte {
	b, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return b
}
