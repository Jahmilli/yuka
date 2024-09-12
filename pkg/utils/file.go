package utils

import (
	"bufio"
	"os"
	"path/filepath"

	"go.uber.org/zap"
)

// WriteToFile overwrite the contents of a file
func WriteToFile(logger *zap.SugaredLogger, s, file string, filePermissions int) error {
	filePath := filepath.Dir(file)
	if err := CreateDirectoriesForPath(logger, filePath); err != nil {
		return err
	}

	// overwrite the existing file contents
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.FileMode(filePermissions))
	if err != nil {
		logger.Warnf("Unable to open the file %s to write to: %v", file, err)
		return err
	}

	defer func(f *os.File) {
		err = f.Close()
		if err != nil {
			logger.Warnf("Unable to write to file [ %s ] %v", file, err)
		}
	}(f)

	wr := bufio.NewWriter(f)
	_, err = wr.WriteString(s)
	if err != nil {
		logger.Warnf("Unable to write to file [ %s ] %v", file, err)
	}
	if err = wr.Flush(); err != nil {
		logger.Warnf("Unable to write to file [ %s ] %v", file, err)
	}
	return nil
}

// FileExists returns true if the file path passed in exists otherwise false
func FileExists(f string) bool {
	if _, err := os.Stat(f); err != nil {
		return false
	}
	return true
}

// CreateDirectoriesForPath will create directories for a given file path if they do not already exist.
// If it exists, then the function will just return. An example would be for the file path, /etc/a/b/c,
// if 'a', 'b', 'c' are directories that do not exist, then this will create those respective directories with the correct structural hierarchy.
func CreateDirectoriesForPath(logger *zap.SugaredLogger, path string) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		logger.Debugf("path %s does not exist, will create directories for it", path)
		// if the directory path doesn't exist then we want to create the directories
		if err = os.MkdirAll(path, os.ModePerm); err != nil { // os.ModePerm gives read-write permissions
			logger.Debugf("Error creating directory: %v", err)
			return err
		}
		return nil
	}
	if os.IsExist(err) {
		// if it already exists then we're happy and can return
		logger.Debugf("Directory already exists: %s", path)
		return nil
	}
	return err
}
