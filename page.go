package main

import (
	_ "embed"
	"html/template"
	"os"
	"sort"
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
	for _, fi := range files {
		link := linkFromFileInfo(path, fi)
		links = append(links, link)
	}
	sort.Slice(links, func(i, j int) bool {
		return links[i].Name < links[j].Name
	})
	upSlice := []Link { upDir(path) }
	if path != "." {
		links = append(upSlice, links...)
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
