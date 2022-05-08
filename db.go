package main

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"log"

	"github.com/Strike-official/petsanjivniBot/schema"
	_ "github.com/go-sql-driver/mysql"
)

func ConnectToRDS() *sql.DB {
	db, err := sql.Open("mysql", Conf.RDSCredentials)

	if err != nil {
		log.Println("[petsanjivniBot][ERROR][ConnectToRDS] Error establishing connection to DB: ", err)
	}
	return db
}

func getTimeslotFromDB(db *sql.DB, date string) []string {
	results, err := db.Query("select timeslot from petsanjivni_timeslot_details where date_of_appointment = '" + date + "';")
	if err != nil {
		log.Println("[petsanjivniBot][ERROR][getTimeslotFromDB] Error fetching from DB: ", err)
	}

	type wrapper struct {
		timeslot string `json:"timeslot"`
	}
	var wrapperObj wrapper
	var timeslots []string
	for results.Next() {
		err = results.Scan(&wrapperObj.timeslot)
		if err != nil {
			log.Println("[petsanjivniBot][ERROR][getTimeslotFromDB] Error scaning next row DB: ", err)
		}
		timeslots = append(timeslots, wrapperObj.timeslot)
	}
	return timeslots
}

func writeAppointmentToDb(db *sql.DB, request schema.Strike_Meta_Request_Structure, dateOfAppointment string, id string) {

	timeSlot := request.User_session_variables.TimeSlot[0]
	species := request.User_session_variables.PetSpecies[0]

	b := make([]byte, 6)
	if _, err := rand.Read(b); err != nil {
		log.Println("[petsanjivniBot][ERROR][writeAppointmentToDB] Error generating random Id: ", err)
	}
	booking_id := fmt.Sprintf("%X", b)

	_, err := db.Query("INSERT INTO petsanjivni_timeslot_details (booking_id,date_of_appointment,timeslot,species,id) VALUES ('" + booking_id + "','" + dateOfAppointment + "','" + timeSlot + "','" + species + "','" + id + "')")

	if err != nil {
		log.Println("[petsanjivniBot][ERROR][writeAppointmentToDB]: Error writing to DB", err)
	}
}

func writeUserToDb(db *sql.DB, request schema.Strike_Meta_Request_Structure) string {

	// Check if user id is already present
	saved_id, err := getIdFromDb(db, request)
	if err == nil {
		return saved_id
	}

	user_id := request.Bybrisk_session_variables.UserId
	business_id := request.Bybrisk_session_variables.BusinessId
	latitude := fmt.Sprintf("%v", request.Bybrisk_session_variables.Location.Latitude)
	longitude := fmt.Sprintf("%v", request.Bybrisk_session_variables.Location.Longitude)
	username := request.Bybrisk_session_variables.Username
	address := request.Bybrisk_session_variables.Address
	phone := request.Bybrisk_session_variables.Phone

	b := make([]byte, 6)
	if _, err := rand.Read(b); err != nil {
		log.Println("[petsanjivniBot][ERROR][writeUserToDb] Error generating random Id: ", err)
		return "NA"
	}
	id := fmt.Sprintf("%X", b)
	id = "petsanjivni_" + id

	query := "INSERT INTO bybrisk_user_details (id,user_id,business_id,latitude,longitude,username,address,phone) SELECT * FROM (SELECT '" + id + "','" + user_id + "','" + business_id + "'," + latitude + "," + longitude + ",'" + username + "','" + address + "','" + phone + "') AS tmp WHERE NOT EXISTS (SELECT user_id,business_id FROM bybrisk_user_details WHERE user_id='" + user_id + "' and business_id='" + business_id + "') LIMIT 1;"
	_, err = db.Query(query)

	if err != nil {
		log.Println("[petsanjivniBot][ERROR][writeUserToDb]: Error writing to DB", err)
	}
	return id
}

func getIdFromDb(db *sql.DB, request schema.Strike_Meta_Request_Structure) (string, error) {
	user_id := request.Bybrisk_session_variables.UserId
	business_id := request.Bybrisk_session_variables.BusinessId

	var emptyString string

	results, err := db.Query("select id from bybrisk_user_details where user_id = '" + user_id + "' and business_id='" + business_id + "';")
	if err != nil {
		log.Println("[petsanjivniBot][ERROR][getIdFromDb] Error fetching from DB: ", err)
		return emptyString, err
	}

	type wrapper struct {
		id string `json:"id"`
	}

	var wrapperObj wrapper
	var wrapperObjArr []string

	for results.Next() {
		err = results.Scan(&wrapperObj.id)
		if err != nil {
			log.Println("[petsanjivniBot][ERROR][getIdFromDb] Error scaning next row DB: ", err)
			return emptyString, err
		}
		wrapperObjArr = append(wrapperObjArr, wrapperObj.id)
	}

	if len(wrapperObjArr) == 0 {
		return emptyString, fmt.Errorf("User donot exist for this business!")
	}
	return wrapperObjArr[0], nil
}
