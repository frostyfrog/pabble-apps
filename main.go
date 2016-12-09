package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.WriteHeader(404)
		w.Write([]byte("404 Not Found<br />Try <a href='/'>Home</a>?"))
		return
	}
	fmt.Fprintf(w, "<a href='/apps'>Apps</a>")
}

func RecurseFolder(w http.ResponseWriter, path string, f os.FileInfo, lvl int) {
	for i := 0; i < lvl; i++ {
		w.Write([]byte("="))
	}
	fmt.Fprintf(w, "> %s<br />", f.Name())
	if f.IsDir() {
		fpath := fmt.Sprintf("%s/%s", path, f.Name())
		folder, err := ioutil.ReadDir(fpath)
		if err != nil {
			log.Println(err)
			return
		}
		for _, f1 := range folder {
			RecurseFolder(w, fpath, f1, lvl+1)
		}
	}
}

func AppsHandler(w http.ResponseWriter, r *http.Request) {
	files, err := ioutil.ReadDir("PebbleAppStore/apps/0/0")
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		w.Header().Add("content-type", "text/html")
		//fmt.Fprintf(w, ">%s<br />", file.Name())
		RecurseFolder(w, "PebbleAppStore/apps/0/0", file, 0)
	}
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/apps", AppsHandler)
	loggedRouter := handlers.LoggingHandler(os.Stdout, r)
	http.Handle("/", r)
	http.ListenAndServe(":8080", loggedRouter)
}
