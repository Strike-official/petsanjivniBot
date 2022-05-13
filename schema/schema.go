package schema

import "time"

type Strike_Meta_Request_Structure struct {

	// Bybrisk variable from strike bot
	//
	Bybrisk_session_variables Bybrisk_session_variables_struct `json: "bybrisk_session_variables"`

	// Our own variable from previous API
	//
	User_session_variables User_session_variables_struct `json: "user_session_variables"`
}

type Bybrisk_session_variables_struct struct {

	// User ID on Bybrisk
	//
	UserId string `json:"userId"`

	// Our own business Id in Bybrisk
	//
	BusinessId string `json:"businessId"`

	// Handler Name for the API chain
	//
	Handler string `json:"handler"`

	// Current location of the user
	//
	Location GeoLocation_struct `json:"location"`

	// Username of the user
	//
	Username string `json:"username"`

	// Address of the user
	//
	Address string `json:"address"`

	// Phone number of the user
	//
	Phone string `json:"phone"`
}

type GeoLocation_struct struct {
	// Latitude
	//
	Latitude float64 `json:"latitude"`

	// Longitude
	//
	Longitude float64 `json:"longitude"`
}

type User_session_variables_struct struct {
	TimeSlot                []string `json:"time_slot,omitempty"`
	DateOfAppointment       []string `json:"date_of_appointment,omitempty"`
	PetSpecies              []string `json:"pet_species,omitempty"`
	CancelAppointmentString []string `json:"cancel_appointment_string,omitempty"`
}

type Wrapper struct {
	Item_description string    `json:"item_description"`
	Item_total       float64   `json:"item_total"`
	Quantity         string    `json:"quantity"`
	Order_time       time.Time `json:"order_time"`
	Delivery_date    string    `json:"delivery_date"`
}

type Booking_Details_Db_Wrapper struct {
	BookingID         string    `json:"booking_id"`
	DateOfAppointment string    `json:"date_of_appointment"`
	Timeslot          string    `json:"timeslot"`
	Species           string    `json:"species"`
	ID                string    `json:"id"`
	DateCreated       time.Time `json:"date_created"`
}
