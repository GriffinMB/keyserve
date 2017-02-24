package main

import (
    	"crypto/tls"
	"golang.org/x/crypto/acme/autocert"
	"flag"
	"net/http"
	"io/ioutil"
	"html/template"
	"github.com/russross/blackfriday"
)

var path string = "/keybase/public/griffinmb/blog"

func handler(w http.ResponseWriter, r *http.Request) {
    	reqPath := r.URL.Path

	if (reqPath == "/") {
		reqPath = "/index"
	}
    	
    	data, err := ioutil.ReadFile(path + reqPath + ".md")
    	if (err != nil) {
        	http.NotFound(w, r)
        	return
    	}
    	output := blackfriday.MarkdownBasic(data)
    	tpl := `
	<!doctype html>
	<html>
	    <head>
	        <meta charset="UTF-8">
	        <title>Griffin's Blog</title>
	        <link rel="stylesheet" type="text/css" href="/static/style.css">
	    </head>
	    <body>
	        {{.}}
	    </body>
	</html>
    	`
    	t, err := template.New("layout").Parse(tpl)
    	t.Execute(w, template.HTML(string(output)))
}

func main() {
    	prod := flag.Bool("prod", false, "run in production mode")
    	flag.Parse()
  	
	fs := http.StripPrefix("/static/", http.FileServer(http.Dir(path + "/static")))
	mux := http.NewServeMux()

	mux.Handle("/static/", fs)
	mux.HandleFunc("/", handler)

	if (*prod) {
		certManager := autocert.Manager{
    		        Prompt:     autocert.AcceptTOS,
    			HostPolicy: autocert.HostWhitelist("www.griffinbyatt.com"),
    		        Cache:      autocert.DirCache("/certs"),
		}

	    	server := &http.Server{
			Addr: ":443",
			Handler: mux,
    	                TLSConfig: &tls.Config{
        			GetCertificate: certManager.GetCertificate,
    	                },
	    	}

	        server.ListenAndServeTLS("", "") 
	} else {
		server := &http.Server{
			Addr:    "0.0.0.0:8080",
			Handler: mux,
		}
		server.ListenAndServe()
	}
}
