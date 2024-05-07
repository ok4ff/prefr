package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Flags struct {
	prefix    string
	directory string
}

var inFlags Flags

func init() {
	flag.StringVar(&inFlags.directory, "d", "", "directory with files")
	flag.StringVar(&inFlags.prefix, "p", "", "prefix for remove")
}

func main() {
	flag.Parse()
	if inFlags.directory == "" {
		fmt.Println("needed set flag -p (prefix for remote)")
		os.Exit(1)
	}

	if inFlags.prefix == "" {
		fmt.Println("needed set flag -p (prefix for remote)")
		os.Exit(1)
	}

	if _, err := os.Stat(inFlags.directory); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	files := make(map[int][]string)
	level := 0

	if err := filesForRename(inFlags.directory, files, level); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	keys := make([]int, 0, len(files))
	for k := range files {
		keys = append(keys, k)
	}

	sort.Ints(keys)

	for i := len(keys) - 1; i > -1; i-- {
		files := files[i]
		for _, path := range files {
			directory, filename := filepath.Split(path)
			filename = strings.ReplaceAll(filename, inFlags.prefix, "")
			newName := filepath.Join(directory, filename)
			if i == 2 {
				fmt.Println(i, path, "->", filepath.Join(directory, filename))
			}
			if path != newName {
				if err := os.Rename(path, newName); err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
			}
		}
	}

	os.Exit(0)
}

func filesForRename(directory string, files map[int][]string, level int) error {
	ents, err := os.ReadDir(directory)
	if err != nil {
		return err
	}

	if len(files[level]) == 0 {
		files[level] = make([]string, 0, len(ents))
	} else {
		old := files[level]
		files[level] = make([]string, 0, len(old)+len(ents))
		files[level] = append(files[level], old...)
	}

	for _, e := range ents {
		filename := filepath.Join(directory, e.Name())
		files[level] = append(files[level], filename)
		if e.IsDir() {
			if err := filesForRename(filename, files, level+1); err != nil {
				return nil
			}
		}
	}

	return nil
}
