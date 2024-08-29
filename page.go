package main

import (
	_ "embed"
	"html/template"
	"os"
	"strings"
)

type pageData struct {
	Links   []Link
	WebPath string
}

//go:embed index.html
var pageHtml string

func WritePage(path string, files []os.FileInfo) string {
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
	var webPath string
	if path == "." {
		webPath = "/"
	} else if strings.HasPrefix(path, ".") {
		webPath = path[1:]
	} else {
		webPath = path
	}
	data := pageData{
		Links:   links,
		WebPath: webPath,
	}
	err := t.Execute(&b, data)
	if err != nil {
		panic(err)
	}
	return b.String()
}
