package dbrepo

import (
	"errors"

	"github.com/kjfs/dieffe_sensor/internal/models"
)

/*
This is used to store specific functions / methods for testDB
*/

func (m *testDBRepo) InsertUser(user models.User) (int, error) {
	// if the nick_name is kjfs then fail. Otherwise pass.
	if user.UserName == "lokko" {
		return 0, errors.New("wrong nick name")
	}
	return 1, nil
}

// GetSensorByID returns the sensor struct specific to a sensor id
func (m *testDBRepo) GetSensorByID(sensorID int) (models.Sensor, error) {
	var sensor models.Sensor
	if sensorID > 2 {
		return sensor, errors.New("just an error")
	}
	return sensor, nil
}

// GetUserByID return a user from user table by ID
func (m *testDBRepo) GetUserByID(userID int) (models.User, error) {
	var user models.User
	if userID > 1000 {
		return user, errors.New("user id is above 1000. this is not possible. error")
	}
	return user, nil
}

// Authenticate ...
func (m *testDBRepo) Authenticate(email, testPassword string) (int, string, error) {
	if email == "k@invalid.de" {
		return 0, "", errors.New("invalid email address")
	} else {
		return 1, "", nil
	}
}

// Muss angepaßt werden
func (m *testDBRepo) AuthenticateAdmin(email, testPassword string) (int, int, error) {
	if email == "k@invalid.de" {
		return 0, 0, errors.New("invalid email address")
	} else {
		return 1, 0, nil
	}
}

// GetAllUsers returns a slice of users
func (m *testDBRepo) GetAllUsers() ([]models.User, error) {
	var users []models.User
	return users, nil
}

// GetNewUsers returns a slice of users
func (m *testDBRepo) GetNewUsers() ([]models.User, error) {
	var users []models.User
	return users, nil
}

// UpdateUser updates registrations table by Id
func (m *testDBRepo) UpdateUser(userID int, user models.User) (bool, error) {
	return true, nil
}

// UpdateProcessedVal updates processed id from registration table by Id
func (m *testDBRepo) UpdateProcessedVal(registrationId, processedVal int) (bool, error) {
	return true, nil
}

// DeleteUser deletes an user by Id.
func (m *testDBRepo) DeleteUser(userID int) (bool, error) {
	return true, nil
}

// AllSensors ...
func (m *testDBRepo) AllSensors() ([]models.Sensor, error) {
	var sensors []models.Sensor
	return sensors, nil
}

func (m *testDBRepo) GetUserEmail(userID int) (string, error) {
	return "", nil
}
