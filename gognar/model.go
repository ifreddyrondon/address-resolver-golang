package gognar

type ResponseError struct {
	Message string `json:"message"`
	Errors  string `json:"error"`
	Status  int    `json:"status"`
}
