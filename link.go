package main

import (
	"os"
	"strconv"
)

type Link struct {
	Name       string
	Size       int64
	HumanSize  string
	Href       string
	IsDownload bool
}

const kilobyte int64 = 1024
const megabyte int64 = kilobyte * kilobyte
const gigabyte int64 = megabyte * kilobyte

func humanSize(size int64) string {
	if size < kilobyte {
		return strconv.FormatInt(size, 10) + "B"
	}
	if size >= kilobyte && size < megabyte {
		kbSize := float64(size) / float64(kilobyte)
		return strconv.FormatFloat(kbSize, 'f', 2, 64) + "kB"
	}
	if size >= megabyte && size < gigabyte {
		mbSize := float64(size) / float64(megabyte)
		return strconv.FormatFloat(mbSize, 'f', 2, 64) + "MB"
	}
	return ""
}

func linkFromFileInfo(path string, fi os.FileInfo) Link {
	fileName := fi.Name()
	filePath := path + "/" + fileName
	var href string
	if fi.IsDir() {
		href = filePath
	} else {
		href = downloadPath + filePath
	}
	return Link{
		Name:       fileName,
		Size:       fi.Size(),
		HumanSize:  humanSize(fi.Size()),
		Href:       href,
		IsDownload: !fi.IsDir(),
	}
}
