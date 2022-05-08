package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Strike-official/petsanjivniBot/schema"
	"github.com/strike-official/go-sdk/strike"
)

var fullAvailableTimeslots []string = []string{"10:00 AM", "10:15 AM", "10:30 AM", "10:45 AM", "11:00 AM", "11:15 AM", "11:30 AM", "11:45 AM", "12:00 PM", "12:15 PM", "12:30 PM", "12:45 PM",
	"01:00 PM", "01:15 PM", "01:30 PM", "01:45 PM", "02:00 PM", "02:15 PM", "02:30 PM", "02:45 PM", "03:00 PM", "03:15 PM", "03:30 PM", "03:45 PM",
	"04:00 PM", "04:15 PM", "04:30 PM", "04:45 PM", "05:00 PM", "05:15 PM", "05:30 PM", "05:45 PM", "06:00 PM", "06:15 PM", "06:30 PM", "06:45 PM",
	"07:00 PM", "07:15 PM", "07:30 PM", "07:45 PM", "08:00 PM", "08:15 PM", "08:30 PM", "08:45 PM"}

var weekdaysAvailableTimeslots []string = []string{"10:00 AM", "10:15 AM", "10:30 AM", "10:45 AM", "11:00 AM", "11:15 AM", "11:30 AM", "11:45 AM", "12:00 PM", "12:15 PM", "12:30 PM", "12:45 PM",
	"01:00 PM", "01:15 PM", "01:30 PM", "01:45 PM", "05:00 PM", "05:15 PM", "05:30 PM", "05:45 PM", "06:00 PM", "06:15 PM", "06:30 PM", "06:45 PM",
	"07:00 PM", "07:15 PM", "07:30 PM", "07:45 PM", "08:00 PM", "08:15 PM", "08:30 PM", "08:45 PM"}

var weekendsAvailableTimeslots []string = []string{"10:00 AM", "10:15 AM", "10:30 AM", "10:45 AM", "11:00 AM", "11:15 AM", "11:30 AM", "11:45 AM", "12:00 PM", "12:15 PM", "12:30 PM", "12:45 PM",
	"01:00 PM", "01:15 PM", "01:30 PM", "01:45 PM"}

func SelectDate(request schema.Strike_Meta_Request_Structure) *strike.Response_structure {
	name := request.Bybrisk_session_variables.Username
	strike_object := strike.Create("select_date", Conf.BaseURL+"/select_timeslot_species")

	question_object1 := strike_object.Question("date_of_appointment").
		QuestionText().SetTextToQuestion("Hi "+name+"! let's select a date for your appointment", "Text Description, getting used for testing purpose.")

	question_object1.Answer(false).DateInput("Get delivery date")

	return strike_object
}

