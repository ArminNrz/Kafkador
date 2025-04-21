package api

type CreateBikerRequest struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phoneNumber"`
}

type CreateBikerResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	PhoneNumber string `json:"phoneNumber"`
}
