package utils

// Claims represents the JWT payload used throughout the project.
// It includes the user identifier and any additional fields you may need.
// Extend this struct as required (e.g., roles, permissions, etc.).
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
}
