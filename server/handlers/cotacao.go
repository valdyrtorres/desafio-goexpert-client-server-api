package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"context"
	"log"
	"time"

	"database/sql"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
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

	// ******* Inicio processo de registro no banco
	// Conecta ao banco de dados SQLite
	db, err := sql.Open("sqlite3", "./cotacoes.db?_timeout=10&_journal_mode=WAL")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Cria a tabela se não existir
	createTableSQL := `CREATE TABLE IF NOT EXISTS cotacoes (
		"id" INTEGER PRIMARY KEY AUTOINCREMENT,
		"code" TEXT,
		"codein" TEXT,
		"name" TEXT,
		"high" TEXT,
		"low" TEXT,
		"varBid" TEXT,
		"pctChange" TEXT,
		"bid" TEXT,
		"ask" TEXT,
		"timestamp" TEXT,
		"create_date" TEXT
	);`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}

	// Inicia uma transação
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	// Insere os dados no banco de dados
	for _, resultado := range *cotacao {
		insertSQL := `INSERT INTO cotacoes (code, codein, name, high, low, varBid, pctChange, bid, ask, timestamp, create_date) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
		_, err = tx.Exec(insertSQL, resultado.Code, resultado.Codein, resultado.Nome, resultado.High, resultado.Low, resultado.VarBid, resultado.PctChange, resultado.Bid, resultado.Ask, resultado.Timestamp, resultado.CreateDate)
		if err != nil {
			tx.Rollback() // Desfaz a transação em caso de erro
			log.Fatal(err)
		}
	}

	// Confirma a transação
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Cotação registrada com sucesso!")
	// ****** Fim processo de registro no banco

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
