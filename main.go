package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
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
	date() time.Time
}

func main() {
	usage := `Photo Sort - sorts your photos

Usage:
  photosort (year|month|day|hour|minute) <inputdir> [--recursive] [--cleanup]

Options:
  -h --help     Show this screen.
  -v, --verbose Show whats being done
  -d, --dryrun  Simulate sorting
  --version     Show version.`

	arguments, err := docopt.ParseArgs(usage, os.Args[1:], VERSION)
	if err != nil {
		log.Fatal(err)
	}
	var photoList []*Photo
	searchPath := path.Clean(arguments["<inputdir>"].(string))
	photoList = gatherFiles(searchPath, photoList,
		arguments["--recursive"].(bool))
	var granularity string
	for _, g := range []string{"year", "month", "day", "hour", "minute"} {
		if arguments[g].(bool) {
			granularity = g
		}
	}
	for _, p := range photoList {
		p.load()
		p.move(granularity)
	}
}

func gatherFiles(searchPath string, photoList []*Photo,
	recursive bool) []*Photo {
	allowedExtensions := []string{".jpg", ".jpeg"}
	files, err := ioutil.ReadDir(searchPath)
	if err != nil {
		log.Fatal(err)
	}
	if photoList == nil {
		photoList = make([]*Photo, 0, len(files))
	}
	for _, f := range files {
		if f.IsDir() {
			if recursive {
				photoList = gatherFiles(fmt.Sprintf("%s/%s", searchPath,
					f.Name()), photoList, recursive)
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
		photoList = append(photoList,
			&Photo{path: fmt.Sprintf("%s/%s", searchPath, f.Name()),
				exifData: nil})
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

func (p *Photo) fileName() string {
	pathSplit := strings.Split(p.path, "/")
	return pathSplit[len(pathSplit)-1]
}

func (p *Photo) fileDir() string {
	pathSplit := strings.Split(p.path, "/")
	return strings.Join(pathSplit[:len(pathSplit)-1], "/")
}

func (p *Photo) move(granularity string) {
	var destPath string
	switch granularity {
	case "year":
		destPath = fmt.Sprintf("%s/%d", p.fileDir(), p.date().Year())
	case "month":
		destPath = fmt.Sprintf("%s/%d/%d", p.fileDir(), p.date().Year(),
			p.date().Month())
	case "day":
		destPath = fmt.Sprintf("%s/%d/%d/%d", p.fileDir(), p.date().Year(),
			p.date().Month(),
			p.date().Day())
	case "hour":
		destPath = fmt.Sprintf("%s/%d/%d/%d/%d", p.fileDir(), p.date().Year(),
			p.date().Month(),
			p.date().Day(), p.date().Hour())
	case "minute":
		destPath = fmt.Sprintf("%s/%d/%d/%d/%d/%d", p.fileDir(),
			p.date().Year(), p.date().Month(), p.date().Day(), p.date().Hour(),
			p.date().Minute())
	default:
		log.Fatalf("unknown granularity: %s", granularity)
	}

	// err := os.MkdirAll(destPath, os.FileMode(int(0755)))
	// if err != nil {
	// 	log.Fatal(err)
	// }

	fmt.Printf("%s => %s/%s\n", p.path, destPath, p.fileName())
}
