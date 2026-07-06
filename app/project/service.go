package project

import (
	"context"

	"nest-api/app/workspace"
	"nest-api/internal/database"
	"nest-api/internal/ent"
	entproject "nest-api/internal/ent/project"
	"nest-api/internal/utils"
	bizerr "nest-api/pkg/errors"
)

type Service struct{}

func (Service) List(ctx context.Context, userID int64, params ListRequest) ([]Item, error) {
	if _, err := workspace.Require(ctx, userID, params.WorkspaceID, workspace.ActionProjectRead); err != nil {
		return nil, err
	}

	projects, err := database.DB.Project.
		Query().
		Where(entproject.WorkspaceIDEQ(params.WorkspaceID)).
		WithCreator().
		Order(ent.Desc(entproject.FieldID)).
		All(ctx)
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

	_, err := database.DB.Project.
		Create().
		SetWorkspaceID(params.WorkspaceID).
		SetName(params.Name).
		SetCreatedBy(userID).
		Save(ctx)
	return err
}

func (Service) Update(ctx context.Context, userID int64, params UpdateRequest) error {
	if _, err := workspace.Require(ctx, userID, params.WorkspaceID, workspace.ActionProjectUpdate); err != nil {
		return err
	}

	exists, err := database.DB.Project.
		Query().
		Where(
			entproject.IDEQ(params.ProjectID),
			entproject.WorkspaceIDEQ(params.WorkspaceID),
		).
		Exist(ctx)
	if err != nil {
		return err
	}
	if !exists {
		return bizerr.New("项目不存在")
	}

	_, err = database.DB.Project.
		UpdateOneID(params.ProjectID).
		SetName(params.Name).
		Save(ctx)
	return err
}

func (Service) Delete(ctx context.Context, userID int64, params DeleteRequest) error {
	if _, err := workspace.Require(ctx, userID, params.WorkspaceID, workspace.ActionProjectDelete); err != nil {
		return err
	}

	n, err := database.DB.Project.
		Delete().
		Where(
			entproject.IDEQ(params.ProjectID),
			entproject.WorkspaceIDEQ(params.WorkspaceID),
		).
		Exec(ctx)
	if err != nil {
		return err
	}
	if n == 0 {
		return bizerr.New("项目不存在")
	}
	return nil
}
