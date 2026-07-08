package interfaces

import (
	"context"
	"strings"

	"nest-api/app/workspace"
	"nest-api/internal/database"
	"nest-api/internal/ent"
	entapi "nest-api/internal/ent/api"
	entfolder "nest-api/internal/ent/folder"
	entproject "nest-api/internal/ent/project"
	"nest-api/internal/utils"
	bizerr "nest-api/pkg/errors"
)

type Service struct{}

func (Service) List(ctx context.Context, userID int64, params ListRequest) ([]Item, error) {
	if _, err := workspace.Require(ctx, userID, params.WorkspaceID, workspace.ActionProjectRead); err != nil {
		return nil, err
	}
	if err := ensureProject(ctx, params.WorkspaceID, params.ProjectID); err != nil {
		return nil, err
	}

	rows, err := database.DB.API.
		Query().
		Where(entapi.ProjectIDEQ(params.ProjectID)).
		WithFolder().
		WithUpdater().
		Order(ent.Desc(entapi.FieldID)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	items := make([]Item, 0, len(rows))
	for _, row := range rows {
		folderName := ""
		if row.Edges.Folder != nil {
			folderName = row.Edges.Folder.Name
		}
		items = append(items, Item{
			ID:            row.ID,
			ProjectID:     row.ProjectID,
			FolderID:      row.FolderID,
			Name:          row.Name,
			Method:        row.Method,
			URL:           row.URL,
			Status:        row.Status,
			FolderName:    folderName,
			UpdatedBy:     row.UpdatedBy,
			UpdatedByName: workspace.UserDisplayName(row.Edges.Updater),
			CreatedAt:     row.CreatedAt.Format(utils.DateTimeFormat),
			UpdatedAt:     row.UpdatedAt.Format(utils.DateTimeFormat),
		})
	}
	return items, nil
}

func (Service) Create(ctx context.Context, userID int64, params CreateRequest) (int64, error) {
	if _, err := workspace.Require(ctx, userID, params.WorkspaceID, workspace.ActionProjectCreate); err != nil {
		return 0, err
	}
	if err := ensureProject(ctx, params.WorkspaceID, params.ProjectID); err != nil {
		return 0, err
	}
	if err := ensureFolder(ctx, params.ProjectID, params.FolderID); err != nil {
		return 0, err
	}

	method := strings.ToUpper(params.Method)
	status := params.Status
	if status == 0 {
		status = 1
	}

	row, err := database.DB.API.
		Create().
		SetProjectID(params.ProjectID).
		SetFolderID(params.FolderID).
		SetName(params.Name).
		SetMethod(method).
		SetURL(params.URL).
		SetStatus(status).
		SetCreatedBy(userID).
		SetUpdatedBy(userID).
		Save(ctx)
	if err != nil {
		return 0, err
	}
	return row.ID, nil
}

func (Service) Update(ctx context.Context, userID int64, params UpdateRequest) error {
	if _, err := workspace.Require(ctx, userID, params.WorkspaceID, workspace.ActionProjectUpdate); err != nil {
		return err
	}
	if err := ensureProject(ctx, params.WorkspaceID, params.ProjectID); err != nil {
		return err
	}
	if err := ensureInterface(ctx, params.ProjectID, params.InterfaceID); err != nil {
		return err
	}

	status := params.Status
	if status == 0 {
		status = 2
	}

	_, err := database.DB.API.
		UpdateOneID(params.InterfaceID).
		SetName(params.Name).
		SetMethod(strings.ToUpper(params.Method)).
		SetURL(params.URL).
		SetStatus(status).
		SetUpdatedBy(userID).
		Save(ctx)
	return err
}

func (Service) Delete(ctx context.Context, userID int64, params DeleteRequest) error {
	if _, err := workspace.Require(ctx, userID, params.WorkspaceID, workspace.ActionProjectDelete); err != nil {
		return err
	}
	if err := ensureProject(ctx, params.WorkspaceID, params.ProjectID); err != nil {
		return err
	}
	if err := ensureInterface(ctx, params.ProjectID, params.InterfaceID); err != nil {
		return err
	}

	return database.DB.API.DeleteOneID(params.InterfaceID).Exec(ctx)
}

func ensureProject(ctx context.Context, workspaceID, projectID int64) error {
	exists, err := database.DB.Project.
		Query().
		Where(
			entproject.IDEQ(projectID),
			entproject.WorkspaceIDEQ(workspaceID),
		).
		Exist(ctx)
	if err != nil {
		return err
	}
	if !exists {
		return bizerr.New("项目不存在")
	}
	return nil
}

func ensureFolder(ctx context.Context, projectID, folderID int64) error {
	exists, err := database.DB.Folder.
		Query().
		Where(
			entfolder.IDEQ(folderID),
			entfolder.ProjectIDEQ(projectID),
		).
		Exist(ctx)
	if err != nil {
		return err
	}
	if !exists {
		return bizerr.New("文件夹不存在")
	}
	return nil
}

func ensureInterface(ctx context.Context, projectID, interfaceID int64) error {
	exists, err := database.DB.API.
		Query().
		Where(
			entapi.IDEQ(interfaceID),
			entapi.ProjectIDEQ(projectID),
		).
		Exist(ctx)
	if err != nil {
		return err
	}
	if !exists {
		return bizerr.New("接口不存在")
	}
	return nil
}
