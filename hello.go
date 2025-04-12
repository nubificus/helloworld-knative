package main

import (
	"fmt"
	"net/http"
	"strings"
	"text/template"
)

type Header struct {
	Key   string
	Value string
}

type PageData struct {
	Headers        []Header
	ImageEventURL  string
	ImageRuntime   string
	ImageNubisURL string
}

func determineImageURL(host string) string {
	// Map host to corresponding image URL
	switch {
	case strings.Contains(host, "hellofc"):
		return "https://s3.nbfc.io/hypervisor-logos/firecracker.png"
	case strings.Contains(host, "helloqemu"):
		return "https://s3.nbfc.io/hypervisor-logos/qemu.png"
	case strings.Contains(host, "helloclh"):
		return "https://s3.nbfc.io/hypervisor-logos/clh.png"
	case strings.Contains(host, "hellors"):
		return "https://s3.nbfc.io/hypervisor-logos/dragonball.png"
	case strings.Contains(host, "hellouruncfc"):
		return "https://s3.nbfc.io/hypervisor-logos/uruncfc.png"
	case strings.Contains(host, "hellouruncqemu"):
		return "https://s3.nbfc.io/hypervisor-logos/uruncqemu.png"
	default:
		return "https://s3.nbfc.io/hypervisor-logos/container.png"
	}
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	// Extract headers from the request
	var headers []Header
	for name, values := range r.Header {
		for _, value := range values {
			headers = append(headers, Header{Key: name, Value: value})
		}
	}

	// Get the host from the request headers to determine the image
	host := r.Header.Get("Host")
	imageURL := determineImageURL(host)

	// Prepare the data for the template
	pageData := PageData{
		Headers:        headers,
		ImageEventURL:  "https://s3.nbfc.io/hypervisor-logos/athk8s.png",
		ImageRuntime:   imageURL,
		ImageNubisURL: "https://s3.nbfc.io/hypervisor-logos/nubis-logo-scaled.png",
	}

	// Define the HTML template with a list format and images
	const replyTemplate = `
<!DOCTYPE html>
<html lang="en">
<head><meta charset="UTF-8"><meta name="viewport" content="width=device-width, initial-scale=1.0"><title>Request Headers</title></head>
<body>
<h1>Hello from Knative!</h1>
<div class="logo-row">
    <img src="{{.ImageEventURL}}" alt="Event Logo" />
    <img src="{{.ImageRuntime}}" alt="RuntimeClass" />
</div>
<h2>Request Headers</h2>
<ul class="header-list">
{{range .Headers}}
	<li><strong>{{.Key}}:</strong> {{.Value}}</li>
{{end}}
</ul>
<h2>Brought to you by</h2>
<img src="{{.ImageNubisURL}}" alt="Nubis Logo" style="max-width: 200px;">
</body>
</html>
`

	// Parse and execute the template
	tmpl, err := template.New("webpage").Parse(replyTemplate)
	if err != nil {
		http.Error(w, "Error generating HTML", http.StatusInternalServerError)
		return
	}

	// Render the HTML page with the headers and image data
	err = tmpl.Execute(w, pageData)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}

func main() {
	// Set up the HTTP server and routes
	http.HandleFunc("/", handleRequest)
	fmt.Println("Listening on port 8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
