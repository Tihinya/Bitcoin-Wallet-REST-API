package main

import (
	"bitcoin-wallet/config"
	"bitcoin-wallet/controllers"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	config.ParseConfig("./dev_config.json")
	r := mux.NewRouter()
	r.HandleFunc("/transactions/get", controllers.GetListTransactions).Methods("GET")
	r.HandleFunc("/balance/get", controllers.GetCurrentBalance).Methods("GET")
	r.HandleFunc("/transfer", controllers.Transfer).Methods("POST")
	r.HandleFunc("/top-up", controllers.TopUp).Methods("POST")

	http.Handle("/", r)

	log.Println("Ctrl + Click on the link: http://localhost:" + config.ConfigFile.Port)
	log.Println("To stop the server press `Ctrl + C`")

	log.Fatalln(http.ListenAndServe(":"+config.ConfigFile.Port, nil))
}
