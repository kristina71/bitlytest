package models

type NotFound struct {
}

func NotFoundError() error {
	return NotFound{}
}

func (m NotFound) Error() string {
	return "not found"
}

type BadRequest struct {
	message string
}

func BadRequestError(message string) error {
	return BadRequest{message: message}
}

func (m BadRequest) Error() string {
	return "bad request: " + m.message
}
