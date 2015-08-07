package grabber

import (
	"errors"
	"io/ioutil"
	"os"
	"time"
)

type device struct {
	deviceid int
	name     string
	address  string
	// Other information
}

func touch(filename string) error {
	os.Chtimes(filename, time.Now(), time.Now())
	file, err := os.OpenFile(filename, os.O_CREATE, 0664)
	if err != nil {
		return err
	}
	file.Close()
	return nil
}

func removeDir(path string) error {
	src, err := os.Stat(path)
	if err != nil {
		return err
	}
	if !src.IsDir() {
		return errors.New("Path is not a directory")
	}

	fileList, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}

	for _, file := range fileList {
		if file.Name()[0] == '.' {
			continue
		}

		os.Remove(path + "/" + file.Name())
	}

	return nil
}

func stringInSlice(needle string, haystack []string) bool {
	for _, item := range haystack {
		if item == needle {
			return true
		}
	}
	return false
}

func intInSlice(needle int, haystack []int) bool {
	for _, item := range haystack {
		if item == needle {
			return true
		}
	}
	return false
}
