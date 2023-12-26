package main

import (
	"cloudnaitivego/internal"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {

	if err := initializeTransactionLog(); err != nil {
		panic(err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/v1/{key}", src.KeyValuePutHandler).Methods("PUT")
	r.HandleFunc("/v1/{key}", src.KeyValueGetHandler).Methods("GET")
	log.Fatal(http.ListenAndServe(":8082", r))
}

func initializeTransactionLog() error {
	var err error
	transact, err := src.NewFileTransactionLogger("transaction.log")
	if err != nil {
		return fmt.Errorf("failed to create transaction logger: %w", err)
	}
	events, errors := transact.ReadEvents()
	e, ok := src.Event{}, true

	for ok && err == nil {
		select {
		case err, ok = <-errors:
		case e, ok = <-events:
			switch e.EventType {
			case src.EventDelete:
				err = src.Delete(e.Key)
			case src.EventPut:
				err = src.Put(e.Key, e.Value)
			}

		}
	}
	transact.Run()
	return err
}