func SelectTimeSlotAndSpecies(request schema.Strike_Meta_Request_Structure) *strike.Response_structure {

	strike_object := strike.Create("select_timeslot", Conf.BaseURL+"/save_data?date="+request.User_session_variables.DateOfAppointment[0])
	question_object1 := strike_object.Question("time_slot").
		QuestionText().SetTextToQuestion("Select timeslot", "Text Description, getting used for testing purpose.")

	answer_object1 := question_object1.Answer(false).AnswerCardArray(strike.VERTICAL_ORIENTATION)

	availableTimeslots := getTimeslotsForDate(request.User_session_variables.DateOfAppointment[0])

	weekday := getWeekdayFromDate(request.User_session_variables.DateOfAppointment[0])
	for i, v := range availableTimeslots {
		if weekday == "Sunday" {
			if v == "booked" {
				answer_object1 = answer_object1.AnswerCard().SetHeaderToAnswer(10, strike.HALF_WIDTH).AddTextRowToAnswer(strike.H6, weekendsAvailableTimeslots[i], "#acadad", true)
			} else {
				answer_object1 = answer_object1.AnswerCard().SetHeaderToAnswer(1, strike.HALF_WIDTH).AddTextRowToAnswer(strike.H5, v, "#009646", true)
			}
		} else {
			if v == "booked" {
				answer_object1 = answer_object1.AnswerCard().SetHeaderToAnswer(10, strike.HALF_WIDTH).AddTextRowToAnswer(strike.H6, weekdaysAvailableTimeslots[i], "#acadad", true)
			} else {
				answer_object1 = answer_object1.AnswerCard().SetHeaderToAnswer(1, strike.HALF_WIDTH).AddTextRowToAnswer(strike.H5, v, "#009646", true)
			}
		}

	}

	question_object2 := strike_object.Question("pet_species").
		QuestionText().SetTextToQuestion("Which pet to book for?", "Text Description, getting used for testing purpose.")

	answer_object2 := question_object2.Answer(false).AnswerCardArray(strike.VERTICAL_ORIENTATION)

	answer_object2.AnswerCard().SetHeaderToAnswer(1, strike.HALF_WIDTH).AddTextRowToAnswer(strike.H5, "Dog", "#009646", true).
		AnswerCard().SetHeaderToAnswer(1, strike.HALF_WIDTH).AddTextRowToAnswer(strike.H5, "Cat", "#009646", true).
		AnswerCard().SetHeaderToAnswer(1, strike.HALF_WIDTH).AddTextRowToAnswer(strike.H5, "Cow", "#009646", true)

	return strike_object
}

func SaveData(request schema.Strike_Meta_Request_Structure, dateOfAppointment string) *strike.Response_structure {

	// save the data to RDS
	id := writeUserToDb(RDS_db_connection, request)
	writeAppointmentToDb(RDS_db_connection, request, dateOfAppointment, id)
	// book on google calender
	bookOnCalender(request, dateOfAppointment)

	strike_object := strike.Create("select_timeslot", "")
	question_object := strike_object.Question("").
		QuestionCard().SetHeaderToQuestion(10, strike.HALF_WIDTH).AddTextRowToQuestion(strike.H4, "Appointment confirmed!", "#009646", true)

	answer_object := question_object.Answer(false).AnswerCardArray(strike.VERTICAL_ORIENTATION)
	answer_object.AnswerCard().SetHeaderToAnswer(10, strike.HALF_WIDTH).AddTextRowToAnswer(strike.H4, "Bring your "+request.User_session_variables.PetSpecies[0]+" on "+dateOfAppointment+" \nTimeslot - "+request.User_session_variables.TimeSlot[0], "#73767a", true).
		AnswerCard().SetHeaderToAnswer(10, strike.HALF_WIDTH).AddTextRowToAnswer(strike.H4, "Get in touch - \n+91-9617879980\n+91-7582974462", "black", false)

	return strike_object
}

func getTimeslotsForDate(date string) []string {

	var tempTimeslot []string

	// Get the booked slots by date from DB
	bookedTimeslots := getTimeslotFromDB(RDS_db_connection, date)

	weekday := getWeekdayFromDate(date)
	if weekday == "Sunday" {

		tempTimeslot = append(tempTimeslot, weekendsAvailableTimeslots...)
		for _, v := range bookedTimeslots {
			index := getIndexForWeekendTimeslot(v)
			tempTimeslot = booked(tempTimeslot, index)
		}
	} else {
		tempTimeslot = append(tempTimeslot, weekdaysAvailableTimeslots...)
		for _, v := range bookedTimeslots {
			index := getIndexForWeekdayTimeslot(v)
			tempTimeslot = booked(tempTimeslot, index)
		}
	}

	return tempTimeslot
}

func booked(s []string, index int) []string {
	if index == -1 {
		return s
	}
	s[index] = "booked"
	return s
}

