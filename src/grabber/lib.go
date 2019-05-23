package grabber

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

type host struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Dtype   string `json:"dtype"`
	Method  string `json:"method"`
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

func filterStrings(s []string, f func(string) bool) []string {
	ret := make([]string, 0)
	for _, e := range s {
		if f(e) {
			ret = append(ret, e)
		}
	}
	return ret
}

func dirlistToFilenames(dl []os.FileInfo) []string {
	ret := make([]string, 0, len(dl))
	for _, f := range dl {
		ret = append(ret, f.Name())
	}
	return ret
}
