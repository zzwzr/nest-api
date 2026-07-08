package folder

type TreeNode struct {
	ID        string     `json:"id"`
	ProjectID int64      `json:"project_id"`
	Name      string     `json:"name"`
	Type      string     `json:"type"`
	Method    string     `json:"method,omitempty"`
	Children  []TreeNode `json:"children,omitempty"`
}

type ProjectScopeRequest struct {
	WorkspaceID int64 `form:"workspace_id" json:"workspace_id" binding:"required,min=1"`
	ProjectID   int64 `form:"project_id" json:"project_id" binding:"required,min=1"`
}

type CreateRequest struct {
	WorkspaceID int64  `json:"workspace_id" binding:"required,min=1"`
	ProjectID   int64  `json:"project_id" binding:"required,min=1"`
	ParentID    int64  `json:"parent_id"`
	Name        string `json:"name" binding:"required,min=1,max=100"`
}

type UpdateRequest struct {
	WorkspaceID int64  `json:"workspace_id" binding:"required,min=1"`
	ProjectID   int64  `json:"project_id" binding:"required,min=1"`
	FolderID    int64  `json:"folder_id" binding:"required,min=1"`
	Name        string `json:"name" binding:"required,min=1,max=100"`
}

type DeleteRequest struct {
	WorkspaceID int64 `form:"workspace_id" json:"workspace_id" binding:"required,min=1"`
	ProjectID   int64 `form:"project_id" json:"project_id" binding:"required,min=1"`
	FolderID    int64 `form:"folder_id" json:"folder_id" binding:"required,min=1"`
}
