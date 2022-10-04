package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Strike-official/petsanjivniBot/schema"
	"github.com/strike-official/go-sdk/strike"
)

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

	answer_object1 := question_object1.Answer(false).AnswerCardArray(strike.HORIZONTAL_ORIENTATION)

	availableTimeslots := getTimeslotsForDate(request.User_session_variables.DateOfAppointment[0])

	weekday := getWeekdayFromDate(request.User_session_variables.DateOfAppointment[0])
	for i, v := range availableTimeslots {
		if weekday == "Sunday" {
			if v == "booked" {
				answer_object1 = answer_object1.AnswerCard().SetHeaderToAnswer(10, "WRAP").AddTextRowToAnswer(strike.H6, weekendsAvailableTimeslots[i], "#acadad", true)
			} else {
				answer_object1 = answer_object1.AnswerCard().SetHeaderToAnswer(1, "WRAP").AddTextRowToAnswer(strike.H5, v, "#009646", true)
			}
		} else {
			if v == "booked" {
				answer_object1 = answer_object1.AnswerCard().SetHeaderToAnswer(10, "WRAP").AddTextRowToAnswer(strike.H6, weekdaysAvailableTimeslots[i], "#acadad", true)
			} else {
				answer_object1 = answer_object1.AnswerCard().SetHeaderToAnswer(1, "WRAP").AddTextRowToAnswer(strike.H5, v, "#009646", true)
			}
		}

	}

	question_object2 := strike_object.Question("pet_species").
		QuestionText().SetTextToQuestion("Which pet to book for?", "Text Description, getting used for testing purpose.")

	answer_object2 := question_object2.Answer(false).AnswerCardArray(strike.HORIZONTAL_ORIENTATION)

	answer_object2.AnswerCard().SetHeaderToAnswer(1, "WRAP").AddTextRowToAnswer(strike.H5, "Dog", "#009646", true).
		AnswerCard().SetHeaderToAnswer(1, "WRAP").AddTextRowToAnswer(strike.H5, "Cat", "#009646", true).
		AnswerCard().SetHeaderToAnswer(1, "WRAP").AddTextRowToAnswer(strike.H5, "Rabbit", "#009646", true).
		AnswerCard().SetHeaderToAnswer(1, "WRAP").AddTextRowToAnswer(strike.H5, "Other", "#009646", true)

	return strike_object
}

func SaveData(request schema.Strike_Meta_Request_Structure, dateOfAppointment string) *strike.Response_structure {

	// save the data to RDS
	id := writeUserToDb(RDS_db_connection, request)
	writeAppointmentToDb(RDS_db_connection, request, dateOfAppointment, id)
	// book on google calender
	bookOnCalender(request, dateOfAppointment)
	// schedule push notification
	pushNotification(request, dateOfAppointment)

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

	return removeElapsedTimeSlots(tempTimeslot, date)
}

func removeElapsedTimeSlots(selectedTimeSlots []string, date string) []string {
	loc, _ := time.LoadLocation("Asia/Kolkata")
	currTimeStamp := time.Now().In(loc)
	currDate := currTimeStamp.Format("02 Jan 2006")

	log.Println("removeElapsedTimeSlots:", selectedTimeSlots, "date:", date, "currDate:", currDate)
	if date == currDate {
		var timeSlotsWithoutElapsedTime []string
		log.Println("It's today!")
		for _, selectedTimeSlot := range selectedTimeSlots {
			if selectedTimeSlot == "booked" {
				continue
			}
			hourIn24HourFormat := getHourFromTimeSlot(selectedTimeSlot)
			hour, minute, _ := currTimeStamp.Clock()
			if hourIn24HourFormat == hour {
				minuteFromTimeSlot := getMinuteFromTimeSlot(selectedTimeSlot)
				if minuteFromTimeSlot > minute {
					timeSlotsWithoutElapsedTime = append(timeSlotsWithoutElapsedTime, selectedTimeSlot)
				}
			}
			if hourIn24HourFormat > hour {
				timeSlotsWithoutElapsedTime = append(timeSlotsWithoutElapsedTime, selectedTimeSlot)
			}

		}
		selectedTimeSlots = timeSlotsWithoutElapsedTime
	}
	return selectedTimeSlots
}

