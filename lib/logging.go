package lib

import (
	"io/ioutil"
	"log"
	"os"
)

func SetQuiet(q bool) {
	if q {
		log.SetOutput(ioutil.Discard)
	} else {
		log.SetOutput(os.Stdout)
	}
}
