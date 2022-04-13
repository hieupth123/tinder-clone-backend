package user

type (
	GetListResponse struct {
		Uuid      string `json:"uuid"`
		LastName  string `json:"lastName"`
		FirstName string `json:"firstName"`
		Gender    string `json:"gender"`
		Age       int    `json:"age"`
		Picture   string `json:"picture"`
	}
)