func getHourFromTimeSlot(selectedTimeSlot string) int {
	selectedTimeSlotArr := strings.Split(selectedTimeSlot, ":")
	selectedHour := selectedTimeSlotArr[0]
	selectedHour = trimZero(selectedHour)

	selectedHourInt, err := strconv.Atoi(selectedHour)
	if err != nil {
		log.Println("Error parsing selectedHour string to int:", err)
	}
	return formatHourIn24TimeFormat(selectedTimeSlot, selectedHourInt)
}

func getMinuteFromTimeSlot(selectedTimeSlot string) int {
	selectedTimeSlotArr := strings.Split(selectedTimeSlot, ":")
	remainderTimeslotStringArr := strings.Split(selectedTimeSlotArr[1], " ")
	selectedMinute := remainderTimeslotStringArr[0]
	selectedMinute = trimZero(selectedMinute)
	selectedMinuteInt, err := strconv.Atoi(selectedMinute)
	if err != nil {
		log.Println("Error parsing selectedMinute string to int:", err)
	}
	return selectedMinuteInt
}

func formatHourIn24TimeFormat(selectedTimeSlot string, selectedHourInt int) int {
	hourIn24HourFormat := selectedHourInt
	selectedTimeSlotArr := strings.Split(selectedTimeSlot, ":")
	remainderTimeslotStringArr := strings.Split(selectedTimeSlotArr[1], " ")
	if remainderTimeslotStringArr[1] == "PM" {
		if hourIn24HourFormat == 12 {
			hourIn24HourFormat = 0
		}
		hourIn24HourFormat = hourIn24HourFormat + 12
	}

	if remainderTimeslotStringArr[1] == "AM" && hourIn24HourFormat == 12 {
		hourIn24HourFormat = 0
	}
	return hourIn24HourFormat
}
func trimZero(s string) string {
	iszero := s[0:1]
	if iszero == "0" {
		s = s[1:]
	}
	return s
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

func pushNotification(request schema.Strike_Meta_Request_Structure, dateOfAppointment string) {
	timeString, dateString := getTimeAndDateStringUTC(request, dateOfAppointment)
	species := request.User_session_variables.PetSpecies[0]
	response := strike.Notification(request.Bybrisk_session_variables.UserId, request.Bybrisk_session_variables.BusinessId).SetContent("🔔 You have an appointment for your " + species + " in 10 minutes").SetTargetTimeUTC(timeString).SetTargetDateUTC(dateString).Do()
	fmt.Println("----->response from notification:", *response.NotificationID, *response.Result, *response.Status)
}

func getTimeAndDateStringUTC(request schema.Strike_Meta_Request_Structure, dateOfAppointment string) (string, string) {
	timeSlot := request.User_session_variables.TimeSlot[0]
	timeToParse := dateOfAppointment + " " + timeSlot
	location, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		log.Println("Error loading time location:", err)
	}
	formatedTime, err := time.ParseInLocation("02 Jan 2006 03:04 PM", timeToParse, location)
	if err != nil {
		log.Println("Error parsing time:", err)
	}
	formattedtimeUTC := formatedTime.UTC().Add(-600000000000)
	hour, min, second := formattedtimeUTC.Clock()
	year, month, day := formattedtimeUTC.Date()
	timeString := appendZeroToInt(hour) + ":" + appendZeroToInt(min) + ":" + appendZeroToInt(second)
	dateString := appendZeroToInt(year) + "-" + appendZeroToInt(int(month)) + "-" + appendZeroToInt(day)
	return timeString, dateString
}

func appendZeroToInt(i int) string {
	if i/10 == 0 {
		return "0" + strconv.Itoa(i)
	}
	return strconv.Itoa(i)
}