func getWeekdayFromDate(d string) string {
	dArr := strings.Split(d, " ")
	newdate := dArr[2] + "-" + dArr[1] + "-" + dArr[0]
	t, err := time.Parse("2006-Jan-02", newdate)
	if err != nil {
		log.Println("[petsanjivniBot][ERROR][getWeekdayFromDate] Error parsing time: ", err)
	}
	return t.Weekday().String()
}

func getIndexForWeekdayTimeslot(ts string) int {
	switch ts {
	case "10:00 AM":
		return 0
	case "10:15 AM":
		return 1
	case "10:30 AM":
		return 2
	case "10:45 AM":
		return 3
	case "11:00 AM":
		return 4
	case "11:15 AM":
		return 5
	case "11:30 AM":
		return 6
	case "11:45 AM":
		return 7
	case "12:00 PM":
		return 8
	case "12:15 PM":
		return 9
	case "12:30 PM":
		return 10
	case "12:45 PM":
		return 11
	case "01:00 PM":
		return 12
	case "01:15 PM":
		return 13
	case "01:30 PM":
		return 14
	case "01:45 PM":
		return 15
	case "05:00 PM":
		return 16
	case "05:15 PM":
		return 17
	case "05:30 PM":
		return 18
	case "05:45 PM":
		return 19
	case "06:00 PM":
		return 20
	case "06:15 PM":
		return 21
	case "06:30 PM":
		return 22
	case "06:45 PM":
		return 23
	case "07:00 PM":
		return 24
	case "07:15 PM":
		return 25
	case "07:30 PM":
		return 26
	case "07:45 PM":
		return 27
	case "08:00 PM":
		return 28
	case "08:15 PM":
		return 29
	case "08:30 PM":
		return 30
	case "08:45 PM":
		return 31
	}
	return -1
}

func getIndexForWeekendTimeslot(ts string) int {
	switch ts {
	case "10:00 AM":
		return 0
	case "10:15 AM":
		return 1
	case "10:30 AM":
		return 2
	case "10:45 AM":
		return 3
	case "11:00 AM":
		return 4
	case "11:15 AM":
		return 5
	case "11:30 AM":
		return 6
	case "11:45 AM":
		return 7
	case "12:00 PM":
		return 8
	case "12:15 PM":
		return 9
	case "12:30 PM":
		return 10
	case "12:45 PM":
		return 11
	case "01:00 PM":
		return 12
	case "01:15 PM":
		return 13
	case "01:30 PM":
		return 14
	case "01:45 PM":
		return 15
	}
	return -1
}

func bookOnCalender(request schema.Strike_Meta_Request_Structure, dateOfAppointment string) {

	date := dateOfAppointment
	timeSlot := request.User_session_variables.TimeSlot[0]
	name := request.Bybrisk_session_variables.Username
	phone := request.Bybrisk_session_variables.Phone
	species := request.User_session_variables.PetSpecies[0]
	location := fmt.Sprintf("%f", request.Bybrisk_session_variables.Location.Latitude) + "," + fmt.Sprintf("%f", request.Bybrisk_session_variables.Location.Longitude)

	data := date + ";" + timeSlot + ";" + name + ";" + phone + ";" + species + ";" + location
	data = strings.Replace(data, " ", "%20", -1)

	response, err := http.Get(Conf.GCApi + "?data=" + data)
	if err != nil {
		log.Println("[petsanjivniBot][ERROR][bookOnCalender] Error Calender GET API: ", err, " data: ", data)
	}

	type googleScriptResponse struct {
		Status string `json:"status"`
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println("[petsanjivniBot][ERROR][bookOnCalender] Error ioutil.ReadAll: ", err)
	}

	var resp googleScriptResponse

	if err := json.Unmarshal(responseData, &resp); err != nil {
		log.Println("[petsanjivniBot][ERROR][bookOnCalender] Error Unmarshalling response data: ", err)
	}

	if resp.Status != "OK" {
		log.Println("[petsanjivniBot][ERROR][bookOnCalender] Error response Status not OK")
	}
}
