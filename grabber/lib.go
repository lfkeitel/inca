package grabber

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

type host struct {
	name    string
	address string
	dtype   string
	method  string
}

type dtype struct {
	deviceType string
	method     string
	scriptfile string
	args       string
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

		os.RemoveAll(filepath.Join(path, file.Name()))
	}

	return nil
}
