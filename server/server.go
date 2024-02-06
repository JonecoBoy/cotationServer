package main

import (
	"encoding/json"
	"fmt"
	"github.com/JonecoBoy/cotationServer/server/cotation"
	"github.com/JonecoBoy/cotationServer/server/db"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
)

func main() {

	mux := http.NewServeMux()
	db.DatabaseBuilder()
	// podia ter passado anonima
	mux.HandleFunc("/", HomeHandler)

	mux.HandleFunc("/cotacao", getExchange)

	log.Print("Listening...")
	http.ListenAndServe(":8080", mux)

}
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}

func getExchange(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	from := queryParams.Get("from")
	if from == "" {
		from = "usd"
	}
	to := queryParams.Get("to")
	if to == "" {
		to = "brl"
	}
	mode := queryParams.Get("mode")

	c, err := cotation.GetCotation(from, to)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	err = db.InsertCotation(&c)
	if err != nil {
		panic(err)
	}

	if mode != "detailed" {
		w.Write([]byte(c.Bid))
	} else {
		// Marshal the Cotation struct to JSON
		jsonData, err := json.Marshal(c)
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Set the Content-Type header to application/json
		w.Header().Set("Content-Type", "application/json")

		// Write the JSON data to the response
		w.Write(jsonData)
	}

}
