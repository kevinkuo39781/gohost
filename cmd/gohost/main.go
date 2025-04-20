package main

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	uploadDir     = "./data/sites"
	maxUploadSize = 20 << 20 // 20MB
)

func main() {
	// Create upload dir if not exists
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		log.Fatalf("Unable to create data directory: %v", err)
	}

	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/g/", serveHandler)
	http.HandleFunc("/", indexHandler)

	fmt.Println("GoHost running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./public/index.html")
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		http.Error(w, "File too large", http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "File upload failed", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	path := r.FormValue("path")
	if path == "" || strings.Contains(path, "..") {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}

	targetDir := filepath.Join(uploadDir, path)
	os.RemoveAll(targetDir)
	os.MkdirAll(targetDir, 0755)

	tmpZip := filepath.Join(uploadDir, path+".zip")
	out, err := os.Create(tmpZip)
	if err != nil {
		http.Error(w, "Unable to save file", http.StatusInternalServerError)
		return
	}
	defer out.Close()
	io.Copy(out, file)

	if err := unzip(tmpZip, targetDir); err != nil {
		http.Error(w, "Failed to unzip", http.StatusInternalServerError)
		return
	}
	os.Remove(tmpZip)
	fmt.Fprintf(w, "Uploaded and deployed to /g/%s\n", path)
}

func serveHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/g/")
	staticPath := filepath.Join(uploadDir, path)
	http.ServeFile(w, r, staticPath)
}

func unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		fpath := filepath.Join(dest, f.Name)
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", fpath)
		}
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}
		os.MkdirAll(filepath.Dir(fpath), os.ModePerm)
		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}
		rc, err := f.Open()
		if err != nil {
			return err
		}
		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()
		if err != nil {
			return err
		}
	}
	return nil
}
