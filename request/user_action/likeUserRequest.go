package user_action
type (
	LikeUserRequest struct {
		UserUuid string `uri:"user_uuid"`
		GuestUuid string `uri:"guest_uuid"`
	}
)