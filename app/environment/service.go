package environment

import (
	"context"

	"nest-api/app/workspace"
	"nest-api/internal/database"
	"nest-api/internal/ent"
	entenv "nest-api/internal/ent/environment"
	entproject "nest-api/internal/ent/project"
	"nest-api/internal/utils"
	bizerr "nest-api/pkg/errors"
)

type Service struct{}

func (Service) List(ctx context.Context, userID int64, params ListRequest) ([]Item, error) {
	if err := ensureProject(ctx, userID, params.WorkspaceID, params.ProjectID); err != nil {
		return nil, err
	}

	rows, err := database.DB.Environment.
		Query().
		Where(entenv.ProjectIDEQ(params.ProjectID)).
		Order(ent.Desc(entenv.FieldIsDefault), ent.Asc(entenv.FieldID)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	items := make([]Item, 0, len(rows))
	for _, row := range rows {
		items = append(items, Item{
			ID:        row.ID,
			ProjectID: row.ProjectID,
			Name:      row.Name,
			BaseURL:   row.BaseURL,
			IsDefault: row.IsDefault,
			CreatedAt: row.CreatedAt.Format(utils.DateTimeFormat),
		})
	}
	return items, nil
}

func (Service) Create(ctx context.Context, userID int64, params CreateRequest) error {
	if err := ensureProject(ctx, userID, params.WorkspaceID, params.ProjectID); err != nil {
		return err
	}
	if _, err := workspace.Require(ctx, userID, params.WorkspaceID, workspace.ActionProjectCreate); err != nil {
		return err
	}

	isDefault := params.IsDefault
	count, err := database.DB.Environment.
		Query().
		Where(entenv.ProjectIDEQ(params.ProjectID)).
		Count(ctx)
	if err != nil {
		return err
	}
	if count == 0 {
		isDefault = true
	}

	if isDefault {
		if err := clearDefault(ctx, params.ProjectID); err != nil {
			return err
		}
	}

	_, err = database.DB.Environment.
		Create().
		SetProjectID(params.ProjectID).
		SetName(params.Name).
		SetBaseURL(params.BaseURL).
		SetIsDefault(isDefault).
		SetCreatedBy(userID).
		Save(ctx)
	return err
}

func (Service) Update(ctx context.Context, userID int64, params UpdateRequest) error {
	if err := ensureProject(ctx, userID, params.WorkspaceID, params.ProjectID); err != nil {
		return err
	}
	if _, err := workspace.Require(ctx, userID, params.WorkspaceID, workspace.ActionProjectUpdate); err != nil {
		return err
	}
	if err := ensureEnvironment(ctx, params.ProjectID, params.EnvironmentID); err != nil {
		return err
	}

	if params.IsDefault {
		if err := clearDefault(ctx, params.ProjectID); err != nil {
			return err
		}
	}

	_, err := database.DB.Environment.
		UpdateOneID(params.EnvironmentID).
		SetName(params.Name).
		SetBaseURL(params.BaseURL).
		SetIsDefault(params.IsDefault).
		Save(ctx)
	return err
}

func (Service) Delete(ctx context.Context, userID int64, params DeleteRequest) error {
	if err := ensureProject(ctx, userID, params.WorkspaceID, params.ProjectID); err != nil {
		return err
	}
	if _, err := workspace.Require(ctx, userID, params.WorkspaceID, workspace.ActionProjectDelete); err != nil {
		return err
	}
	if err := ensureEnvironment(ctx, params.ProjectID, params.EnvironmentID); err != nil {
		return err
	}

	n, err := database.DB.Environment.
		Delete().
		Where(
			entenv.IDEQ(params.EnvironmentID),
			entenv.ProjectIDEQ(params.ProjectID),
		).
		Exec(ctx)
	if err != nil {
		return err
	}
	if n == 0 {
		return bizerr.New("环境不存在")
	}
	return nil
}

func ensureProject(ctx context.Context, userID, workspaceID, projectID int64) error {
	if _, err := workspace.Require(ctx, userID, workspaceID, workspace.ActionProjectRead); err != nil {
		return err
	}

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

func ensureEnvironment(ctx context.Context, projectID, environmentID int64) error {
	exists, err := database.DB.Environment.
		Query().
		Where(
			entenv.IDEQ(environmentID),
			entenv.ProjectIDEQ(projectID),
		).
		Exist(ctx)
	if err != nil {
		return err
	}
	if !exists {
		return bizerr.New("环境不存在")
	}
	return nil
}

func clearDefault(ctx context.Context, projectID int64) error {
	_, err := database.DB.Environment.
		Update().
		Where(
			entenv.ProjectIDEQ(projectID),
			entenv.IsDefaultEQ(true),
		).
		SetIsDefault(false).
		Save(ctx)
	return err
}
