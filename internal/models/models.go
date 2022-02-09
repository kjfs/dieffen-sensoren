package models

import (
	"time"
)

// Users is the Usermodel, which describes our Users table.
type User struct {
	Id         int
	FirstName  string
	LastName   string
	Email      string
	Phone      string
	UserName   string
	Passwd     string
	Access_lvl int
	Processed  int
	CreatedAt  time.Time
	UpdatedAt  time.Time
} // Needs to describe all fields in the table (check schema).

type SensorSearch struct {
	Id        int
	SensorId  int
	StartDate time.Time
	EndDate   time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Sensor struct {
	Id         int
	SensorName string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type MeasurementKPI struct {
	Id        int
	SensorId  int
	KpiName   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type MeasurementValue struct {
	Id        int
	KpiId     int
	KpiValue  int
	CreatedAt time.Time
	UpdatedAt time.Time
}

// MailData holds an Email message
type MailData struct {
	To      string
	From    string
	Subject string
	//Content template.HTML
	Content  string
	Template string
}
