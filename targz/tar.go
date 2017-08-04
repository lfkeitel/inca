package tarGz

/**
 * The code for this package was taken from MadCrazy's question on StackOverflow
 * URL: http://stackoverflow.com/questions/13611100/how-to-write-a-directory-not-just-the-files-in-it-to-a-tar-gz-file-in-golang
 */

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"

	"github.com/lfkeitel/verbose"
)

var appLogger *verbose.Logger

func init() {
	appLogger = verbose.New("tarGz-log")

	fileLogger, err := verbose.NewFileHandler("logs/tar.log")
	if err != nil {
		panic("Failed to open logging directory")
	}

	appLogger.AddHandler("file", fileLogger)
}

func handleError(err error) {
	if err != nil {
		appLogger.Error(err.Error())
	}
}

func tarGzWrite(_path string, tw *tar.Writer, fi os.FileInfo) {
	fr, err := os.Open(_path)
	handleError(err)
	defer fr.Close()

	h := new(tar.Header)
	h.Name = _path
	h.Size = fi.Size()
	h.Mode = int64(fi.Mode())
	h.ModTime = fi.ModTime()

	err = tw.WriteHeader(h)
	handleError(err)

	_, err = io.Copy(tw, fr)
	handleError(err)
	return
}

func iterDirectory(dirPath string, tw *tar.Writer) {
	dir, err := os.Open(dirPath)
	handleError(err)
	defer dir.Close()
	fis, err := dir.Readdir(0)
	handleError(err)
	for _, fi := range fis {
		curPath := dirPath + "/" + fi.Name()
		if fi.IsDir() {
			iterDirectory(curPath, tw)
		} else {
			tarGzWrite(curPath, tw, fi)
		}
	}
	return
}

func TarGz(outFilePath string, inPath string) {
	// file write
	fw, err := os.Create(outFilePath)
	handleError(err)
	defer fw.Close()

	// gzip write
	gw := gzip.NewWriter(fw)
	defer gw.Close()

	// tar write
	tw := tar.NewWriter(gw)
	defer tw.Close()

	iterDirectory(inPath, tw)
	return
}