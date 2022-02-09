package dbrepo

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/kjfs/dieffe_sensor/internal/models"
	"golang.org/x/crypto/bcrypt"
)

/*
This is used to store specific functions / methods for Postgres
*/

func (m *postgresDBRepo) InsertUser(user models.User) (int, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var newID int

	stmt := `
	insert into users 
	(first_name, last_name, email, phone, nick_name, passwd, access_lvl, created_at, updated_at) 
	values 
	($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id
	`

	err := m.Db.QueryRowContext(ctx, stmt,
		user.FirstName,
		user.LastName,
		user.Email,
		user.Phone,
		user.UserName,
		user.Passwd,
		user.Access_lvl,
		time.Now(),
		time.Now(),
	).Scan(&newID)

	if err != nil {
		log.Println("EIN FEHLER: ", err)
		return 0, err
	}

	return newID, nil
}

// GetAllUsers returns a slice of users
func (m *postgresDBRepo) GetAllUsers() ([]models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var users []models.User
	query := `
	select 
		u.id, u.first_name, u.last_name, u.email, u.phone, u.nick_name, u.passwd, u.access_lvl, u.processed, u.created_at, u.updated_at
	from 
		users u
	order by
		u.created_at asc
		`
	rows, err := m.Db.QueryContext(ctx, query)
	if err != nil {
		return users, err
	}
	defer rows.Close()
	for rows.Next() {
		var u models.User
		err := rows.Scan(
			&u.Id,
			&u.FirstName,
			&u.LastName,
			&u.Email,
			&u.Phone,
			&u.UserName,
			&u.Passwd,
			&u.Access_lvl,
			&u.Processed,
			&u.CreatedAt,
			&u.UpdatedAt,
		)
		if err != nil {
			return users, err
		}
		users = append(users, u)
	}
	if err = rows.Err(); err != nil {
		return users, err
	}
	return users, nil
}

// GetNewRegistrations returns a slie of registrations
func (m *postgresDBRepo) GetNewUsers() ([]models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var users []models.User
	query := `
	select 
		u.id, u.first_name, u.last_name, u.email, u.phone, u.nick_name, u.passwd, u.access_lvl, u.processed, u.created_at, u.updated_at
	from 
		users u 
	where 
		processed = 0
	order by
		u.created_at asc
		`
	rows, err := m.Db.QueryContext(ctx, query)
	if err != nil {
		return users, err
	}
	defer rows.Close()
	for rows.Next() {
		var u models.User
		err := rows.Scan(
			&u.Id,
			&u.FirstName,
			&u.LastName,
			&u.Email,
			&u.Phone,
			&u.UserName,
			&u.Passwd,
			&u.Access_lvl,
			&u.Processed,
			&u.CreatedAt,
			&u.UpdatedAt,
		)
		if err != nil {
			return users, err
		}
		users = append(users, u)
	}
	if err = rows.Err(); err != nil {
		return users, err
	}
	return users, nil
}

// GetRegistrationById returns a one struct of registration / reservation
func (m *postgresDBRepo) GetUserByID(userID int) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user models.User

	query := `
	select 
		u.id, u.first_name, u.last_name, u.email, u.phone, u.nick_name, u.passwd, u.access_lvl, u.processed, u.created_at, u.updated_at
	from 
		users u
	where 
		u.id = $1
		`
	row := m.Db.QueryRowContext(ctx, query, userID)
	err := row.Scan(
		&user.Id,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Phone,
		&user.UserName,
		&user.Passwd,
		&user.Access_lvl,
		&user.Processed,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return user, err
	}
	return user, nil
}

// updates user from user table by ID
// password will be handled in a different way
func (m *postgresDBRepo) UpdateUser(userID int, user models.User) (bool, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	stmt :=
		`
	update 
		users
	set
		first_name = $1, last_name = $2, email = $3, access_lvl = $4, updated_at = $5, processed = $6
	where
		id = $7
	`

	_, err := m.Db.ExecContext(ctx, stmt,
		user.FirstName,
		user.LastName,
		user.Email,
		user.Access_lvl,
		time.Now(),
		user.Processed,
		userID,
	)

	if err != nil {
		return false, errors.New("could not update user")
	}

	return true, nil
}

// DeleteRegistration deletes a registration by Id.
func (m *postgresDBRepo) DeleteUser(userID int) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	stmt := `
	delete 
	from 
		users
	where
		id = $1
		`
	_, err := m.Db.ExecContext(ctx, stmt, userID)
	if err != nil {
		log.Println("EIN FEHLER:", err)
		return false, errors.New("could not update user")
	}
	return true, nil
}

// AllSensors ...
func (m *postgresDBRepo) AllSensors() ([]models.Sensor, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	var sensors []models.Sensor
	query := `
		select 
			id, sensor_name, created_at, updated_at
		from 
			sensors
		order by 
			sensor_name
		`
	rows, err := m.Db.QueryContext(ctx, query)
	if err != nil {
		log.Println("EIN FEHLER:", err)
		return sensors, err
	}
	defer rows.Close()
	for rows.Next() {
		var sensor models.Sensor
		err := rows.Scan(
			&sensor.Id,
			&sensor.SensorName,
			&sensor.CreatedAt,
			&sensor.UpdatedAt,
		)
		if err != nil {
			return sensors, err
		}
		sensors = append(sensors, sensor)
	}
	if err = rows.Err(); err != nil {
		return sensors, err
	}
	return sensors, nil
}

// Authentication mechanism for normal user
func (m *postgresDBRepo) Authenticate(email, password string) (int, string, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var id int
	var hashedPassword string

	query := `

	select 
		id, passwd
	from
		users
	where
		email = $1
	
	`
	// lets query the db for a userID
	row := m.Db.QueryRowContext(ctx, query, email)
	err := row.Scan(&id, &hashedPassword)
	if err != nil {
		log.Println(err)
		return id, "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword || err == bcrypt.ErrHashTooShort {
		log.Println(err)
		return 0, "", errors.New("incorrect password")
	} else if err != nil {
		return 0, "", err
	}
	return id, hashedPassword, nil

}

// Authentication mechanism for admin user
func (m *postgresDBRepo) AuthenticateAdmin(email, password string) (int, int, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var id int
	var accessLevel int

	query := `

	select 
		id, access_lvl
	from
		users
	where
		email = $1
	
	`
	// lets query the db for a userID
	row := m.Db.QueryRowContext(ctx, query, email)
	err := row.Scan(&id, &accessLevel)
	if err != nil {
		log.Println(err)
		return id, 0, err
	}
	return id, accessLevel, nil

}

// UpdateProcessedVal updates processed value from users table by Id
func (m *postgresDBRepo) UpdateProcessedVal(userID, processedVal int) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	stmt := `
	update 
		users 
	set
		processed = $1, updated_at = $2
	where
		id = $3
		`
	_, err := m.Db.ExecContext(ctx, stmt,
		processedVal,
		time.Now(),
		userID,
	)
	if err != nil {
		log.Println("error: update processed value:", err)
		return false, errors.New("could not update registration")
	}
	return true, nil
}

// GetUserEmail get email value from user table by Id
func (m *postgresDBRepo) GetUserEmail(userID int) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	var email string

	query := `
	select 
		email
	from
		users
	where
		id = $1
		`

	row := m.Db.QueryRowContext(ctx, query, userID)
	err := row.Scan(&email)
	if err != nil {
		log.Println("error in getUserEmail method: ", err)
		return "", errors.New("could not get user email")
	}
	return email, nil
}
