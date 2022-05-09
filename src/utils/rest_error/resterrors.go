package rest_error

import (
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strconv"
)

var logger *log.Logger

func init() {
	//Creating Root Log.Entry ------------------------------------------------
	logger = log.New()
	logger.SetLevel(log.DebugLevel) // DebugLevel for verbose logging
	logger.SetFormatter(&log.JSONFormatter{})
	hostname, err := os.Hostname()
	if err != nil {
		logger.Debugf("Error while trying to get host name, err = %v", err)
		hostname = "error"
	}
	pid := os.Getpid()
	entry := logger.WithFields(log.Fields{
		"hostname": hostname,
		"appname":  "url-shortener",
		"pid":      strconv.Itoa(pid),
	})
	main_entry := entry.WithFields(log.Fields{
		"package": "main",
	})
	main_entry.Debug("Into this world, we're thrown!")
}

// A RestErr is an error that is used when the required input fails validation.
// swagger:response RestErr
type RestErr interface {
	Message() string
	Status() int
	Error() string
	Causes() []interface{}
}

type restErr struct {
	ErrMessage string        `json:"message"`
	ErrStatus  int           `json:"status"`
	ErrError   string        `json:"error"`
	ErrCauses  []interface{} `json:"causes"`
}

func (e restErr) Error() string {
	return fmt.Sprintf("message: %s - status: %d - error: %s - causes: %v",
		e.ErrMessage, e.ErrStatus, e.ErrError, e.ErrCauses)
}

func (e restErr) Message() string {
	return e.ErrMessage
}

func (e restErr) Status() int {
	return e.ErrStatus
}

func (e restErr) Causes() []interface{} {
	return e.ErrCauses
}

func NewRestError(message string, status int, err string, causes []interface{}) RestErr {
	return restErr{
		ErrMessage: message,
		ErrStatus:  status,
		ErrError:   err,
		ErrCauses:  causes,
	}
}

func NewRestErrorFromBytes(bytes []byte) (RestErr, error) {
	var apiErr restErr
	if err := json.Unmarshal(bytes, &apiErr); err != nil {
		return nil, errors.New("invalid json")
	}
	return apiErr, nil
}

func NewBadRequestError(message string) RestErr {
	logger.Errorf(" %s", message)
	return restErr{
		ErrMessage: message,
		ErrStatus:  http.StatusBadRequest,
		ErrError:   "bad_request",
	}
}

func NewNotFoundError(message string) RestErr {
	logger.Errorf(" %s", message)
	return restErr{
		ErrMessage: message,
		ErrStatus:  http.StatusNotFound,
		ErrError:   "not_found",
	}
}

func NewUnauthorizedError(message string) RestErr {
	return restErr{
		ErrMessage: message,
		ErrStatus:  http.StatusUnauthorized,
		ErrError:   "unauthorized",
	}
}

func NewInternalServerError(message string, err error) RestErr {
	logger.Errorf(" %s", err)
	result := restErr{
		ErrMessage: message,
		ErrStatus:  http.StatusInternalServerError,
		ErrError:   "internal_server_error",
	}
	if err != nil {
		result.ErrCauses = append(result.ErrCauses, err.Error())
	}
	return result
}

func ShowLogError(err string) {
	logger.Errorln(err)
}
