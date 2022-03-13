package users

import (
	"encoding/json"
	"fmt"
	"io"
	"regexp"

	"github.com/go-playground/validator/v10"
	"github.com/rs/xid"
)

// structure of User
type User struct {
	ID       string `json:"_id" db:"_id" validate:"required"`
	Username string `json:"username" db:"username" validate:"required,username"`
	Fullname string `json:"fullname" db:"fullname" validate:"required"`
	Email    string `json:"email" db:"email" validate:"email"`
	Password string `json:"password" db:"password" validate:"min=8"`
}

// custom validate for username tag
// that return boolean if validation works well
// otherwise return false
func validateUsername(fl validator.FieldLevel) bool {
	// username format is like the other username pattern
	// satisfied => [admin123, admin-123, admin, 123admin, 123_admin_12, ...etc]
	// unsatisfied => [!admin123, admin*123, admin#, $123admin, ...etc]
	// so use regex to validate those chars
	re := regexp.MustCompile(`^(\w*[-._]?[\d|\w]+)$`)
	matches := re.FindString(fl.Field().String())

	// if len of matches is 0
	// then return false
	if len(matches) == 0 {
		return false
	}

	// otherwise return true
	return true
}

// validate struct with validate tag
// that attach to each key
func (u *User) Validate() error {
	validate := validator.New()

	// custom validation for username tag
	validate.RegisterValidation("username", validateUsername)

	// return validate
	return validate.Struct(u)
}

// function to change format of JSON
// to User using decoding function provided
// by encoding/json package then return the
// encode data of user
//
// it would write to http.Request
func (u *User) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(u)
}

// function to change format of User
// to JSON using encoding function provided
// by encoding/json package then return the
// encode data of user
//
// it would write to http.ResponseWriter
func (u *User) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(u)
}

// Users is collection of user
type Users []*User

// function to change format of Users
// to JSON uisng Encoding function provided
// by encoding/json package then return the
// encode data of the Users
// ----------------------------------------
func (u *Users) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(u)
}

// function to list Users
func ListUsers() Users {
	return userList
}

// function to create new User
func CreateUser(u *User) *User {
	// init new ID
	// then add user to users
	u.ID = xid.New().String()
	userList = append(userList, u)
	return u
}

// function to update User by given _id
func UpdateUser(_id string, u *User) (*User, error) {
	// get idx of user
	// by given _id
	idx, err := findUser(_id)
	if err != nil {
		return nil, err
	}

	// replace data on current idx
	u.ID = _id
	userList[idx] = u

	return u, nil
}

var ErrUserNotFound = fmt.Errorf("User not found")

func findUser(_id string) (int, error) {
	for i, u := range userList {
		if u.ID == _id {
			return i, nil
		}
	}

	return -1, ErrUserNotFound
}

// list of users
var userList = Users{}
