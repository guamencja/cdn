package main

import (
	"fmt"
	"log"
	"strings"
	"net/http"
	"io/ioutil"

	"github.com/gorilla/mux"
)

func UploadEndpoint(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(100 << 20) // 100 mb
	file, handler, err := r.FormFile("file")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer file.Close()

	buff := make([]byte, handler.Size) // btw. przy sprawdzaniu tego, prawie mi kakuter jebnal, bo dalem 100 << 20 size arraya XDDDDDDDDDDDDDDDDDdd
	file.Read(buff)
	
	if err := ioutil.WriteFile(fmt.Sprintf("./files/%s", handler.Filename), buff, 0644); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// nie mam innego pomyslu jak wylaczyc directory listing, ok?
func h(next http.Handler) (http.Handler) {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		w.Header().Set("Cache-Control", "public, max-age=31536000")
		next.ServeHTTP(w, r)
	})
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/upload", UploadEndpoint)
	r.PathPrefix("/files/").Handler(h(http.StripPrefix("/files/", http.FileServer(http.Dir("./files/")))))

	server := &http.Server{
		Handler: r,
		Addr: fmt.Sprintf(":%s", "8080"), // kiedys napisze config
	}

	log.Fatal(server.ListenAndServe())
}