package admin

type UserItem struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Account  string `json:"account"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
	IsAdmin  bool   `json:"is_admin"`
	Status   int8   `json:"status"`
}

type WorkspaceItem struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	OwnerID   int64  `json:"owner_id"`
	OwnerName string `json:"owner_name"`
	CreatedAt string `json:"created_at"`
}

type TransferWorkspaceRequest struct {
	OwnerID int64 `json:"owner_id" binding:"required,min=1"`
}
