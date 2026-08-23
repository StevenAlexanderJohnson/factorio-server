package models

type CommandRequest struct {
	Command string `json:"command"`
}

type CommandResponse struct {
	Response string `json:"response"`
}
