package project

import (
	"github.com/docker/docker/builder/dockerignore"
	"github.com/docker/docker/pkg/fileutils"
	"io/ioutil"
	"os"
	"path/filepath"
)

func ContextFiles(path string) ([]string, error) {
	f, err := os.Open(filepath.Join(path, ".dockerignore"))
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	defer f.Close()

	var excludes []string
	if err == nil {
		excludes, err = dockerignore.ReadAll(f)
		if err != nil {
			return nil, err
		}
		// Include the '.dockerignore' file:
		excludes = append(excludes, ".dockerignore")
	}

	return RecurseDir(path, excludes)
}

func RecurseDir(path string, excludes []string) ([]string, error) {
	filenames := []string{}

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		prefixedPath := filepath.Join(path, file.Name())

		didMatch, err := fileutils.Matches(prefixedPath, excludes)
		if err != nil {
			return nil, err
		}

		if !didMatch {
			if file.IsDir() {
				childFilenames, err := RecurseDir(prefixedPath, excludes)
				if err != nil {
					return nil, err
				}
				filenames = append(filenames, childFilenames...)
			} else {
				filenames = append(filenames, prefixedPath)
			}
		}
	}

	return filenames, nil
}
