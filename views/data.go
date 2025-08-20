package views

import (
	"errors"
	"lenslocked.com/models"
)

const (
	AlertLvlError   = "danger"
	AlertLvlWarning = "warning"
	AlertLvlInfo    = "info"
	AlertLvlSuccess = "success"
	AlertMsgGeneric = "Something went wrong. Please try " +
		"again, and contact us if the problem persists."
)

type PublicError interface {
	error
	Public() string
}

func (d *Data) SetAlert(err error) {
	var msg string
	var pErr PublicError
	if errors.As(err, &pErr) {
		msg = pErr.Public()
	}

	d.Alert = &Alert{
		Level:   AlertLvlError,
		Message: msg,
	}
}

func (d *Data) AlertError(msg string) {
	d.Alert = &Alert{
		Level:   AlertLvlError,
		Message: msg,
	}
}

type Data struct {
	Alert *Alert
	User  *models.User
	Yield interface{}
}

type Alert struct {
	Level   string
	Message string
}
