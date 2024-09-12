package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type Dolar struct {
	Bid string `json:"bid"`
}

func main() {
	req, err := http.NewRequest(http.MethodGet, "http://localhost:8080/cotacao/USD-BRL", nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao criar requisição: %v\n", err)
		log.Fatal(err)
	}

	// O client.go precisará receber do server.go apenas o valor atual do câmbio (campo "bid" do JSON).
	// Utilizando o package "context", o client.go terá um timeout máximo de 300ms
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Millisecond*300))
	defer cancel()
	req = req.WithContext(ctx)

	c := &http.Client{}
	res, err := c.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao fazer requisição: %v\n", err)
		log.Fatal(err)
	}
	defer res.Body.Close()

	out, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao ler resposta: %v\n", err)
		log.Fatal(err)
	}
	log.Println(string(out))

	var data Dolar
	err = json.Unmarshal(out, &data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao fazer parse da resposta: %v\n", err)
	}

	file, err := os.Create("cotacao.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao criar arquivo: %v\n", err)
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("Dólar: %s", data.Bid))
	fmt.Println("Arquivo criado com sucesso!")
	fmt.Println("Dólar: ", data.Bid)
}
