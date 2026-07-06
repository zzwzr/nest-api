package admin

import (
	"context"

	"nest-api/app/workspace"
	"nest-api/internal/database"
	"nest-api/internal/ent"
	"nest-api/internal/ent/user"
	entworkspace "nest-api/internal/ent/workspace"
	"nest-api/internal/ent/workspacemember"
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
		Order(ent.Desc(entworkspace.FieldID)).
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

	ws, err := database.DB.Workspace.Get(ctx, workspaceID)
	if err != nil {
		if ent.IsNotFound(err) {
			return bizerr.New("工作空间不存在")
		}
		return err
	}
	if ws.OwnerID == ownerID {
		return nil
	}

	tx, err := database.DB.Tx(ctx)
	if err != nil {
		return err
	}

	_, err = tx.Workspace.
		UpdateOneID(workspaceID).
		SetOwnerID(ownerID).
		Save(ctx)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	if err = syncOwnerMember(ctx, tx, workspaceID, ws.OwnerID, ownerID); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}

func syncOwnerMember(ctx context.Context, tx *ent.Tx, workspaceID, oldOwnerID, newOwnerID int64) error {
	oldMember, err := tx.WorkspaceMember.
		Query().
		Where(
			workspacemember.WorkspaceIDEQ(workspaceID),
			workspacemember.UserIDEQ(oldOwnerID),
		).
		Only(ctx)
	if err == nil {
		if _, err = tx.WorkspaceMember.UpdateOneID(oldMember.ID).SetRole(workspace.RoleAdmin).Save(ctx); err != nil {
			return err
		}
	} else if !ent.IsNotFound(err) {
		return err
	}

	newMember, err := tx.WorkspaceMember.
		Query().
		Where(
			workspacemember.WorkspaceIDEQ(workspaceID),
			workspacemember.UserIDEQ(newOwnerID),
		).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			_, err = tx.WorkspaceMember.
				Create().
				SetWorkspaceID(workspaceID).
				SetUserID(newOwnerID).
				SetRole(workspace.RoleOwner).
				Save(ctx)
			return err
		}
		return err
	}

	_, err = tx.WorkspaceMember.UpdateOneID(newMember.ID).SetRole(workspace.RoleOwner).Save(ctx)
	return err
}
