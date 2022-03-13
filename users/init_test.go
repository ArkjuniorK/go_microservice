package users

import (
	"testing"

	"github.com/rs/xid"
)

// testing validation for User
// struct ini init.go
func TestValidate(t *testing.T) {
	id := xid.New().String()

	// create user
	u := &User{
		ID:       id,
		Username: "adfad12",
		Fullname: "ArkjuniorK Retvrnz",
		Email:    "arkjuniork@email.com",
		Password: "meong123",
	}

	// validate
	err := u.Validate()
	if err != nil {
		t.Fatal(err)
	}

}
