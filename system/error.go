package system

import log "github.com/sirupsen/logrus"

func HandleError(err error) {
	if err != nil {
		log.Panic(err)
	}
}
