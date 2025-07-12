package views

import (
	"log"
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
	if pErr, ok := err.(PublicError); ok {
		msg = pErr.Public()
	} else {
		log.Println("Public", err)
		msg = AlertMsgGeneric
	}

	d.Alert = &Alert{
		Level:   AlertLvlError,
		Message: msg,
	}
}

func (d *Data) AlertError(msg string) {
	d.Alert = &Alert{
		Level: AlertLvlError,
		Message: msg,
	}
}

type Data struct {
	Alert *Alert
	Yield interface{}
}

type Alert struct {
	Level   string
	Message string
}
