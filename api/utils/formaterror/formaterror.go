package formaterror

import (
	"errors"
	"strings"
)

func FormatError(err string) error {

	if strings.Contains(err, "username") {
		return errors.New("Username already Exists")
	}

	if strings.Contains(err, "email") {
		return errors.New("Email Already Exists")
	}

	if strings.Contains(err, "title") {
		return errors.New("Title Already Exists")
	}

	if strings.Contains(err, "hashedPassword") {
		return errors.New("Password is Incorrect")
	}

	return errors.New("Provided Details are Incorrect")
}