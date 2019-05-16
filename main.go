package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/docopt/docopt-go"
	"github.com/rwcarlsen/goexif/exif"
)

// VERSION of the program
const VERSION = "1"

// Photo is a picture of stb or sth
type Photo struct {
	path     string
	exifData *exif.Exif
}

// ExifReader provides exif data from a photo
type ExifReader interface {
	load()
	date() time.Time
}

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
	var photoList []*Photo
	photoList = gatherFiles(arguments["<inputdir>"].(string), photoList, arguments["--recursive"].(bool))
	for _, p := range photoList {
		p.load()
		fmt.Println(p.date().Year())
	}
}

func gatherFiles(path string, photoList []*Photo, recursive bool) []*Photo {
	allowedExtensions := []string{".jpg", ".jpeg"}
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	if photoList == nil {
		photoList = make([]*Photo, 0, len(files))
	}
	for _, f := range files {
		if f.IsDir() {
			if recursive {
				photoList = gatherFiles(fmt.Sprintf("%s/%s", path, f.Name()), photoList, recursive)
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
		photoList = append(photoList, &Photo{path: fmt.Sprintf("%s/%s", path, f.Name()), exifData: nil})
	}
	return photoList
}

func (p *Photo) load() {
	fHandle, err := os.Open(p.path)
	defer fHandle.Close()
	if err != nil {
		log.Fatal(err)
	}
	p.exifData, err = exif.Decode(fHandle)
	if err != nil {
		log.Fatal(err)
	}
}

func (p *Photo) date() time.Time {
	time, err := p.exifData.DateTime()
	if err != nil {
		log.Fatal(err)
	}
	return time
}
