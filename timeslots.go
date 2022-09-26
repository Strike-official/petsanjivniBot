package main

var fullAvailableTimeslots []string = []string{"10:00 AM", "10:15 AM", "10:30 AM", "10:45 AM", "11:00 AM", "11:15 AM", "11:30 AM", "11:45 AM", "12:00 PM", "12:15 PM", "12:30 PM", "12:45 PM",
	"01:00 PM", "01:15 PM", "01:30 PM", "01:45 PM", "02:00 PM", "02:15 PM", "02:30 PM", "02:45 PM", "03:00 PM", "03:15 PM", "03:30 PM", "03:45 PM",
	"04:00 PM", "04:15 PM", "04:30 PM", "04:45 PM", "05:00 PM", "05:15 PM", "05:30 PM", "05:45 PM", "06:00 PM", "06:15 PM", "06:30 PM", "06:45 PM",
	"07:00 PM", "07:15 PM", "07:30 PM", "07:45 PM", "08:00 PM", "08:15 PM", "08:30 PM", "08:45 PM", "09:00 PM", "09:15 PM", "09:30 PM", "09:45 PM"}

var weekdaysAvailableTimeslots []string = []string{"10:00 AM", "10:15 AM", "10:30 AM", "10:45 AM", "11:00 AM", "11:15 AM", "11:30 AM", "11:45 AM", "12:00 PM", "12:15 PM", "12:30 PM", "12:45 PM",
	"01:00 PM", "01:15 PM", "01:30 PM", "01:45 PM", "02:00 PM", "02:15 PM", "02:30 PM", "02:45 PM", "03:00 PM", "03:15 PM", "03:30 PM", "03:45 PM",
	"04:00 PM", "04:15 PM", "04:30 PM", "04:45 PM", "05:00 PM", "05:15 PM", "05:30 PM", "05:45 PM", "06:00 PM", "06:15 PM", "06:30 PM", "06:45 PM",
	"07:00 PM", "07:15 PM", "07:30 PM", "07:45 PM", "08:00 PM", "08:15 PM", "08:30 PM", "08:45 PM", "09:00 PM", "09:15 PM", "09:30 PM", "09:45 PM"}

var weekendsAvailableTimeslots []string = []string{"10:00 AM", "10:15 AM", "10:30 AM", "10:45 AM", "11:00 AM", "11:15 AM", "11:30 AM", "11:45 AM", "12:00 PM", "12:15 PM", "12:30 PM", "12:45 PM",
	"01:00 PM", "01:15 PM", "01:30 PM", "01:45 PM", "02:00 PM", "02:15 PM", "02:30 PM", "02:45 PM", "03:00 PM", "03:15 PM", "03:30 PM", "03:45 PM",
	"04:00 PM", "04:15 PM", "04:30 PM", "04:45 PM", "05:00 PM", "05:15 PM", "05:30 PM", "05:45 PM", "06:00 PM", "06:15 PM", "06:30 PM", "06:45 PM",
	"07:00 PM", "07:15 PM", "07:30 PM", "07:45 PM", "08:00 PM", "08:15 PM", "08:30 PM", "08:45 PM", "09:00 PM", "09:15 PM", "09:30 PM", "09:45 PM"}

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
	case "02:00 PM":
		return 16
	case "02:15 PM":
		return 17
	case "02:30 PM":
		return 18
	case "02:45 PM":
		return 19
	case "03:00 PM":
		return 20
	case "03:15 PM":
		return 21
	case "03:30 PM":
		return 22
	case "03:45 PM":
		return 23
	case "04:00 PM":
		return 24
	case "04:15 PM":
		return 25
	case "04:30 PM":
		return 26
	case "04:45 PM":
		return 27
	case "05:00 PM":
		return 28
	case "05:15 PM":
		return 29
	case "05:30 PM":
		return 30
	case "05:45 PM":
		return 31
	case "06:00 PM":
		return 32
	case "06:15 PM":
		return 33
	case "06:30 PM":
		return 34
	case "06:45 PM":
		return 35
	case "07:00 PM":
		return 36
	case "07:15 PM":
		return 37
	case "07:30 PM":
		return 38
	case "07:45 PM":
		return 39
	case "08:00 PM":
		return 40
	case "08:15 PM":
		return 41
	case "08:30 PM":
		return 42
	case "08:45 PM":
		return 43
	case "09:00 PM":
		return 44
	case "09:15 PM":
		return 45
	case "09:30 PM":
		return 46
	case "09:45 PM":
		return 47
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
	case "02:00 PM":
		return 16
	case "02:15 PM":
		return 17
	case "02:30 PM":
		return 18
	case "02:45 PM":
		return 19
	case "03:00 PM":
		return 20
	case "03:15 PM":
		return 21
	case "03:30 PM":
		return 22
	case "03:45 PM":
		return 23
	case "04:00 PM":
		return 24
	case "04:15 PM":
		return 25
	case "04:30 PM":
		return 26
	case "04:45 PM":
		return 27
	case "05:00 PM":
		return 28
	case "05:15 PM":
		return 29
	case "05:30 PM":
		return 30
	case "05:45 PM":
		return 31
	case "06:00 PM":
		return 32
	case "06:15 PM":
		return 33
	case "06:30 PM":
		return 34
	case "06:45 PM":
		return 35
	case "07:00 PM":
		return 36
	case "07:15 PM":
		return 37
	case "07:30 PM":
		return 38
	case "07:45 PM":
		return 39
	case "08:00 PM":
		return 40
	case "08:15 PM":
		return 41
	case "08:30 PM":
		return 42
	case "08:45 PM":
		return 43
	case "09:00 PM":
		return 44
	case "09:15 PM":
		return 45
	case "09:30 PM":
		return 46
	case "09:45 PM":
		return 47
	}
	return -1
}
