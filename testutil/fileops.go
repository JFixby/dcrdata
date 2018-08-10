package testutil

import (
	"os"
	"path/filepath"
)

// RemoveContents clears target folder content
func RemoveContents(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}

// DirChild returns full path to a child of a target directory
func DirChild(dir string, child string) string {
	target, err := filepath.Abs(dir)
	if err != nil {
		panic("Failed to list folder: " + target)
	}
	return filepath.Join(target, child)
}

// ListDir returns list of full paths for a target directory content
func ListDir(target string) []string {
	target, err := filepath.Abs(target)
	if err != nil {
		panic("Failed to list folder: " + target)
	}
	d, err := os.Open(target)
	if err != nil {
		panic("Failed to list folder: " + target)
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		panic("Failed to list folder: " + target)
	}

	var result = make([]string, len(names))
	for i, name := range names {
		file := filepath.Join(target, name)
		result[i] = file
	}

	return result
}
