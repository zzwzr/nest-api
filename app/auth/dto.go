package auth

type RegisterRequest struct {
	Name            string `json:"name" binding:"required,min=1,max=50"`
	Account         string `json:"account" binding:"required,min=3,max=50"`
	Email           string `json:"email" binding:"required,email,max=100"`
	Password        string `json:"password" binding:"required,min=6"`
	ConfirmPassword string `json:"confirm_password" binding:"required,min=6"`
}

type LoginRequest struct {
	Account  string `json:"account" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6"`
}

type UserProfile struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Account  string `json:"account"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
	IsAdmin  bool   `json:"is_admin"`
	Status   int8   `json:"status"`
}

type LoginResponse struct {
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
	User         UserProfile `json:"user"`
}

type SiteResponse struct {
	Installed bool   `json:"installed"`
	SiteURL   string `json:"site_url"`
}
