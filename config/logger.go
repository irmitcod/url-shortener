package config

import (
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
)

func NewLogger() *log.Entry {
	// Creating Root Log.Entry ------------------------------------------------
	logger := log.New()
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

	return entry
}
