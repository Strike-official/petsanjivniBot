package main

import (
	"github.com/Strike-official/petsanjivniBot/schema"
	"github.com/strike-official/go-sdk/strike"
)

func CancelAppointment(request schema.Strike_Meta_Request_Structure) *strike.Response_structure {

	// Get all appointment from RDS

	strike_object := strike.Create("cancel_appointment", Conf.BaseURL+"/cancel")
	question_object := strike_object.Question("cancel_appointment_string").QuestionText().SetTextToQuestion("Which appointment to cancel?", "")

	answer_object := question_object.Answer(false).AnswerCardArray(strike.VERTICAL_ORIENTATION)

	answer_object.AnswerCard().SetHeaderToAnswer(2, strike.HALF_WIDTH).AddTextRowToAnswer(strike.H4, "Appointment for Dog", "#8a4c9c", true).
		AddTextRowToAnswer(strike.H5, "15 Mar 2022 @ 04:00 PM", "#c354e3", false).
		AnswerCard().SetHeaderToAnswer(2, strike.HALF_WIDTH).AddTextRowToAnswer(strike.H4, "Appointment for Cat", "#8a4c9c", true).
		AddTextRowToAnswer(strike.H5, "21 Mar 2022 @ 12:15 PM", "#c354e3", false).
		AnswerCard().SetHeaderToAnswer(2, strike.HALF_WIDTH).AddTextRowToAnswer(strike.H4, "Appointment for Cow", "#8a4c9c", true).
		AddTextRowToAnswer(strike.H5, "05 May 2022 @ 08:45 PM", "#c354e3", false)

	return strike_object
}

func Cancel(request schema.Strike_Meta_Request_Structure) *strike.Response_structure {
	// Remove the particular timeslot from RDS

	// Clear the calendar event using ID

	strike_object := strike.Create("cancel_appointment", "")
	strike_object.Question("").QuestionCard().SetHeaderToQuestion(10, strike.HALF_WIDTH).AddTextRowToQuestion(strike.H4, "Appointment cancelled!", "#8a4c9c", true)

	return strike_object

}
