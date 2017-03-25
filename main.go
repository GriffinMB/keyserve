package main

import (
	"flag"
	"fmt"
	"math/rand"
	"github.com/russross/blackfriday"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

var path string
var cache = make(map[string][]byte)
var resetTime int64
var cacheBust int
var tpl string

type TemplateData struct {
	Data      template.HTML
	CacheBust int
}

func findOrCache(reqPath string) (data []byte, err error) {
	if data, ok := cache[reqPath]; ok {
		return data, nil
	}

	data, err = ioutil.ReadFile(path + reqPath + ".md")

	if err == nil {
		cache[reqPath] = data
	}

	return data, err
}

func handler(w http.ResponseWriter, r *http.Request) {
	reqPath := r.URL.Path

	if reqPath == "/" {
		reqPath = "/index"
	}

	data, err := findOrCache(reqPath)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	output := blackfriday.MarkdownBasic(data)
	tplData := TemplateData{Data: template.HTML(string(output)), CacheBust: cacheBust}
	t, err := template.New("layout").Parse(tpl)
	t.Execute(w, tplData)
}

func reset() {
	resetPath := path + "/reset.txt"
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for {
		info, err := os.Stat(resetPath)
		if os.IsNotExist(err) {
			continue
		}

		modTime := info.ModTime().Unix()
		if resetTime == 0 || modTime > resetTime {
			resetTime = modTime
			cacheBust = r.Int()
			cache = make(map[string][]byte)
		}
		time.Sleep(5 * time.Second)
	}
}

func main() {
	port := flag.String("port", "8080", "specify port")
	uname := flag.String("uname", "", "Keybase username")
	flag.Parse()

	if *uname == "" {
		fmt.Println("Must specify a Keybase username.")
		os.Exit(1)
	}

	path = "/keybase/public/" + *uname + "/blog"

	go reset()

	tpl = `
	<!doctype html>
	<html>
	    <head>
	        <meta charset="UTF-8">
	        <title>Griffin's Blog</title>
	        <link rel="stylesheet" type="text/css" href="/static/style.css?{{.CacheBust}}">
	    </head>
	    <body>
	    	<div class="main">
	            {{.Data}}
	        </div>
	    </body>
	</html>
    	`

	fs := http.StripPrefix("/static/", http.FileServer(http.Dir(path+"/static")))
	mux := http.NewServeMux()

	mux.Handle("/static/", fs)
	mux.HandleFunc("/", handler)

	server := &http.Server{
		Addr:    "0.0.0.0:" + *port,
		Handler: mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
	server.ListenAndServe()

}
