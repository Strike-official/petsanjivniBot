package main

import (
	"testing"

	"github.com/Strike-official/petsanjivniBot/schema"
)

func TestPushNotification(t *testing.T) {
	request := schema.Strike_Meta_Request_Structure{
		Bybrisk_session_variables: schema.Bybrisk_session_variables_struct{
			UserId:     "623efb9195ba637fe92fb07b",
			BusinessId: "624acdf02c5817b22d7b303f",
		},
		User_session_variables: schema.User_session_variables_struct{
			TimeSlot:   []string{"09:10 AM"},
			PetSpecies: []string{"Cat"},
		},
	}
	date := "2022-Sep-26"
	pushNotification(request, date)
}
