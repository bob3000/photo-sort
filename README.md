# Photo Sort

Photo Sort reads the creation time of jpeg photos and puts them into subfolders
organized by year/month/day/hour/minute.

It works but I wrote it just for fun though. No guarantees. If you're looking
for a serious command line tool to sort your photos I recommend using
[exiftool](https://sourceforge.net/projects/exiftool/).

## Usage

```
Photo Sort - sorts your photos

Usage:
  photosort [options] (year|month|day|hour|minute) <inputdir>

Options:
  -h --help        Show this screen.
  -v, --verbose    Show whats being done
  -d, --dryrun     Simulate sorting
  -r, --recursive  Scan directories recursively
  -c, --cleanup    Remove empty source directories when run recursively
  --version        Show version.
```