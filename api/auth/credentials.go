package auth

// Credentials struct models the structure of a user, both in the request body, and in the DB
type Credentials struct {
	Password string `json:"user_password", db:"user_password"`
	Username string `json:"user_name", db:"user_name"`
	Email    string `json:"user_email", db:"user_email"`
}
