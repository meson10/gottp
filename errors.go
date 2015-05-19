package gottp

import (
	"net/http"

	utils "gopkg.in/simversity/gottp.v3/utils"
)

type HttpError struct {
	Status  int
	Message string
}

func (e HttpError) SendOverWire() utils.Q {
	if e.Status == 0 {
		e.Status = http.StatusInternalServerError
	}

	if len(e.Message) == 0 {
		e.Message = ERROR
	}

	return utils.Q{
		"data":    nil,
		"status":  e.Status,
		"message": e.Message,
	}
}
