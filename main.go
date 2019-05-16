package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/docopt/docopt-go"
)

// VERSION of the program
const VERSION = "1"

func main() {
	usage := `Photo Sort - sorts your photos

Usage:
  photosort (year|month|day|hour|minute) <inputdir> [--recursive] [--cleanup]

Options:
  -h --help     Show this screen.
  --version     Show version.`

	arguments, err := docopt.ParseArgs(usage, os.Args[1:], VERSION)
	if err != nil {
		log.Fatal(err)
	}
	fileList := make([]string, 10)
	fileList = gatherFiles(arguments["<inputdir>"].(string), fileList, arguments["--recursive"].(bool))
	fmt.Println(fileList)
}

func gatherFiles(path string, fileList []string, recursive bool) []string {
	allowedExtensions := []string{".jpg", ".jpeg"}
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		if f.IsDir() {
			if recursive {
				fileList = gatherFiles(fmt.Sprintf("%s/%s", path, f.Name()), fileList, recursive)
			}
			continue
		}
		hasAllowedExtension := false
		for _, i := range allowedExtensions {
			if strings.HasSuffix(strings.ToLower(f.Name()), i) {
				hasAllowedExtension = true
			}
		}
		if !hasAllowedExtension {
			continue
		}
		fileList = append(fileList, f.Name())
	}
	return fileList
}
