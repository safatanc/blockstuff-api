package auth

type Auth struct {
	UserID string `json:"user_id,omitempty"`
	Token  string `json:"token,omitempty"`
}
