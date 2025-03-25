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
	fmt.Println("Starting server on 😸 : 5000...")
	// demojson.DemoJSon()
	basicService.WorkRequest()

	http.HandleFunc("/", greet)
	err := http.ListenAndServe(":5000", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}

}
