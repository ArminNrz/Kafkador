package model

type CreateBikerMessage struct {
	CorrelationID string `json:"correlation_id"`
	Payload       any    `json:"payload"`
}

type CreateBikerResponseMessage struct {
	CorrelationID string `json:"correlation_id"`
	StatusCode    int64  `json:"status_code"`
	Message       string `json:"message"`
	ID            string `json:"id"`
	Name          string `json:"name"`
	PhoneNumber   string `json:"phoneNumber"`
}
