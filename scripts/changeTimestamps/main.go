package changeTimestamps

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func main() {
	endPath := flag.String("path", "", "directory path to change timestamps within")
	dateStr := flag.String("date", "03-21-2026", "date to use for all files")
	flag.Parse()
	if endPath == nil || *endPath == "" {
		panic("path is required")
	}
	dateStrs := strings.Split(*dateStr, "-")
	if len(dateStrs) != 3 {
		panic("date format is invalid. Must be MM-DD-YYYY")
	}
	mo, err := strconv.ParseInt(dateStrs[0], 10, 8)
	if err != nil {
		panic("month must be an integer")
	}
	dy, err := strconv.ParseInt(dateStrs[1], 10, 8)
	if err != nil {
		panic("day must be an integer")
	}
	yr, err := strconv.ParseInt(dateStrs[2], 10, 64)
	if err != nil {
		panic("year must be an integer")
	}
	modificationTime := time.Date(int(yr), time.Month(mo), int(dy), 0, 0, 0, 0, time.UTC)

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	root := wd + "/" + *endPath // TODO: ensure ok
	err = filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Printf("Error accessing path %s: %v\n", path, err)
			return err
		}

		if d.IsDir() {
			return nil
		}
		// Process the file. Change the access and modification times
		return os.Chtimes(path, modificationTime, modificationTime) // TODO: ensure this also changes creation time
	})
	if err != nil {
		log.Fatalf("Error walking the path %s: %v\n", root, err)
	}

}
