package error

type Error struct {
	Code     string `json:"code"`
	Message  string `json:"message"`
	Internal error
	Metadata map[string]any `json:"metadata,omitempty"`
}

func (e *Error) Error() string {
	return e.Message
}
