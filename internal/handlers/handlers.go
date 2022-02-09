package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
	"github.com/kjfs/dieffe_sensor/helpers"
	"github.com/kjfs/dieffe_sensor/internal/config"
	"github.com/kjfs/dieffe_sensor/internal/driver"
	"github.com/kjfs/dieffe_sensor/internal/forms"
	"github.com/kjfs/dieffe_sensor/internal/models"
	"github.com/kjfs/dieffe_sensor/internal/render"
	"github.com/kjfs/dieffe_sensor/internal/repository"
	"github.com/kjfs/dieffe_sensor/internal/repository/dbrepo"
	"golang.org/x/crypto/bcrypt"
)

// The repository used by the handlers.
var Repo *Repository // Repository Pattern: This Variable will use the Appconfig. Capital Letter, wich means this will be visible out side of this pakcage.

// is the repository type
type Repository struct {
	App *config.AppConfig
	Db  repository.DatabaseRepo // Repository pattern
} // Repository Pattern: We import the Appconfig from config into handlers.

// Creates a new repository.
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		Db:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
} // Repository Pattern: We have to use Repo Variable:

// Creates new repo for unit tests.
func NewTestingRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
		Db:  dbrepo.NewTestingRepo(a),
	}
}

// It sets the Repository for the handlers.
func NewHandlers(r *Repository) {
	Repo = r
} // Repository Pattern: r sets a Variable.

// Home is the handler for the home page
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	render.Template(w, "home.page.tmpl", &models.TemplateData{}, r)
}

// About is the handler for the about page
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	render.Template(w, "about.page.tmpl", &models.TemplateData{}, r)

}

// Renders the radioactivity page
func (m *Repository) Radioactivity(w http.ResponseWriter, r *http.Request) {
	if !m.App.Session.Exists(r.Context(), "user_id") {
		m.App.Session.Put(r.Context(), "flash", "Please log in with your credentials.")
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}
	render.Template(w, "radioactivity.page.tmpl", &models.TemplateData{}, r)
}

// Renders the temperature page
func (m *Repository) Temperature(w http.ResponseWriter, r *http.Request) {
	if !m.App.Session.Exists(r.Context(), "user_id") {
		m.App.Session.Put(r.Context(), "flash", "Please log in with your credentials.")
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}
	render.Template(w, "temperature.page.tmpl", &models.TemplateData{}, r)
}

// Renders the contact page
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.Template(w, "contact.page.tmpl", &models.TemplateData{}, r)
}

// Renders the missing - permission page
func (m *Repository) MissingPermissions(w http.ResponseWriter, r *http.Request) {
	if !m.App.Session.Exists(r.Context(), "user_id") {
		m.App.Session.Put(r.Context(), "flash", "Please log in with your credentials.")
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}
	render.Template(w, "missing-permissions.page.tmpl", &models.TemplateData{}, r)
}

// Renders the login page
func (m *Repository) ShowLogin(w http.ResponseWriter, r *http.Request) {
	render.Template(w, "login.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
	}, r)

}

// Renders the login page
func (m *Repository) ShowLoginPost(w http.ResponseWriter, r *http.Request) {
	err := m.App.Session.RenewToken(r.Context()) // Prevents Session Fixation Attacks. good for LogIn and LogOut
	if err != nil {
		m.App.ErrorLog.Println("could not renew token: ", err)
		return
	}
	err = r.ParseForm()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse form")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	form := forms.New(r.PostForm)
	form.Required("email", "password")
	form.IsEmail("email")
	form.MinLength("password", 12)
	form.MinLength("email", 5)
	if !form.Valid() {
		m.App.Session.Put(r.Context(), "error", "Your login credentials are invalid")
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}
	// Authenticate Method
	email := r.Form.Get("email")
	passwd := r.Form.Get("password")
	userId, _, err := m.Db.Authenticate(email, passwd)
	if err != nil {
		//log.Println("error during authentification: ", err)
		m.App.Session.Put(r.Context(), "error", "Your login credentials are invalid")
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}
	_, accessLevel, err := m.Db.AuthenticateAdmin(email, passwd)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "account access level issue")
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}
	m.App.Session.Put(r.Context(), "flash", "Well done - You are logged in")
	m.App.Session.Put(r.Context(), "user_id", userId)
	m.App.Session.Put(r.Context(), "access_lvl", accessLevel)
	render.Template(w, "login.page.tmpl", &models.TemplateData{
		Form: form,
	}, r)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Logs a user out.
func (m *Repository) Logout(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.Destroy(r.Context())
	_ = m.App.Session.RenewToken(r.Context())

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// AdminDashboard shows admin dashboard
func (m *Repository) AdminDashboard(w http.ResponseWriter, r *http.Request) {
	if m.App.Session.GetInt(r.Context(), "access_lvl") != 2 {
		m.App.Session.Put(r.Context(), "flash", "missing permissions")
		http.Redirect(w, r, "/missing-permissions", http.StatusSeeOther)
		return
	}
	render.Template(w, "admin-dashboard.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
	}, r)
}

