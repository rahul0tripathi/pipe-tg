package entity

type PipeMessage struct {
	ID        string `json:"id"`
	Message   string `json:"message"`
	Channel   string `json:"channel"`
	Timestamp int    `json:"timestamp"`
}
