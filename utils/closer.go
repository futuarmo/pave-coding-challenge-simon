package utils

type closable interface {
	Close()
}

func CloseIfNotNil(opened closable) {
	if opened != nil {
		opened.Close()
	}
}
