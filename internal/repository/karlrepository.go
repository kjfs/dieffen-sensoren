package repository

import "github.com/kjfs/dieffe_sensor/internal/models"


type DatabaseRepo interface {
	InsertUser(user models.User) (int, error)
	GetAllUsers() ([]models.User, error)
	GetNewUsers() ([]models.User, error)
	GetUserByID(userID int) (models.User, error)
	GetUserEmail(userID int) (string, error)
	UpdateUser(userID int, user models.User) (bool, error)
	DeleteUser(userID int) (bool, error)
	AllSensors() ([]models.Sensor, error)
	Authenticate(email, password string) (int, string, error)
	AuthenticateAdmin(email, password string) (int, int, error)
	UpdateProcessedVal(userID, processedVal int) (bool, error)
}
