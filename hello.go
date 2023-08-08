package main

import (
  "fmt"
  "log"
  "net/http"
  "os"
  "strings"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("On my way to unikernels!!!")
	fmt.Fprintf(w, "Hello--unifan !\n")

	fmt.Println("Received request")
	fmt.Println(prettyPrint(r))      // echo to local output
	fmt.Fprintf(w, prettyPrint(r)) // echo to caller
}

func prettyPrint(req *http.Request) string {
	b := &strings.Builder{}
	fmt.Fprintf(b, "%v %v %v %v\n", req.Method, req.URL, req.Proto, req.Host)
	for k, vv := range req.Header {
		for _, v := range vv {
			fmt.Fprintf(b, "  %v: %v\n", k, v)
		}
	}

	if req.Method == "POST" {
		req.ParseForm()
		fmt.Fprintln(b, "Body:")
		for k, v := range req.Form {
			fmt.Fprintf(b, "  %v: %v\n", k, v)
		}
	}

	return b.String()
}

func main() {
  log.Print("helloworld: starting server...")

  http.HandleFunc("/", handler)

  port := os.Getenv("PORT")
  if port == "" {
    port = "8080"
  }

  log.Printf("helloworld: listening on port %s", port)
  log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
