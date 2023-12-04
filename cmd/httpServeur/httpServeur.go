package main

import (
	"log"
	"net/http"

	handle "archi.org/aeroportGo/internal/httpHandlers"
)

func main() {
	log.Default().Println("Démarrage du serveur")
	http.HandleFunc("/", handle.HelloHandler)
	http.HandleFunc("/api/average", handle.AverageHandler)
	http.HandleFunc("/api/valueInterval", handle.IntervalHandler)
	http.HandleFunc("/api/lastValues", handle.LastValuesHandler)
	err := http.ListenAndServe(":80", nil)
	log.Fatal(err)
}
