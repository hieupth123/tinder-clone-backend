package user

type (
	DetailRequest struct {
		Uuid string `uri:"uuid"`
	}

	GetDetailResponse struct {
		Uuid      string   `json:"uuid"`
		LastName  string   `json:"lastName"`
		FirstName string   `json:"firstName"`
		Gender    string   `json:"gender"`
		Age       int      `json:"age"`
		Picture   string   `json:"picture"`
		Email     string   `json:"email"`
		Phone     string   `json:"phone"`
		Matches   []string `json:"matches"`
	}
)
