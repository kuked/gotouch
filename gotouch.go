package gotouch

import (
	"os"
	"time"

	"github.com/djherbis/atime"
)

// Create creates new empty file.
func Create(name string) error {
	f, err := os.Create(name)
	if err != nil {
		return err
	}

	if err := f.Close(); err != nil {
		return err
	}
	return nil
}

// UpdateAtime updates access time of the file. It leaves modification time as it is.
func UpdateAtime(name string, atime time.Time) error {
	fi, err := os.Stat(name)
	if err != nil {
		return err
	}

    mtime := fi.ModTime()
	if err := os.Chtimes(name, atime, mtime); err != nil {
		return err
	}
	return nil
}

// UpdateMtime updates modification time of the file. It leaves access time as it is.
func UpdateMtime(name string, mtime time.Time) error {
	atime, err := atime.Stat(name)
	if err != nil {
		return err
	}

	if err := os.Chtimes(name, atime, mtime); err != nil {
		return err
	}
	return nil
}

// UpdateTime updates access time and modification time of the file.
func UpdateTime(name string, amtime time.Time) error {
	if err := os.Chtimes(name, amtime, amtime); err != nil {
		return err
	}
	return nil
}
