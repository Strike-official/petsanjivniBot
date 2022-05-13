package main

import (
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/Strike-official/petsanjivniBot/schema"
	"github.com/gorilla/mux"
	"github.com/strike-official/go-sdk/strike"
)

var Conf *AppConfig
var RDS_db_connection *sql.DB

func determineListenAddress() (string, error) {
	port := os.Getenv("PORT")
	if port == "" {
		port = Conf.Port
	}
	return ":" + port, nil
}

func main() {
	err := InitAppConfig("config.json")
	if err != nil {
		log.Println("[petsanjivniBot][ERROR] Error Initializing AppConfig: ", err)
	}
	Conf = GetAppConfig()

	addr, err := determineListenAddress()
	if err != nil {
		log.Fatal(err)
	}
	RDS_db_connection = ConnectToRDS()
	r := mux.NewRouter()

	r.HandleFunc("/petsanjivniBot/select_date", select_date).Methods("POST")
	r.HandleFunc("/petsanjivniBot/select_timeslot_species", select_timeslot_species).Methods("POST")
	r.HandleFunc("/petsanjivniBot/save_data", save_data).Methods("POST")
	r.HandleFunc("/petsanjivniBot/cancel_appointment", cancel_appointment).Methods("POST")
	r.HandleFunc("/petsanjivniBot/cancel", cancel).Methods("POST")
	r.HandleFunc("/petsanjivniBot/history", history).Methods("POST")
	http.Handle("/petsanjivniBot/", r)

	if err := http.ListenAndServe(addr, nil); err != nil {
		panic(err)
	}
}

//#####################################################HANDLERS###########################################################
func select_date(w http.ResponseWriter, r *http.Request) {
	WriteResponse(w, SelectDate(DecodeRequestToStruct(r.Body)))
}

func select_timeslot_species(w http.ResponseWriter, r *http.Request) {
	WriteResponse(w, SelectTimeSlotAndSpecies(DecodeRequestToStruct(r.Body)))
}

func save_data(w http.ResponseWriter, r *http.Request) {
	date := r.URL.Query()["date"][0]
	WriteResponse(w, SaveData(DecodeRequestToStruct(r.Body), date))
}

func cancel_appointment(w http.ResponseWriter, r *http.Request) {
	WriteResponse(w, CancelAppointment(DecodeRequestToStruct(r.Body)))
}

func cancel(w http.ResponseWriter, r *http.Request) {
	WriteResponse(w, Cancel(DecodeRequestToStruct(r.Body)))
}

func history(w http.ResponseWriter, r *http.Request) {
	WriteResponse(w, History(DecodeRequestToStruct(r.Body)))
}

//##############################################HELPER##########################################################

func DecodeRequestToStruct(r io.Reader) schema.Strike_Meta_Request_Structure {
	decoder := json.NewDecoder(r)
	var request schema.Strike_Meta_Request_Structure
	err := decoder.Decode(&request)
	if err != nil {
		log.Println("[petsanjivniBot][ERROR][DecodeRequestToStruct] Error in decoding request: ", err)
	}
	log.Println("[petsanjivniBot][INFO]: ", request)
	return request
}

func WriteResponse(w http.ResponseWriter, s *strike.Response_structure) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(s.ToJson())
}
