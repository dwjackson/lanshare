package main

import (
	"bufio"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"io/ioutil"
)

const downloadPath string = "/download/"
const templateHtmlFileName string = "index.html.tmpl"

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

func main() {
	path := "."

	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		log.Printf("%s %s", req.Method, req.URL)
		reqPath := strings.Replace(req.URL.Path, "..", "", -1)
		reqPath = strings.Replace(reqPath, "//", "/", -1)
		if reqPath == "/" {
			path = "."
		} else {
			path = "." + reqPath
		}
		files, err := readDir(path)
		if err != nil {
			log.Fatal(err)
		}
		pageHtml := writePage(path, files)
		_, err = res.Write([]byte(pageHtml))
		if err != nil {
			log.Fatal(err)
		}
	})

	http.HandleFunc("/favicon.ico", func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(404)
		_, err := res.Write([]byte{})
		if err != nil {
			log.Fatal(err)
		}
	})

	http.HandleFunc(downloadPath, func(res http.ResponseWriter, req *http.Request) {
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

	addr := ":8080"
	fmt.Printf("listening on %s...\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

type Link struct {
	Name       string
	Href       string
	IsDownload bool
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
		Href:       href,
		IsDownload: !fi.IsDir(),
	}
}

func writePage(path string, files []os.FileInfo) string {
	content, err := ioutil.ReadFile(templateHtmlFileName)
	if err != nil {
		panic("Could not read tempate HTML file")
	}
	pageHtml := string(content)

	t := template.Must(template.New("pageHtml").Parse(pageHtml))
	b := strings.Builder{}
	var links []Link = []Link{
		upDir(path),
	}
	for _, fi := range files {
		link := linkFromFileInfo(path, fi)
		links = append(links, link)
	}
	t.Execute(&b, links)
	return b.String()
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
