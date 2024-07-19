package main

import (
	"fmt"
	"net/http"
	"text/template"

	"github.com/billzayy/url-shorten/pkg"
	"github.com/rs/cors"
)

func main() {

	mux := http.NewServeMux()

	options := cors.Options{
		AllowedOrigins:   []string{"*"}, // Adjust according to your needs
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		AllowCredentials: true,
	}

	c := cors.New(options)
	corsHandler := c.Handler(mux)

	mux.HandleFunc("POST /api/short-url", pkg.GenerateURL)
	mux.HandleFunc("GET /redirect/", pkg.RedirectURL)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("./templates/index.html"))
		tmpl.Execute(w, nil)
	})

	fmt.Println("Listening on port: 8080")
	err := http.ListenAndServe(":8080", corsHandler)

	if err != nil {
		panic(err)
	}
}
