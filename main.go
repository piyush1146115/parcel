package main

import (
	"github.com/gorilla/mux"
	"github.com/piyush1146115/parcel/handler"
	"net/http"
)

func main() {

	sm := mux.NewRouter()

	postR := sm.Methods(http.MethodPost).Subrouter()
	postR.HandleFunc("/api/v1/parcel", handler.HandleNewParcelRequest)

}
