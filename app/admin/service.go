package admin

import (
	"context"

	"nest-api/internal/database"
	"nest-api/internal/ent"
	"nest-api/internal/ent/user"
	"nest-api/internal/ent/workspace"
	"nest-api/internal/utils"
	bizerr "nest-api/pkg/errors"
)

type Service struct{}

func (Service) ListUsers(ctx context.Context) ([]UserItem, error) {
	users, err := database.DB.User.
		Query().
		Order(ent.Desc(user.FieldID)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	items := make([]UserItem, 0, len(users))
	for _, u := range users {
		items = append(items, UserItem{
			ID:      u.ID,
			Name:    u.Name,
			Account: u.Account,
			Email:   u.Email,
			Avatar:  u.Avatar,
			IsAdmin: u.IsAdmin,
			Status:  u.Status,
		})
	}

	return items, nil
}

func (Service) ListWorkspaces(ctx context.Context) ([]WorkspaceItem, error) {
	workspaces, err := database.DB.Workspace.
		Query().
		WithOwner().
		Order(ent.Desc(workspace.FieldID)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	items := make([]WorkspaceItem, 0, len(workspaces))
	for _, ws := range workspaces {
		ownerName := ""
		if ws.Edges.Owner != nil {
			ownerName = ws.Edges.Owner.Name
			if ownerName == "" {
				ownerName = ws.Edges.Owner.Account
			}
		}

		items = append(items, WorkspaceItem{
			ID:        ws.ID,
			Name:      ws.Name,
			OwnerID:   ws.OwnerID,
			OwnerName: ownerName,
			CreatedAt: ws.CreatedAt.Format(utils.DateTimeFormat),
		})
	}

	return items, nil
}

func (Service) TransferWorkspace(ctx context.Context, workspaceID, ownerID int64) error {
	exists, err := database.DB.User.
		Query().
		Where(user.IDEQ(ownerID)).
		Exist(ctx)
	if err != nil {
		return err
	}
	if !exists {
		return bizerr.New("目标用户不存在")
	}

	_, err = database.DB.Workspace.
		UpdateOneID(workspaceID).
		SetOwnerID(ownerID).
		Save(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return bizerr.New("工作空间不存在")
		}
		return err
	}

	return nil
}
