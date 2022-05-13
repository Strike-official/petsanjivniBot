package main

import (
	"log"
	"time"

	"github.com/Strike-official/petsanjivniBot/schema"
	"github.com/strike-official/go-sdk/strike"
)

func History(request schema.Strike_Meta_Request_Structure) *strike.Response_structure {

	name := request.Bybrisk_session_variables.Username

	strike_object := strike.Create("history", "")
	question_object := strike_object.Question("").QuestionText().SetTextToQuestion("Hi "+name+", Here is your appointment history", "desc")

	answer_object := question_object.Answer(false).AnswerCardArray(strike.VERTICAL_ORIENTATION)

	// Get booking data from db
	data, err := getBookingsFromDb(RDS_db_connection, request)
	if err != nil || len(data) == 0 {
		log.Println("[petsanjivniBot][ERROR][History] Error getting booking data from DB: ", err)

		answer_object.AnswerCard().SetHeaderToAnswer(10, strike.HALF_WIDTH).AddTextRowToAnswer(strike.H2, "No Appointments yet!", "black", true)
		return strike_object
	}

	// Show booking data

	for _, v := range data {
		answer_object = answer_object.AnswerCard().SetHeaderToAnswer(10, strike.HALF_WIDTH).AddTextRowToAnswer(strike.H4, "ID - "+v.BookingID, "#474747", true)

		now := time.Now().Format("02 Jan 2006")
		if v.DateOfAppointment == now {
			answer_object = answer_object.AddTextRowToAnswer(strike.H5, "Appointment for today @ "+v.Timeslot, "#de680d", true)
		} else {
			timeStampString := v.DateOfAppointment
			layOut := "02 Jan 2006"
			timeStamp, err := time.Parse(layOut, timeStampString)
			if err != nil {
				log.Println("Error parsing timestamp: ", err)
			}

			if timeStamp.Unix() > time.Now().Unix() {
				answer_object = answer_object.AddTextRowToAnswer(strike.H5, "Appointment on "+v.DateOfAppointment+" @ "+v.Timeslot, "#009646", true)
			} else {
				answer_object = answer_object.AddTextRowToAnswer(strike.H5, "Appointment was on "+v.DateOfAppointment+" @ "+v.Timeslot, "#ff002b", true)
			}

		}

		answer_object.AddTextRowToAnswer(strike.H5, "Booked on "+v.DateCreated.Format("Mon, 02 Jan 2006"), "black", true).
			AddTextRowToAnswer(strike.H5, "for "+v.Species, "black", false)
	}

	return strike_object

}
