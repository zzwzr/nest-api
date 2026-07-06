package workspace

type Item struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	OwnerID   int64  `json:"owner_id"`
	OwnerName string `json:"owner_name,omitempty"`
	Role      uint8  `json:"role"`
	CreatedAt string `json:"created_at"`
}

type CreateRequest struct {
	Name string `json:"name" binding:"required,min=1,max=100"`
}

type GetRequest struct {
	WorkspaceID int64 `form:"workspace_id" binding:"required,min=1"`
}

type UpdateRequest struct {
	WorkspaceID int64  `json:"workspace_id" binding:"required,min=1"`
	Name        string `json:"name" binding:"required,min=1,max=100"`
}

type DeleteRequest struct {
	WorkspaceID int64 `form:"workspace_id" json:"workspace_id" binding:"required,min=1"`
}

type TransferOwnerRequest struct {
	WorkspaceID int64 `json:"workspace_id" binding:"required,min=1"`
	OwnerID     int64 `json:"owner_id" binding:"required,min=1"`
}
