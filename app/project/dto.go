package project

type Item struct {
	ID          int64  `json:"id"`
	WorkspaceID int64  `json:"workspace_id"`
	Name        string `json:"name"`
	CreatedBy   int64  `json:"created_by"`
	CreatorName string `json:"creator_name,omitempty"`
	CreatedAt   string `json:"created_at"`
}

type ListRequest struct {
	WorkspaceID int64 `form:"workspace_id" binding:"required,min=1"`
}

type CreateRequest struct {
	WorkspaceID int64  `json:"workspace_id" binding:"required,min=1"`
	Name        string `json:"name" binding:"required,min=1,max=100"`
}

type UpdateRequest struct {
	WorkspaceID int64  `json:"workspace_id" binding:"required,min=1"`
	ProjectID   int64  `json:"project_id" binding:"required,min=1"`
	Name        string `json:"name" binding:"required,min=1,max=100"`
}

type DeleteRequest struct {
	WorkspaceID int64 `form:"workspace_id" json:"workspace_id" binding:"required,min=1"`
	ProjectID   int64 `form:"project_id" json:"project_id" binding:"required,min=1"`
}
