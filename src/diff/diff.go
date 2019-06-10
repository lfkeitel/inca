package diff

import (
	"bytes"
	"io"
	"os"
)

// SameFileContents read file1 and file2 in chuncks and returns
// if the files have the same contents. Size is checked first.
func SameFileContents(file1, file2 string) (bool, error) {
	// First check size, this should catch the majority of differences
	f1stat, err := os.Stat(file1)
	if err != nil {
		return false, err
	}

	f2stat, err := os.Stat(file2)
	if err != nil {
		return false, err
	}

	if f1stat.Size() != f2stat.Size() {
		return false, nil
	}

	// Then check contents
	fh1, err := os.Open(file1)
	if err != nil {
		return false, err
	}
	defer fh1.Close()

	fh2, err := os.Open(file2)
	if err != nil {
		return false, err
	}
	defer fh2.Close()

	// Read file in small chunks to save space and bail early
	f1buff := make([]byte, 4096)
	f2buff := make([]byte, 4096)
	for {
		f1rn, f1err := fh1.Read(f1buff)
		if f1err != nil && f1err != io.EOF {
			return false, f1err
		}

		f2rn, f2err := fh2.Read(f2buff)
		if f2err != nil && f2err != io.EOF {
			return false, f2err
		}

		// Reads weren't equal, this shouldn't happen because
		// we check the size first, but you never know.
		if f1rn != f2rn {
			return false, nil
		}

		// If one is EOF, the other must be as well
		if f1err == io.EOF || f2err == io.EOF {
			break
		}

		// Finally, actually check the contents
		if !bytes.Equal(f1buff, f2buff) {
			return false, nil
		}
	}

	return true, nil
}
