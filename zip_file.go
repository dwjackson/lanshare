package main

import (
	"archive/zip"
	"os"
)

type TempZipFile struct {
	file *os.File
}

func (t *TempZipFile) Close() error {
	err := t.file.Close()
	if err != nil {
		return err
	}
	err = os.Remove(t.file.Name())
	if err != nil {
		return err
	}
	return nil
}

func CreateTemporaryZipFile(path string) (TempZipFile, error) {
	file, err := os.CreateTemp("", "all_files")
	if err != nil {
		return TempZipFile{}, err
	}

	zipWriter := zip.NewWriter(file)
	files, err := readDir(path)
	if err != nil {
		return TempZipFile{}, err
	}

	for _, fi := range files {
		if fi.Name()[0] == '.' || fi.IsDir() {
			continue
		}
		fileName := fi.Name()
		zipEntry, err := zipWriter.Create(fileName)
		if err != nil {
			return TempZipFile{}, err
		}
		filePath := path + "/" + fileName
		fileBytes, err := os.ReadFile(filePath)
		if err != nil {
			return TempZipFile{}, err
		}
		_, err = zipEntry.Write(fileBytes)
		if err != nil {
			return TempZipFile{}, err
		}
	}

	err = zipWriter.Close()
	if err != nil {
		return TempZipFile{}, err
	}

	return TempZipFile{file}, nil
}
