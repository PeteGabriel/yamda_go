package models

//ErrorProblem represents the specification rfc7807
//that deals in a standardized way with errors in HTTP APIs.
type ErrorProblem struct {
	Type     string `json:"type,omitempty"`
	Title    string `json:"title"`
	Status   int    `json:"status"`
	Detail   string `json:"detail"`
	Instance string `json:"instance,omitempty"`
}
