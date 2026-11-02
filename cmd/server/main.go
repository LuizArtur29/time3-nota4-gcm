package main

import (
	"log"
	"net/http"

	"github.com/LuizArtur29/time3-nota4-gcm/internal/calculator"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/multiplicar", calculator.MultiplyHandler)

	log.Println("Microsserviço de multiplicação iniciado na porta 8003")

	if err := http.ListenAndServe(":8003", mux); err != nil {
		log.Fatal(err)
	}
}
