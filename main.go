package main

import (
	"fmt"
	"net/http"

	"github.com/usuario/desafio-go/handlers"
)

func main() {
	// Criar o ServeMux (multiplexador) para roteamento de URLs
	muxDesafio := http.NewServeMux()

	// Registrar a rota "/" para o handler HomeHandler
	muxDesafio.HandleFunc("/", handlers.HomeHandler)

	muxDesafio.HandleFunc("/sobre", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Este é um programa de exemplo em Go.")
	})

	muxDesafio.HandleFunc("/contato", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Entre em contato: contato@example.com")
	})

	// Iniciar o servidor na porta 8080
	fmt.Println("Servidor rodando na porta 8080...")
	http.ListenAndServe(":8080", muxDesafio)
}
