package member

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"

	appworkspace "nest-api/app/workspace"
	"nest-api/internal/database"
	"nest-api/internal/ent"
	"nest-api/internal/ent/workspace"
	"nest-api/internal/ent/workspacemember"
	"nest-api/internal/runtime"
	bizerr "nest-api/pkg/errors"
)

const inviteCodeCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const inviteCodeLength = 6

func generateInviteCode() (string, error) {
	b := make([]byte, inviteCodeLength)
	max := big.NewInt(int64(len(inviteCodeCharset)))
	for i := range b {
		n, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}
		b[i] = inviteCodeCharset[n.Int64()]
	}
	return string(b), nil
}

func ensureInviteCode(ctx context.Context, workspaceID int64) (string, error) {
	ws, err := database.DB.Workspace.
		Query().
		Where(workspace.IDEQ(workspaceID)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return "", bizerr.New("工作空间不存在")
		}
		return "", err
	}
	if strings.TrimSpace(ws.InviteCode) != "" {
		return ws.InviteCode, nil
	}
	return assignInviteCode(ctx, workspaceID)
}

func assignInviteCode(ctx context.Context, workspaceID int64) (string, error) {
	for range 8 {
		code, err := generateInviteCode()
		if err != nil {
			return "", err
		}

		exists, err := database.DB.Workspace.
			Query().
			Where(workspace.InviteCodeEQ(code)).
			Exist(ctx)
		if err != nil {
			return "", err
		}
		if exists {
			continue
		}

		_, err = database.DB.Workspace.
			UpdateOneID(workspaceID).
			SetInviteCode(code).
			Save(ctx)
		if err != nil {
			return "", err
		}
		return code, nil
	}
	return "", bizerr.New("生成邀请码失败，请重试")
}

func buildInviteURL(code string) string {
	base := strings.TrimRight(runtime.SiteURL(), "/")
	if base == "" {
		base = "http://localhost:5173"
	}
	return fmt.Sprintf("%s/invite?inviteCode=%s", base, code)
}

func (Service) GetInviteLink(ctx context.Context, userID int64, params InviteLinkRequest) (*InviteLinkResponse, error) {
	if _, err := appworkspace.Require(ctx, userID, params.WorkspaceID, appworkspace.ActionMemberInvite); err != nil {
		return nil, err
	}

	ws, err := database.DB.Workspace.
		Query().
		Where(workspace.IDEQ(params.WorkspaceID)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, bizerr.New("工作空间不存在")
		}
		return nil, err
	}

	code, err := ensureInviteCode(ctx, params.WorkspaceID)
	if err != nil {
		return nil, err
	}

	return &InviteLinkResponse{
		WorkspaceID:   ws.ID,
		WorkspaceName: ws.Name,
		InviteCode:    code,
		InviteURL:     buildInviteURL(code),
		SiteURL:       runtime.SiteURL(),
	}, nil
}

func (Service) RefreshInviteLink(ctx context.Context, userID int64, params InviteLinkRequest) (*InviteLinkResponse, error) {
	if _, err := appworkspace.Require(ctx, userID, params.WorkspaceID, appworkspace.ActionMemberInvite); err != nil {
		return nil, err
	}

	ws, err := database.DB.Workspace.
		Query().
		Where(workspace.IDEQ(params.WorkspaceID)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, bizerr.New("工作空间不存在")
		}
		return nil, err
	}

	code, err := assignInviteCode(ctx, params.WorkspaceID)
	if err != nil {
		return nil, err
	}

	return &InviteLinkResponse{
		WorkspaceID:   ws.ID,
		WorkspaceName: ws.Name,
		InviteCode:    code,
		InviteURL:     buildInviteURL(code),
		SiteURL:       runtime.SiteURL(),
	}, nil
}

func (Service) PreviewInvite(ctx context.Context, params InvitePreviewRequest) (*InvitePreviewResponse, error) {
	code := strings.TrimSpace(params.InviteCode)
	if code == "" {
		return nil, bizerr.New("邀请码无效")
	}

	ws, err := database.DB.Workspace.
		Query().
		Where(workspace.InviteCodeEQ(code)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, bizerr.New("邀请链接无效或已过期")
		}
		return nil, err
	}

	return &InvitePreviewResponse{
		WorkspaceID:   ws.ID,
		WorkspaceName: ws.Name,
		InviteCode:    code,
	}, nil
}

func (Service) AcceptInvite(ctx context.Context, userID int64, params AcceptInviteRequest) (*AcceptInviteResponse, error) {
	code := strings.TrimSpace(params.InviteCode)
	if code == "" {
		return nil, bizerr.New("邀请码无效")
	}

	ws, err := database.DB.Workspace.
		Query().
		Where(workspace.InviteCodeEQ(code)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, bizerr.New("邀请链接无效或已过期")
		}
		return nil, err
	}

	exists, err := database.DB.WorkspaceMember.
		Query().
		Where(
			workspacemember.WorkspaceIDEQ(ws.ID),
			workspacemember.UserIDEQ(userID),
		).
		Exist(ctx)
	if err != nil {
		return nil, err
	}
	if exists {
		return &AcceptInviteResponse{
			WorkspaceID:   ws.ID,
			WorkspaceName: ws.Name,
			AlreadyMember: true,
		}, nil
	}

	_, err = database.DB.WorkspaceMember.
		Create().
		SetWorkspaceID(ws.ID).
		SetUserID(userID).
		SetRole(appworkspace.RoleEditor).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	return &AcceptInviteResponse{
		WorkspaceID:   ws.ID,
		WorkspaceName: ws.Name,
		AlreadyMember: false,
	}, nil
}
