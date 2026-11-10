package common

import "log"

// HandleError panics if err is not nil.
func HandleError(err error) {
	if err != nil {
		log.Panic(err)
	}
}