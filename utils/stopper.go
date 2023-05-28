package utils

type stoppable interface {
	Stop()
}

func StopIfNotNil(runned stoppable) {
	if runned != nil {
		runned.Stop()
	}
}