// AdminRegistrationNew shows new registrations in admin tool
func (m *Repository) AdminRegistrationNew(w http.ResponseWriter, r *http.Request) {
	newUsers, err := m.Db.GetNewUsers()
	if err != nil {
		m.App.ErrorLog.Println("can't get new registrations from db", err)
		return
	}

	data := make(map[string]interface{})
	data["registration"] = newUsers

	render.Template(w, "admin-registrations-new.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	}, r)
}

// AdminRegistrationAll shows all registrations in admin tool
func (m *Repository) AdminRegistrationAll(w http.ResponseWriter, r *http.Request) {
	allUsers, err := m.Db.GetAllUsers()
	if err != nil {
		m.App.ErrorLog.Println("can't get user data from db", err)
		return
	}

	data := make(map[string]interface{})
	data["registration"] = allUsers

	render.Template(w, "admin-registrations-all.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	}, r)
}

//AdminShowRegistration
func (m *Repository) AdminShowRegistration(w http.ResponseWriter, r *http.Request) {
	exploded := strings.Split(r.RequestURI, "/")
	id, err := strconv.Atoi(exploded[4])
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	source := exploded[3]
	stringMap := make(map[string]string)
	stringMap["src"] = source

	user, err := m.Db.GetUserByID(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	data := make(map[string]interface{})
	data["registration"] = user
	render.Template(w, "admin-registrations-show-one-entry.page.tmpl", &models.TemplateData{
		Form:      forms.New(nil),
		StringMap: stringMap,
		Data:      data,
	}, r)
}

// AdminShowRegistrationPost updates registrations.
func (m *Repository) AdminShowRegistrationPost(w http.ResponseWriter, r *http.Request) {
	exploded := strings.Split(r.RequestURI, "/")
	id, err := strconv.Atoi(exploded[4])
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	source := exploded[3]
	stringMap := make(map[string]string)
	stringMap["src"] = source
	err = r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	user, err := m.Db.GetUserByID(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	accessLevelValue, err := strconv.Atoi(r.Form.Get("access_level"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	processedValue, err := strconv.Atoi(r.Form.Get("processed"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	user.FirstName = r.Form.Get("first_name")
	user.LastName = r.Form.Get("last_name")
	user.UserName = r.Form.Get("nick_name")
	user.Phone = r.Form.Get("phone_number")
	user.Access_lvl = accessLevelValue
	user.Email = r.Form.Get("email_address")
	user.Processed = processedValue

	_, err = m.Db.UpdateUser(id, user)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	m.App.Session.Put(r.Context(), "flash", "Changes saved")
	http.Redirect(w, r, fmt.Sprintf("/admin/registrations-%s", source), http.StatusSeeOther)

}

// AdminProcessRegistration marks a registrations as processed.
func (m *Repository) AdminProcessUser(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	source := chi.URLParam(r, "src")
	result, err := m.Db.UpdateProcessedVal(id, 1)
	log.Printf("Result: %v \n", result)
	if err != nil {
		helpers.ServerError(w, err)
		log.Printf("Error: %v \n", err)
		return
	}

	// Create notification email to send to user once registration request is approved by admin.
	htmlMessage := fmt.Sprintln(`

		<strong>Dieffe IoT - Welcome In The Club</strong><br>
		<p>
		You Request is approved!
		<p>
		Have fun consuming and generating data: https://www.iot-data-stream.de/user/login<br>
		<p>
		And btw, please let us know if you find any bugs (= errors). We reward you in advance with this song: https://youtu.be/05b9FM7BHV0
		<p>
		Thanks a lot, <br>
		Admin<br> 

	`)

	email, err := m.Db.GetUserEmail(id)
	if err != nil {
		helpers.ServerError(w, err)
		log.Printf("Error: %v \n", err)
		return
	}

	msg := models.MailData{
		To:      email,
		From:    "dieffesensoren@gmail.com",
		Subject: "Dieffe IoT - Welcome In The Club",
		Content: htmlMessage,
	}
	m.App.MailChan <- msg

	m.App.Session.Put(r.Context(), "flash", "registration marked as processed")
	http.Redirect(w, r, fmt.Sprintf("/admin/registrations-%s", source), http.StatusSeeOther)

}

// AdminDeleteRegistration deletes a registration.
func (m *Repository) AdminDeleteUser(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	source := chi.URLParam(r, "src")
	result, err := m.Db.DeleteUser(id)
	log.Printf("Result: %v \n", result)
	if err != nil {
		helpers.ServerError(w, err)
		log.Printf("Error: %v \n", err)
		return
	}

	m.App.Session.Put(r.Context(), "flash", "registration deleted")
	http.Redirect(w, r, fmt.Sprintf("/admin/registrations-%s", source), http.StatusSeeOther)

}

/*
Hier sind die Anpassungen für die Dieffe Sensoren page
*/

// Renders the register page
func (m *Repository) UserRegistration(w http.ResponseWriter, r *http.Request) {
	render.Template(w, "registration.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
	}, r)
}

// RegisterPost handles the posting of a Registration Form.
func (m *Repository) UserRegistrationPost(w http.ResponseWriter, r *http.Request) {
	user := models.User{
		Access_lvl: 0,
	}

	err := r.ParseForm()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "parsing error: please notify admin.")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	passwordBytes := []byte(r.Form.Get("password"))
	hashedPasswordBytes, err := bcrypt.GenerateFromPassword(passwordBytes, 12)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "security error. please notify admin.")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	user.FirstName = r.Form.Get("first_name")
	user.LastName = r.Form.Get("last_name")
	user.UserName = r.Form.Get("nick_name")
	user.Phone = r.Form.Get("phone_number")
	user.Email = r.Form.Get("email_address")
	user.Passwd = string(hashedPasswordBytes)

	form := forms.New(r.PostForm) // This will create the form object and send it back to our handlers.
	form.Required("first_name", "last_name", "nick_name", "phone_number", "email_address", "password")
	form.MinLength("first_name", 3)
	form.MinLength("last_name", 3)
	form.MinLength("nick_name", 3)
	form.MinLength("phone_number", 3)
	form.MinLength("email_address", 5)
	form.MinLength("password", 12)
	form.IsEmail("email_address")
	if !form.Valid() {
		data := make(map[string]interface{})
		data["registrationData"] = user
		render.Template(w, "registration.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		}, r)
		return
	} // If !form.Valid() then function will return an false, which we ignore. But it's adding an error if first_name is empty or non existing.

	newUserID, err := m.Db.InsertUser(user)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "database error. please contact admin.")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Create notification email to send to user once he filled out the form.
	htmlMessage := fmt.Sprintf(`

		<strong>Dieffe IoT - Registration Confirmation</strong><br>
		<p>
		Dear %s, <br>
		<p>
		This is to confirm your registration request. <br>
		Once this registration is approved, you will be notified via email.<br>
		In the meantime: Listen, enjoy and get inspired for your next haircut - https://youtu.be/q4xKvHANqjk
		<p>
		Thanks a lot, <br>
		Admin<br> 

	`, user.FirstName)

	msg := models.MailData{
		To:      user.Email,
		From:    "dieffesensoren@gmail.com",
		Subject: "Dieffe IoT - Registration Confirmation Email",
		Content: htmlMessage,
	}
	m.App.MailChan <- msg
	// Create notification to site owner!
	htmlMessage = fmt.Sprintf(`
		<strong>Registration Received</strong><br>
		<p>
		Dear Karl, <br>
		We have received a registration email. <br>
		Please find below the data. <br>
		Registration id %s<br>
		FirstName %s<br>
		LastName %s<br>
		UserName %s<br>
		Phone %s<br>
		Email %s<br>
		Password %s<br>
		<p>
		Please approve: https://www.iot-data-stream.de/admin/registrations-new<br>
		<p>
		Let's have fun!`, strconv.Itoa(newUserID), user.FirstName, user.LastName, user.UserName, user.Phone, user.Email, user.Passwd)
	msg = models.MailData{
		To:      "dieffesensoren@gmail.com",
		From:    "dieffesensoren@gmail.com",
		Subject: "Registration Confirmation Email",
		Content: htmlMessage,
	}
	m.App.MailChan <- msg
	// Create notification based on Email templates
	// htmlMessage = fmt.Sprintf(`
	// 	<strong>Registration Received</strong><br>
	// 	Dear Karl, <br>
	// 	We have received a registration email. <br>
	// 	Please find below the data. <br>
	// 	RegistrationId %s<br>
	// 	FirstName %s<br>
	// 	LastName %s<br>
	// 	UserName %s<br>
	// 	Phone %s<br>
	// 	Email %s<br>
	// 	Password %s <br>`, strconv.Itoa(newUserID), user.FirstName, user.LastName, user.UserName, user.Phone, user.Email, user.Passwd)

	// msg = models.MailData{
	// 	To:       "dieffesensoren@gmail.com",
	// 	From:     "dieffesensoren@gmail.com",
	// 	Subject:  "Registration Confirmation Email",
	// 	Content:  htmlMessage,
	// 	Template: "basic.html",
	// }
	// m.App.MailChan <- msg
	m.App.Session.Put(r.Context(), "user", user)
	http.Redirect(w, r, "/registration-summary", http.StatusSeeOther)
}

// Renders the registration-summary page
func (m *Repository) UserRegistrationSummary(w http.ResponseWriter, r *http.Request) {
	user, ok := m.App.Session.Get(r.Context(), "user").(models.User)
	if !ok {
		m.App.Session.Put(r.Context(), "warning", "For security reasons we redirect you.")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	m.App.Session.Remove(r.Context(), "user")
	data := make(map[string]interface{})
	data["user"] = user
	render.Template(w, "registration-summary.page.tmpl", &models.TemplateData{
		Data: data,
	}, r)
}
