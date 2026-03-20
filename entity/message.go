package entity

type ResponseMsg struct {
	Response *Response
	Err      error
}

type SaveResultMsg struct {
	Err error
}

type CollectionsLoadedMsg struct {
	Collections []Collection
	Err         error
}
