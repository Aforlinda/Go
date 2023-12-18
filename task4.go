package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
)

// File represents information about an uploaded file.
type File struct {
	Name string
	URL  string
}

var (
	files     []File
	filesLock sync.Mutex
)

const (
	uploadFormHTML = `
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta http-equiv="X-UA-Compatible" content="IE=edge">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>File Upload</title>
		<link rel="stylesheet" type="text/css" href="/style.css">
	</head>
	<body>
		<h1>File Upload</h1>
		<form enctype="multipart/form-data" action="/upload" method="post">
			<label for="uploadFile">Choose a file:</label>
			<input type="file" name="uploadFile" id="uploadFile" required>
			<input type="submit" value="Upload">
		</form>
		<h2>Uploaded Files</h2>
		<ul>
			{{range .Files}}
				<li>
					<a href="{{.URL}}" target="_blank">{{.Name}}</a>
					<form action="/delete" method="post" style="display: inline;">
						<input type="hidden" name="file" value="{{.Name}}">
						<input type="submit" value="Delete">
					</form>
				</li>
			{{end}}
		</ul>
	</body>
	</html>
	`

	cssStyle = `
	body {
		background: linear-gradient(to right, #800080, #4B0082);
		color: #ffffff;
		font-family: Arial, sans-serif;
		margin: 0;
		padding: 0;
	}

	h1, h2 {
		text-align: center;
		margin: 20px 0;
	}

	form {
		max-width: 600px;
		margin: auto;
		padding: 20px;
		background-color: rgba(255, 255, 255, 0.8);
		border-radius: 10px;
		box-shadow: 0 0 10px rgba(0, 0, 0, 0.2);
	}

	input[type="file"] {
		width: 100%;
		padding: 10px;
		margin: 8px 0;
		display: inline-block;
		border: 1px solid #ccc;
		box-sizing: border-box;
		border-radius: 5px;
	}

	input[type="submit"] {
		background-color: #4CAF50;
		color: white;
		padding: 10px 15px;
		border: none;
		border-radius: 5px;
		cursor: pointer;
		font-size: 16px;
	}

	input[type="submit"]:hover {
		background-color: #45a049;
	}

	label {
		display: block;
		margin: 10px 0;
	}

	ul {
		list-style-type: none;
		padding: 0;
	}

	li {
		margin-bottom: 10px;
	}

	a {
		color: #ffffff;
		text-decoration: none;
	}
	`
)

func main() {
	// Set up HTTP handlers
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/download", downloadHandler)
	http.HandleFunc("/delete", deleteHandler)
	http.HandleFunc("/style.css", styleHandler)

	// Start the server in a goroutine
	go func() {
		port := 8080
		fmt.Printf("Server listening on :%d...\n", port)
		err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
		if err != nil {
			fmt.Println("Error:", err)
		}
	}()

	// Graceful shutdown
	waitForShutdown()
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	filesLock.Lock()
	defer filesLock.Unlock()

	// Execute the HTML template with the list of uploaded files
	tmpl, err := template.New("index").Parse(uploadFormHTML)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	tmplData := struct {
		Files []File
	}{
		Files: files,
	}

	err = tmpl.Execute(w, tmplData)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the form data to retrieve the uploaded file
	file, header, err := r.FormFile("uploadFile")
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Generate a unique filename for the uploaded file
	filename := header.Filename
	filename = filepath.Join("uploads", filename)

	// Create a new file in the "uploads" directory
	dst, err := os.Create(filename)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Copy the uploaded file to the destination file on the server
	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Update the list of uploaded files
	filesLock.Lock()
	defer filesLock.Unlock()

	files = append(files, File{
		Name: header.Filename,
		URL:  "/download?file=" + header.Filename,
	})

	// Redirect back to the index page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	fileName := r.FormValue("file")
	if fileName == "" {
		http.Error(w, "Bad Request: Please provide a filename", http.StatusBadRequest)
		return
	}

	// Open the requested file
	filePath := filepath.Join("uploads", fileName)
	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "File Not Found", http.StatusNotFound)
		return
	}
	defer file.Close()

	// Set the Content-Disposition header to trigger a download prompt
	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	w.Header().Set("Content-Type", "application/octet-stream")

	// Copy the file to the response writer
	_, err = io.Copy(w, file)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	fileName := r.FormValue("file")
	if fileName == "" {
		http.Error(w, "Bad Request: Please provide a filename", http.StatusBadRequest)
		return
	}

	// Delete the requested file
	filePath := filepath.Join("uploads", fileName)
	err := os.Remove(filePath)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Update the list of uploaded files
	filesLock.Lock()
	defer filesLock.Unlock()

	for i, file := range files {
		if file.Name == fileName {
			files = append(files[:i], files[i+1:]...)
			break
		}
	}

	// Redirect back to the index page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func styleHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/css")
	w.Write([]byte(cssStyle))
}

func waitForShutdown() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	fmt.Println("\nShutting down...")
	// Perform cleanup and shutdown tasks here

	// Exit the application
	os.Exit(0)
}
