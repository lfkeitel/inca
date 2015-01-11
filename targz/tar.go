package tarGz

/**
 * The code for this package was taken from MadCrazy's question on StackOverflow
 * URL: http://stackoverflow.com/questions/13611100/how-to-write-a-directory-not-just-the-files-in-it-to-a-tar-gz-file-in-golang
 */

import (
    "os"
    "io"
    "archive/tar"
    "compress/gzip"

    logger "github.com/dragonrider23/go-logger"
)

var appLogger *logger.Logger

func init() {
    appLogger = logger.New("tarGz-log")
}

func handleError(_e error) {
    if _e != nil {
        appLogger.Error(_e.Error())
    }
}

func TarGzWrite(_path string, tw *tar.Writer, fi os.FileInfo) {
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

func IterDirectory(dirPath string, tw *tar.Writer) {
    dir, err := os.Open(dirPath)
    handleError(err)
    defer dir.Close()
    fis, err := dir.Readdir(0)
    handleError(err)
    for _, fi := range fis {
        curPath := dirPath + "/" + fi.Name()
        if fi.IsDir() {
            //TarGzWrite(curPath, tw, fi)
            IterDirectory(curPath, tw)
        } else {
            appLogger.Info("adding... %s", curPath)
            TarGzWrite(curPath, tw, fi)
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

    IterDirectory(inPath, tw)
    return
}
