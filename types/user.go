package types

type User struct {
	UserId    string `json:"user_id"`
	Email     string             `json:"email"`
	UserName  string             `json:"username"`
	FirstName string             `json:"first_name"`
	LastName  string             `json:"last_name"`
	Avatar    string             `json:"avatar"`
}

type UserCache struct {
	User        User   `json:"user"`
	AccessToken string `json:"access_token"`
}
