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

func main() {
	usage := `Photo Sort - sorts your photos

Usage:
  photosort [options] (year|month|day|hour|minute) <inputdir>

Options:
  -h --help        Show this screen.
  -v, --verbose    Show whats being done
  -d, --dryrun     Simulate sorting
  -r, --recursive  Scan directories recursively
  -c, --cleanup    Remove empty source directories when run recursively
  --version        Show version.`

	arguments, err := docopt.ParseArgs(usage, os.Args[1:], VERSION)
	if err != nil {
		log.Fatal(err)
	}
	verbose := arguments["--verbose"].(bool)
	dryrun := arguments["--dryrun"].(bool)
	doCleanup := arguments["--cleanup"].(bool)
	recursive := arguments["--recursive"].(bool)
	searchPath := path.Clean(arguments["<inputdir>"].(string))

	var photoList []*Photo
	photoList = gatherFiles(searchPath, photoList, dryrun)
	var granularity string
	for _, g := range []string{"year", "month", "day", "hour", "minute"} {
		if arguments[g].(bool) {
			granularity = g
		}
	}
	for _, p := range photoList {
		p.load()
		p.move(granularity, dryrun, verbose || dryrun)
	}
	if doCleanup && recursive && !dryrun {
		cleanup(searchPath, verbose)
	}
}

func cleanup(rootPath string, verbose bool) {
	files, err := ioutil.ReadDir(rootPath)
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		if f.IsDir() {
			fls, err := ioutil.ReadDir(rootPath)
			relPath := path.Join(rootPath, f.Name())
			if err != nil {
				log.Fatal(err)
			}
			if len(fls) == 0 {
				if verbose {
					fmt.Printf("Removing %s", relPath)
				}
				err = os.Remove(relPath)
				if err != nil {
					log.Fatal(err)
				}
			} else {
				cleanup(relPath, verbose)
			}
		}
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

func (p *Photo) move(granularity string, dryrun bool, verbose bool) {
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

	if verbose {
		fmt.Printf("%s => %s/%s\n", p.path, destPath, p.fileName())
	}

	if !dryrun {
		err := os.MkdirAll(destPath, os.FileMode(int(0755)))
		if err != nil {
			log.Fatal(err)
		}
		if p.path != destPath {
			err = os.Rename(p.path, path.Join(destPath, p.fileName()))
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
