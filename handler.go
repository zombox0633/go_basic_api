package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/rs/cors"
	"github.com/zombox0633/api/basicService"
)

func greet(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World! 😸 %s", time.Now())
}

func main() {
	fmt.Println("Starting server on 😸 : 5000...")

	// demojson.DemoJSon()
	basicService.WorkRequest()

	//go get -u github.com/rs/cors
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	handler := c.Handler(http.HandlerFunc(greet))
	http.Handle("/", handler)

	err := http.ListenAndServe(":5000", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}

}
