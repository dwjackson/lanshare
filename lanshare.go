package main

import (
	"archive/zip"
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
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

	path := "."

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
		fileContent, fileErr := readFile(fileName)
		if fileErr != nil {
			log.Fatal(fileErr)
		}
		_, err := res.Write(fileContent)
		if err != nil {
			log.Fatal(err)
		}
	})

	http.HandleFunc("/download_all", func(res http.ResponseWriter, req *http.Request) {
		log.Printf("Download All\n")
		path = "."
		buf := new(bytes.Buffer)
		zipWriter := zip.NewWriter(buf)
		files, err := readDir(path)
		if err != nil {
			log.Fatal(err)
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		for _, fi := range files {
			if fi.Name()[0] == '.' || fi.IsDir() {
				continue
			}
			fileName := fi.Name()
			zipEntry, err := zipWriter.Create(fileName)
			if err != nil {
				log.Fatal(err)
				res.WriteHeader(http.StatusInternalServerError)
				return 
			}
			fileBytes, err := os.ReadFile(fileName)
			if err != nil {
				log.Fatal(err)
				res.WriteHeader(http.StatusInternalServerError)
				return 
			}
			_, err = zipEntry.Write(fileBytes)
			if err != nil {
				log.Fatal(err)
				res.WriteHeader(http.StatusInternalServerError)
				return 
			}
		}
		err = zipWriter.Close()
		if err != nil {
			log.Fatal(err)
			res.WriteHeader(http.StatusInternalServerError)
			return 
		}
		res.Header().Set("Content-Disposition", fmt.Sprintf("attachment;filename=%s", "all_files.zip"))
		res.Write(buf.Bytes())
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

func upDir(path string) Link {
	var href string
	if path == "." {
		href = "/"
	} else {
		pathParts := strings.Split(path, "/")
		pathParts = pathParts[:len(pathParts)-1]
		href = "/" + strings.Join(pathParts, "/")
	}
	return Link{
		Name: "..",
		Href: href,
		Size: 0,
	}
}

func readFile(fileName string) ([]byte, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	stat, statErr := file.Stat()
	if statErr != nil {
		return nil, statErr
	}

	size := stat.Size()
	buf := make([]byte, size)
	r := bufio.NewReader(file)
	_, readErr := r.Read(buf)

	return buf, readErr
}
