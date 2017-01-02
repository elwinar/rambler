package main

import (
	"io/ioutil"

	"github.com/elwinar/rambler/log"
)

func init() {
	logger = log.NewLogger(func(l *log.Logger) {
		l.Output = ioutil.Discard
	})
}
