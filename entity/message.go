package entity

type ResponseMsg struct {
	Response *Response
	Err      error
}

type SaveResultMsg struct {
	Err error
}
