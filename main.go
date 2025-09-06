package main

import (
    "fmt"
    "net/http"
		"davidhampgonsalves/lifedashboard/pkg"
		"io"
		"bytes"
)

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/" {
				http.NotFound(w, r)
				return
			}

			png := pkg.Generate();
			w.Header().Set("Content-Type", "image/png") // <-- set the content-type header
			io.Copy(w, bytes.NewBuffer(png))
    })

    fmt.Println("Server listening on http://localhost:8080")
    http.ListenAndServe(":8080", nil)
}