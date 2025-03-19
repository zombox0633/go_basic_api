package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/zombox0633/api/basicService"
)

func greet(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World! 😸 %s", time.Now())
}

func main() {
	// demojson.DemoJSon()
	basicService.WorkRequest()

	http.HandleFunc("/", greet)
	http.ListenAndServe(":5000", nil)
}
