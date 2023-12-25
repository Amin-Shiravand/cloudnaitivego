package main

import (
	"cloudnaitivego/cloud_Native"
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
	r.HandleFunc("/v1/{key}", cloud_Native.KeyValuePutHandler).Methods("PUT")
	r.HandleFunc("/v1/{key}", cloud_Native.KeyValueGetHandler).Methods("GET")
	log.Fatal(http.ListenAndServe(":8082", r))
}

func initializeTransactionLog() error {
	var err error
	transact, err := cloud_Native.NewFileTransactionLogger("transaction.log")
	if err != nil {
		return fmt.Errorf("failed to create transaction logger: %w", err)
	}
	events, errors := transact.ReadEvents()
	e, ok := cloud_Native.Event{}, true

	for ok && err == nil {
		select {
		case err, ok = <-errors:
		case e, ok = <-events:
			switch e.EventType {
			case cloud_Native.EventDelete:
				err = cloud_Native.Delete(e.Key)
			case cloud_Native.EventPut:
				err = cloud_Native.Put(e.Key, e.Value)
			}

		}
	}
	transact.Run()
	return err
}
