package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"context"
	"log"
	"time"

	"github.com/gorilla/mux"
)

type AweSomeApi struct {
	Code       string `json:"code"`
	Codein     string `json:"codein"`
	Nome       string `json:"name"`
	High       string `json:"high"`
	Low        string `json:"low"`
	VarBid     string `json:"varBid"`
	PctChange  string `json:"pctChange"`
	Bid        string `json:"bid"`
	Ask        string `json:"ask"`
	Timestamp  string `json:"timestamp"`
	CreateDate string `json:"create_date"`
}

func CotacaoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	modoCambioParam, ok := vars["cambio"]
	if !ok || modoCambioParam == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	cotacao, err := PegaCotacao(modoCambioParam)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// eu quero somente o bid para fornecer ao client.go
	for _, valor := range *cotacao {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"bid": valor.Bid})
		return
	}

}

func CotacaoFullHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	modoCambioParam, ok := vars["cambio"]
	if !ok || modoCambioParam == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	cotacao, err := PegaCotacao(modoCambioParam)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(cotacao)

}

func PegaCotacao(modoCambio string) (*map[string]AweSomeApi, error) {

	req, err := http.NewRequest(http.MethodGet, "https://economia.awesomeapi.com.br/json/last/"+modoCambio, nil)
	if err != nil {
		log.Fatal(err)
	}

	// timeout máximo para chamar a API de cotação do dólar deverá ser de 200ms .
	// Utilizando o package "context"
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Millisecond*200))
	defer cancel()
	req = req.WithContext(ctx)

	c := &http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Usando mapa, pois o cambio pode mudar
	var resultado map[string]AweSomeApi

	err = json.Unmarshal(body, &resultado)
	if err != nil {
		return nil, err
	}
	return &resultado, nil
}
