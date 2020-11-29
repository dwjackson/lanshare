package main

import (
	"bufio"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	dir, err := os.Open(".")
	if err != nil {
		log.Fatal("Could not open current directory")
		return
	}
	files, err := dir.Readdir(0)
	if err != nil {
		log.Fatal("Could not read current directory")
		return
	}
	pageHtml := writePage(files)

	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		_, err = res.Write([]byte(pageHtml))
		if err != nil {
			log.Fatal(err)
		}
	})

	http.HandleFunc("/file/", func(res http.ResponseWriter, req *http.Request) {
		fileName := strings.TrimPrefix(req.URL.Path, "/file/")
		fileContent, fileErr := readFile(fileName)
		if fileErr != nil {
			log.Fatal(err)
		}
		_, err = res.Write(fileContent)
		if err != nil {
			log.Fatal(err)
		}
	})

	addr := ":8080"
	fmt.Printf("listening on %s...\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func writePage(files []os.FileInfo) string {
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
      {{$name := $val.Name}}
      <li><a href="/file/{{$name}}">{{$name}}</a></li>
      {{end}}
    </ul>
  </body>
</html>`
	t := template.Must(template.New("pageHtml").Parse(pageHtml))
	b := strings.Builder{}
	t.Execute(&b, files)
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
