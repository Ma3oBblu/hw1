package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type File struct {
	Name string
}

type Folder struct {
	Name    string
	Files   []File
	Folders []Folder
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	files, err := oneDir(path)
	tree := Folder{Name: path}
	if err != nil {
		fmt.Errorf("error")
	}
	for _, file := range files {
		if file.IsDir() {
			newFolder := Folder{Name: file.Name()}
			tree.Folders = pushtree.FoldersnewFolder
			fmt.Println(file.Name())
			fmt.Println(walker(path + "/" + file.Name()))
		} else {
			newFile := File{Name: file.Name()}
			tree.Files = newFile
		}
	}

	return err
}

func oneDir(path string) ([]os.FileInfo, error) {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	files, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		log.Fatal(err)
	}
	return files, nil
}

func walker(path string) error {
	err := filepath.Walk(path,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			stringArray := strings.Split(path, "/")
			last := stringArray[len(stringArray)-1:]
			fmt.Println(last)
			return nil
		})
	return err
}

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}
