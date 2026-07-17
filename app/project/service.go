package project

import (
	"context"

	"nest-api/app/workspace"
	"nest-api/internal/utils"
	bizerr "nest-api/pkg/errors"
)

type Service struct{}

func (Service) List(ctx context.Context, userID int64, params ListRequest) ([]Item, error) {
	if _, err := workspace.Require(ctx, userID, params.WorkspaceID, workspace.ActionProjectRead); err != nil {
		return nil, err
	}

	projects, err := Repo{}.ListByWorkspace(ctx, params.WorkspaceID)
	if err != nil {
		return nil, err
	}

	items := make([]Item, 0, len(projects))
	for _, p := range projects {
		items = append(items, Item{
			ID:          p.ID,
			WorkspaceID: p.WorkspaceID,
			Name:        p.Name,
			CreatedBy:   p.CreatedBy,
			CreatorName: workspace.UserDisplayName(p.Edges.Creator),
			CreatedAt:   p.CreatedAt.Format(utils.DateTimeFormat),
		})
	}
	return items, nil
}

func (Service) Create(ctx context.Context, userID int64, params CreateRequest) error {
	if _, err := workspace.Require(ctx, userID, params.WorkspaceID, workspace.ActionProjectCreate); err != nil {
		return err
	}

	_, err := Repo{}.Create(ctx, params.WorkspaceID, userID, params.Name)
	return err
}

func (Service) Update(ctx context.Context, userID int64, params UpdateRequest) error {
	if _, err := workspace.Require(ctx, userID, params.WorkspaceID, workspace.ActionProjectUpdate); err != nil {
		return err
	}
	if err := EnsureExists(ctx, params.WorkspaceID, params.ProjectID); err != nil {
		return err
	}

	return (Repo{}).UpdateName(ctx, params.ProjectID, params.Name)
}

func (Service) Delete(ctx context.Context, userID int64, params DeleteRequest) error {
	if _, err := workspace.Require(ctx, userID, params.WorkspaceID, workspace.ActionProjectDelete); err != nil {
		return err
	}

	n, err := Repo{}.Delete(ctx, params.WorkspaceID, params.ProjectID)
	if err != nil {
		return err
	}
	if n == 0 {
		return bizerr.New("项目不存在")
	}
	return nil
}
