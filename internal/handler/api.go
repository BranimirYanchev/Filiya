package handler

import (
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

func respondError(err error) (int, string) {
	if strings.Contains(err.Error(), "external") {
		return http.StatusBadRequest, err.Error()[10:]
	} else if strings.Contains(err.Error(), "internal") {
		log.Debugf(err.Error())
		return http.StatusInternalServerError, err.Error()[10:]
	}
	return http.StatusBadRequest, err.Error()
}
