package services

import (
	"crypto/md5"
	"encoding/hex"
	"nes-emu/src/shared/application"
	"os"
	"path"
	"path/filepath"
)

type localFileSystem struct {
}

func NewLocalFileSystem() application.FileSystem {
	return &localFileSystem{}
}

func (fs *localFileSystem) SaveFile(filepath string, data []byte) error {
	dir := path.Dir(filepath)
	err := os.MkdirAll(dir, os.ModePerm)

	if err != nil {
		return err
	}

	return os.WriteFile(filepath, data, 0644)
}

func (fs *localFileSystem) ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func (fs *localFileSystem) ExistsFile(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func (fs *localFileSystem) CopyFile(src string, dst string) error {
	input, err := os.ReadFile(src)

	if err != nil {
		return err
	}

	dir := path.Dir(dst)
	err = os.MkdirAll(dir, os.ModePerm)

	if err != nil {
		return err
	}

	return os.WriteFile(dst, input, 0644)
}

func (fs *localFileSystem) RemoveFile(path string) error {
	return os.Remove(path)
}

func (fs *localFileSystem) ListFiles(path string) ([]string, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var fileNames []string
	for _, file := range files {
		if !file.IsDir() {
			fileNames = append(fileNames, filepath.Join(path, file.Name()))
		}
	}

	return fileNames, nil
}

func (fs *localFileSystem) GetFileName(path string) string {
	filename := filepath.Base(path)
	extension := filepath.Ext(filename)
	return filename[:len(filename)-len(extension)]
}

func (fs *localFileSystem) GetFileExtension(path string) string {
	return filepath.Ext(path)
}

func (fs *localFileSystem) GetFileMD5Hash(path string) (string, error) {
	data, err := fs.ReadFile(path)
	if err != nil {
		return "", err
	}
	if len(data) == 0 {
		return "", nil
	}
	hash := md5.Sum(data)
	return hex.EncodeToString(hash[:]), nil
}
