package main

import (
	"html/template"
	"os"
	"strings"
)

type pageData struct {
	Links   []Link
	WebPath string
}

func WritePage(path string, files []os.FileInfo) string {
	pageHtml := `<!DOCTYPE html>
<html>
  <head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width,initial-scale=1">
    <title>LANshare: {{.WebPath}}</title>
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
      <h1>LANshare: {{.WebPath}}</h1>
      <ul>
        {{range $val := .Links}}
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
