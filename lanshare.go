package main

import (
	"os"
	"log"
	"strings"
	"html/template"
	"net/http"
	"fmt"
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
      <li>{{$name}}</li>
      {{end}}
    </ul>
  </body>
</html>`
	t := template.Must(template.New("pageHtml").Parse(pageHtml))
	b := strings.Builder{}
	t.Execute(&b, files)
	return b.String()
}
