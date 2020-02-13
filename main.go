package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"strings"
)

func uploadFile(w http.ResponseWriter, r *http.Request) {
	fmt.Println("---- File Upload Start ----")
	// upload of 10 MB files.
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		fmt.Fprintln(w, err)
		log.Fatal(err)
		return
	}

	file, header, err := r.FormFile("originalFile")
	if err != nil {
		log.Fatal("Error Retrieving the File")
		log.Fatal(err)
		return
	}
	defer file.Close()

	fmt.Printf("Uploaded File: %+v\n", header.Filename)
	fmt.Printf("File Size: %+v\n", header.Size)
	fmt.Printf("MIME Header: %+v\n", header.Header)

	// Create a temporary file within out tmp dir
	tempFile, err := ioutil.TempFile("tmp", strings.TrimSuffix(header.Filename, path.Ext(header.Filename))+"-*"+path.Ext(header.Filename))
	if err != nil {
		log.Fatal(err)
	}
	defer tempFile.Close()

	io.Copy(tempFile, file)
	defer tempFile.Close()
	fmt.Fprintf(w, "Successfully Uploaded File\n")
}

func runSever() {
	http.HandleFunc("/uploadFile", uploadFile)

	fmt.Printf("Start Server at %s\n", "localhost:8000")
	// server is listening on port
	http.ListenAndServe(":8000", nil)

}

func main() {
	runSever()
}
