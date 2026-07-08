package folder

import (
	"context"
	"fmt"
	"sort"

	"nest-api/app/workspace"
	"nest-api/internal/database"
	"nest-api/internal/ent"
	entapi "nest-api/internal/ent/api"
	entfolder "nest-api/internal/ent/folder"
	entproject "nest-api/internal/ent/project"
	bizerr "nest-api/pkg/errors"
)

type Service struct{}

func (Service) Tree(ctx context.Context, userID int64, params ProjectScopeRequest) ([]TreeNode, error) {
	if err := requireProject(ctx, userID, params); err != nil {
		return nil, err
	}

	folders, err := database.DB.Folder.
		Query().
		Where(entfolder.ProjectIDEQ(params.ProjectID)).
		Order(ent.Asc(entfolder.FieldSortOrder), ent.Asc(entfolder.FieldID)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	ifaces, err := database.DB.API.
		Query().
		Where(entapi.ProjectIDEQ(params.ProjectID)).
		Order(ent.Asc(entapi.FieldSortOrder), ent.Asc(entapi.FieldID)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	return buildTree(folders, ifaces), nil
}

func (Service) Create(ctx context.Context, userID int64, params CreateRequest) (int64, error) {
	if _, err := workspace.Require(ctx, userID, params.WorkspaceID, workspace.ActionProjectCreate); err != nil {
		return 0, err
	}
	if err := ensureProject(ctx, params.WorkspaceID, params.ProjectID); err != nil {
		return 0, err
	}
	if params.ParentID > 0 {
		if err := ensureFolder(ctx, params.ProjectID, params.ParentID); err != nil {
			return 0, err
		}
	}

	row, err := database.DB.Folder.
		Create().
		SetProjectID(params.ProjectID).
		SetParentID(params.ParentID).
		SetName(params.Name).
		SetCreatedBy(userID).
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
	if err := ensureFolder(ctx, params.ProjectID, params.FolderID); err != nil {
		return err
	}

	_, err := database.DB.Folder.
		UpdateOneID(params.FolderID).
		SetName(params.Name).
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
	if err := ensureFolder(ctx, params.ProjectID, params.FolderID); err != nil {
		return err
	}

	folders, err := database.DB.Folder.
		Query().
		Where(entfolder.ProjectIDEQ(params.ProjectID)).
		All(ctx)
	if err != nil {
		return err
	}

	ifaces, err := database.DB.API.
		Query().
		Where(entapi.ProjectIDEQ(params.ProjectID)).
		All(ctx)
	if err != nil {
		return err
	}

	folderIDs := collectDescendants(params.FolderID, folders)
	folderIDSet := make(map[int64]struct{}, len(folderIDs))
	for _, id := range folderIDs {
		folderIDSet[id] = struct{}{}
	}

	for _, item := range ifaces {
		if _, ok := folderIDSet[item.FolderID]; ok {
			if err := database.DB.API.DeleteOneID(item.ID).Exec(ctx); err != nil {
				return err
			}
		}
	}

	for i := len(folderIDs) - 1; i >= 0; i-- {
		if err := database.DB.Folder.DeleteOneID(folderIDs[i]).Exec(ctx); err != nil {
			return err
		}
	}
	return nil
}

func requireProject(ctx context.Context, userID int64, params ProjectScopeRequest) error {
	if _, err := workspace.Require(ctx, userID, params.WorkspaceID, workspace.ActionProjectRead); err != nil {
		return err
	}
	return ensureProject(ctx, params.WorkspaceID, params.ProjectID)
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

func buildTree(folders []*ent.Folder, ifaces []*ent.API) []TreeNode {
	folderChildren := make(map[int64][]*ent.Folder)
	for _, folder := range folders {
		folderChildren[folder.ParentID] = append(folderChildren[folder.ParentID], folder)
	}

	ifaceChildren := make(map[int64][]*ent.API)
	for _, item := range ifaces {
		ifaceChildren[item.FolderID] = append(ifaceChildren[item.FolderID], item)
	}

	for parentID := range folderChildren {
		sort.Slice(folderChildren[parentID], func(i, j int) bool {
			left := folderChildren[parentID][i]
			right := folderChildren[parentID][j]
			if left.SortOrder != right.SortOrder {
				return left.SortOrder < right.SortOrder
			}
			return left.ID < right.ID
		})
	}
	for folderID := range ifaceChildren {
		sort.Slice(ifaceChildren[folderID], func(i, j int) bool {
			left := ifaceChildren[folderID][i]
			right := ifaceChildren[folderID][j]
			if left.SortOrder != right.SortOrder {
				return left.SortOrder < right.SortOrder
			}
			return left.ID < right.ID
		})
	}

	var buildFolders func(parentID int64) []TreeNode
	buildFolders = func(parentID int64) []TreeNode {
		children := folderChildren[parentID]
		nodes := make([]TreeNode, 0, len(children))
		for _, folder := range children {
			node := TreeNode{
				ID:        folderIDKey(folder.ID),
				ProjectID: folder.ProjectID,
				Name:      folder.Name,
				Type:      "folder",
			}
			subFolders := buildFolders(folder.ID)
			apiNodes := make([]TreeNode, 0, len(ifaceChildren[folder.ID]))
			for _, item := range ifaceChildren[folder.ID] {
				apiNodes = append(apiNodes, TreeNode{
					ID:        interfaceIDKey(item.ID),
					ProjectID: item.ProjectID,
					Name:      item.Name,
					Type:      "api",
					Method:    item.Method,
				})
			}
			node.Children = append(subFolders, apiNodes...)
			nodes = append(nodes, node)
		}
		return nodes
	}

	return buildFolders(0)
}

func collectDescendants(rootID int64, folders []*ent.Folder) []int64 {
	childrenMap := make(map[int64][]int64)
	for _, folder := range folders {
		childrenMap[folder.ParentID] = append(childrenMap[folder.ParentID], folder.ID)
	}

	result := []int64{rootID}
	queue := []int64{rootID}
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		for _, childID := range childrenMap[current] {
			result = append(result, childID)
			queue = append(queue, childID)
		}
	}
	return result
}

func folderIDKey(id int64) string {
	return fmt.Sprintf("folder-%d", id)
}

func interfaceIDKey(id int64) string {
	return fmt.Sprintf("api-%d", id)
}
