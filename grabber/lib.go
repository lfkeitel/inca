package grabber

import (
    "os"
    "time"
)

type host struct {
    name    string
    address string
    manufacturer   string
    proto   string
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
