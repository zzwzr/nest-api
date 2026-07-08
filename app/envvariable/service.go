package envvariable

import (
	"context"

	"nest-api/app/workspace"
	"nest-api/internal/database"
	"nest-api/internal/ent"
	entenv "nest-api/internal/ent/environment"
	entenvvar "nest-api/internal/ent/environmentvariable"
	entproject "nest-api/internal/ent/project"
	"nest-api/internal/utils"
	bizerr "nest-api/pkg/errors"
)

type Service struct{}

func (Service) List(ctx context.Context, userID int64, params ListRequest) ([]Item, error) {
	if err := ensureEnvironment(ctx, userID, params.WorkspaceID, params.ProjectID, params.EnvironmentID); err != nil {
		return nil, err
	}

	rows, err := database.DB.EnvironmentVariable.
		Query().
		Where(entenvvar.EnvironmentIDEQ(params.EnvironmentID)).
		Order(ent.Asc(entenvvar.FieldKey)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	items := make([]Item, 0, len(rows))
	for _, row := range rows {
		items = append(items, Item{
			ID:            row.ID,
			EnvironmentID: row.EnvironmentID,
			Key:           row.Key,
			Value:         row.Value,
			Description:   row.Description,
			CreatedAt:     row.CreatedAt.Format(utils.DateTimeFormat),
			UpdatedAt:     row.UpdatedAt.Format(utils.DateTimeFormat),
		})
	}
	return items, nil
}

func (Service) Create(ctx context.Context, userID int64, params CreateRequest) error {
	if err := ensureEnvironment(ctx, userID, params.WorkspaceID, params.ProjectID, params.EnvironmentID); err != nil {
		return err
	}
	if _, err := workspace.Require(ctx, userID, params.WorkspaceID, workspace.ActionProjectCreate); err != nil {
		return err
	}

	exists, err := database.DB.EnvironmentVariable.
		Query().
		Where(
			entenvvar.EnvironmentIDEQ(params.EnvironmentID),
			entenvvar.KeyEQ(params.Key),
		).
		Exist(ctx)
	if err != nil {
		return err
	}
	if exists {
		return bizerr.New("变量名已存在")
	}

	_, err = database.DB.EnvironmentVariable.
		Create().
		SetEnvironmentID(params.EnvironmentID).
		SetKey(params.Key).
		SetValue(params.Value).
		SetDescription(params.Description).
		SetCreatedBy(userID).
		Save(ctx)
	return err
}

func (Service) Update(ctx context.Context, userID int64, params UpdateRequest) error {
	if err := ensureEnvironment(ctx, userID, params.WorkspaceID, params.ProjectID, params.EnvironmentID); err != nil {
		return err
	}
	if _, err := workspace.Require(ctx, userID, params.WorkspaceID, workspace.ActionProjectUpdate); err != nil {
		return err
	}
	if err := ensureVariable(ctx, params.EnvironmentID, params.VariableID); err != nil {
		return err
	}

	exists, err := database.DB.EnvironmentVariable.
		Query().
		Where(
			entenvvar.EnvironmentIDEQ(params.EnvironmentID),
			entenvvar.KeyEQ(params.Key),
			entenvvar.IDNEQ(params.VariableID),
		).
		Exist(ctx)
	if err != nil {
		return err
	}
	if exists {
		return bizerr.New("变量名已存在")
	}

	_, err = database.DB.EnvironmentVariable.
		UpdateOneID(params.VariableID).
		SetKey(params.Key).
		SetValue(params.Value).
		SetDescription(params.Description).
		Save(ctx)
	return err
}

func (Service) Delete(ctx context.Context, userID int64, params DeleteRequest) error {
	if err := ensureEnvironment(ctx, userID, params.WorkspaceID, params.ProjectID, params.EnvironmentID); err != nil {
		return err
	}
	if _, err := workspace.Require(ctx, userID, params.WorkspaceID, workspace.ActionProjectDelete); err != nil {
		return err
	}
	if err := ensureVariable(ctx, params.EnvironmentID, params.VariableID); err != nil {
		return err
	}

	n, err := database.DB.EnvironmentVariable.
		Delete().
		Where(
			entenvvar.IDEQ(params.VariableID),
			entenvvar.EnvironmentIDEQ(params.EnvironmentID),
		).
		Exec(ctx)
	if err != nil {
		return err
	}
	if n == 0 {
		return bizerr.New("变量不存在")
	}
	return nil
}

func ensureEnvironment(ctx context.Context, userID, workspaceID, projectID, environmentID int64) error {
	if err := ensureProject(ctx, userID, workspaceID, projectID); err != nil {
		return err
	}

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

func ensureVariable(ctx context.Context, environmentID, variableID int64) error {
	exists, err := database.DB.EnvironmentVariable.
		Query().
		Where(
			entenvvar.IDEQ(variableID),
			entenvvar.EnvironmentIDEQ(environmentID),
		).
		Exist(ctx)
	if err != nil {
		return err
	}
	if !exists {
		return bizerr.New("变量不存在")
	}
	return nil
}
