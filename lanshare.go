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
)

const downloadPath string = "/download/"

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

func writePage(path string, files []os.FileInfo) string {
	pageHtml := `<!DOCTYPE html>
<html>
  <head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width,initial-scale=1">
    <title>LANshare</title>
    <style>
      body {
        font-family: 'Helvetica', 'Arial', sans-serif;
      }
      #content {
        max-width: 980px;
	margin: auto;
      }
    </style>
  </head>
  <body>
    <div id="content">
      <h1>LANshare</h1>
      <ul>
        {{range $val := .}}
        <li><a href="{{$val.Href}}" {{if $val.IsDownload}}download{{end}}>{{$val.Name}}</a>{{ if $val.IsDownload }} ({{$val.HumanSize}}){{end}}</li>
        {{end}}
      </ul>
    </div>
  </body>
</html>`

	t := template.Must(template.New("pageHtml").Parse(pageHtml))
	b := strings.Builder{}
	var links []Link
	if path != "." {
		links = append(links, upDir(path))
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
