package main

import (
	"bufio"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"errors"
)

const downloadPath string = "/download/"

func readDir(path string) ([]os.FileInfo, error) {
	dir, err := os.Open(path)
	if err != nil {
		return nil, errors.New("Could not open directory")
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
		reqPath := req.URL.Path
		if reqPath != "/" {
			if path == "/" {
				reqPath = strings.TrimLeft(reqPath, "/")
			}
			path += reqPath
		} else {
			path = "."
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
	Name string
	Href string
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
		Name: fileName,
		Href: href,
		IsDownload: !fi.IsDir(),
	}
}

func writePage(path string, files []os.FileInfo) string {
	pageHtml := `<!DOCTYPE html>
<html>
  <head>
    <meta charset="utf-8">
    <title>LANshare</title>
  </head>
  <body>
    <h1>LANshare</h1>
    <ul>
      {{range $val := .}}
      <li><a href="{{$val.Href}}" {{if $val.IsDownload}}download{{end}}>{{$val.Name}}</a></li>
      {{end}}
    </ul>
  </body>
</html>`
	t := template.Must(template.New("pageHtml").Parse(pageHtml))
	b := strings.Builder{}
	var links []Link = []Link{
		Link{
			Name: "..",
			Href: "/",
		},
	}
	for _, fi := range(files) {
		link := linkFromFileInfo(path, fi)
		links = append(links, link)
	}
	t.Execute(&b, links)
	return b.String()
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
