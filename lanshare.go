package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

const downloadPath string = "/download/"
const DEFAULT_MAX_FILE_SIZE = "10MiB"

func readDir(path string) ([]os.FileInfo, error) {
	dir, err := os.Open(path)
	if err != nil {
		return nil, errors.New("Could not open directory: " + path)
	}
	files, err := dir.Readdir(0)
	if err != nil {
		return nil, errors.New("Could not read current directory")
	}
	return files, nil
}

func currentPath(req *http.Request) string {
	var path string
	reqPath := strings.Replace(req.URL.Path, "..", "", -1)
	reqPath = strings.Replace(reqPath, "//", "/", -1)
	if reqPath == "/" {
		path = "."
	} else {
		path = "." + reqPath
	}
	return path
}

func main() {
	portPtr := flag.Int("p", 8080, "Port")
	maxUploadFileSizeString := flag.String("m", DEFAULT_MAX_FILE_SIZE, "Max file size")
	flag.Parse()
	maxUploadFileSize, err := parseFileSize(*maxUploadFileSizeString)
	if err != nil {
		panic(err)
	}
	fmt.Println("Max upload file size in bytes: ", maxUploadFileSize)

	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			http.Error(res, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		log.Printf("%s %s\n", req.Method, req.URL)
		path := currentPath(req)
		files, err := readDir(path)
		if err != nil {
			log.Fatal(err)
		}
		pageHtml := WritePage(path, files)
		_, err = res.Write([]byte(pageHtml))
		if err != nil {
			log.Fatal(err)
		}
	})

	http.HandleFunc("/favicon.ico", func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			http.Error(res, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		http.Error(res, "No such file", http.StatusNotFound)
	})

	http.HandleFunc(downloadPath, func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			http.Error(res, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		log.Printf("%s %s", req.Method, req.URL)
		fileName := strings.TrimPrefix(req.URL.Path, downloadPath)

		file, fileError := os.Open(fileName)
		if fileError != nil {
			log.Fatal(fileError)
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer file.Close()

		serveFile(res, req, file, file.Name())
	})

	http.HandleFunc("/download_all", func(res http.ResponseWriter, req *http.Request) {
		log.Printf("Download All\n")

		reqParams, err := url.ParseQuery(req.URL.RawQuery)
		if err != nil {
			log.Fatal(err)
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		log.Println(reqParams)
		pathParam, ok := reqParams["path"]
		if !ok {
			log.Fatal("No path given")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		path := pathParam[0]
		path = strings.Replace(path, "..", "", -1) // Remove any "go up" directives
		path = "." + path

		zipFile, err := CreateTemporaryZipFile(path)
		if err != nil {
			log.Fatal(err)
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer zipFile.Close()

		serveFile(res, req, zipFile.file, "all_files.zip")
	})

	http.HandleFunc("/upload", func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			http.Error(res, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		err := req.ParseMultipartForm(maxUploadFileSize)
		if err != nil {
			log.Println("Error parsing form: ", err)
			http.Error(res, "Error parsing form", http.StatusInternalServerError)
			return
		}
		uploadedFile, header, err := req.FormFile("file")
		if err != nil {
			log.Println("Error reading file")
			http.Error(res, "Could not read file", http.StatusInternalServerError)
			return
		}
		defer uploadedFile.Close()

		fileName := header.Filename
		file, err := os.Create(fileName)
		if err != nil {
			log.Println("Error creating file: ", err)
			http.Error(res, "Could not save file", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		_, err = io.Copy(file, uploadedFile)
		if err != nil {
			log.Println("Error writing file: ", err)
			http.Error(res, "Error writing file", http.StatusInternalServerError)
			return
		}

		log.Printf("File uploaded: %s", fileName)
		res.WriteHeader(http.StatusOK)
	})

	addr := ":" + strconv.Itoa(*portPtr)
	fmt.Printf("listening on %s...\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func serveFile(res http.ResponseWriter, req *http.Request, file *os.File, fileName string) {
	fileStat, statError := file.Stat()
	if statError != nil {
		log.Fatal(statError)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	contentDisposition := fmt.Sprintf("attachment; fileName=%s", fileName)
	res.Header().Set("Content-Disposition", contentDisposition)
	res.Header().Set("Content-Type", "application/octet-stream")
	fileSize := fileStat.Size()
	res.Header().Set("Content-Length", strconv.FormatInt(fileSize, 10))

	http.ServeContent(res, req, fileName, fileStat.ModTime(), file)
}
