package controller

type errorMessage struct {
	Message string
}

func errorResponse(message string) errorMessage {
	return errorMessage{Message: message}
}
