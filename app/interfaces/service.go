package interfaces

import (
	"context"
	"strings"

	"nest-api/app/workspace"
	"nest-api/internal/database"
	"nest-api/internal/ent"
	entapi "nest-api/internal/ent/api"
	entbodyfield "nest-api/internal/ent/interfacebodyfield"
	entexample "nest-api/internal/ent/interfaceexample"
	entfield "nest-api/internal/ent/interfacefield"
	entheader "nest-api/internal/ent/interfaceheader"
	entquery "nest-api/internal/ent/interfacequeryparam"
	entreqheader "nest-api/internal/ent/interfacerequestheader"
	entresult "nest-api/internal/ent/interfaceresult"
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
		Order(ent.Asc(entapi.FieldSortOrder), ent.Asc(entapi.FieldID)).
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

func (Service) Detail(ctx context.Context, userID int64, params DetailRequest) (*DetailItem, error) {
	if _, err := workspace.Require(ctx, userID, params.WorkspaceID, workspace.ActionProjectRead); err != nil {
		return nil, err
	}
	if err := ensureProject(ctx, params.WorkspaceID, params.ProjectID); err != nil {
		return nil, err
	}
	if err := ensureInterface(ctx, params.ProjectID, params.InterfaceID); err != nil {
		return nil, err
	}

	row, err := database.DB.API.
		Query().
		Where(
			entapi.IDEQ(params.InterfaceID),
			entapi.ProjectIDEQ(params.ProjectID),
		).
		WithFolder().
		WithUpdater().
		WithRequestHeaders(func(q *ent.InterfaceRequestHeaderQuery) {
			q.Order(ent.Asc(entreqheader.FieldSortOrder), ent.Asc(entreqheader.FieldID))
		}).
		WithQueryParams(func(q *ent.InterfaceQueryParamQuery) {
			q.Order(ent.Asc(entquery.FieldSortOrder), ent.Asc(entquery.FieldID))
		}).
		WithBodyFields(func(q *ent.InterfaceBodyFieldQuery) {
			q.Order(ent.Asc(entbodyfield.FieldSortOrder), ent.Asc(entbodyfield.FieldID))
		}).
		WithResponseHeaders(func(q *ent.InterfaceHeaderQuery) {
			q.Order(ent.Asc(entheader.FieldSortOrder), ent.Asc(entheader.FieldID))
		}).
		WithResponseResults(func(q *ent.InterfaceResultQuery) {
			q.Order(ent.Asc(entresult.FieldSortOrder), ent.Asc(entresult.FieldID)).
				WithFields(func(fq *ent.InterfaceFieldQuery) {
					fq.Order(ent.Asc(entfield.FieldSortOrder), ent.Asc(entfield.FieldID))
				})
		}).
		WithResponseExamples(func(q *ent.InterfaceExampleQuery) {
			q.Order(ent.Asc(entexample.FieldSortOrder), ent.Asc(entexample.FieldID))
		}).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	folderName := ""
	if row.Edges.Folder != nil {
		folderName = row.Edges.Folder.Name
	}

	detail := &DetailItem{
		Item: Item{
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
		},
		RequestHeaders: buildParamItems(row.Edges.RequestHeaders),
		RequestBody: RequestBodyConfig{
			Format:   row.RequestBodyFormat,
			DataType: row.RequestBodyDataType,
			Fields:   buildBodyFieldTree(row.Edges.BodyFields),
		},
		QueryParams:      buildQueryParamItems(row.Edges.QueryParams),
		ResponseHeaders:  buildResponseHeaders(row.Edges.ResponseHeaders),
		ResponseResults:  buildResponseResults(row.Edges.ResponseResults),
		ResponseExamples: buildResponseExamples(row.Edges.ResponseExamples),
	}
	return detail, nil
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

	sortOrder, err := nextInterfaceSortOrder(ctx, params.ProjectID, params.FolderID)
	if err != nil {
		return 0, err
	}

	row, err := database.DB.API.
		Create().
		SetProjectID(params.ProjectID).
		SetFolderID(params.FolderID).
		SetName(params.Name).
		SetMethod(method).
		SetURL(params.URL).
		SetStatus(status).
		SetSortOrder(sortOrder).
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

	bodyFormat := params.RequestBody.Format
	if bodyFormat == "" {
		bodyFormat = "json"
	}
	bodyDataType := params.RequestBody.DataType
	if bodyDataType == "" {
		bodyDataType = "Object"
	}

	tx, err := database.DB.Tx(ctx)
	if err != nil {
		return err
	}

	updater := tx.API.
		UpdateOneID(params.InterfaceID).
		SetName(params.Name).
		SetMethod(strings.ToUpper(params.Method)).
		SetURL(params.URL).
		SetStatus(status).
		SetRequestBodyFormat(bodyFormat).
		SetRequestBodyDataType(bodyDataType).
		SetUpdatedBy(userID)

	if params.FolderID > 0 {
		if err := ensureFolder(ctx, params.ProjectID, params.FolderID); err != nil {
			_ = tx.Rollback()
			return err
		}
		updater = updater.SetFolderID(params.FolderID)
	}

	if err := updater.Exec(ctx); err != nil {
		_ = tx.Rollback()
		return err
	}

	if err := replaceRequestHeaders(ctx, tx, params.InterfaceID, params.RequestHeaders); err != nil {
		_ = tx.Rollback()
		return err
	}
	if err := replaceQueryParams(ctx, tx, params.InterfaceID, params.QueryParams); err != nil {
		_ = tx.Rollback()
		return err
	}
	if err := replaceBodyFields(ctx, tx, params.InterfaceID, params.RequestBody.Fields); err != nil {
		_ = tx.Rollback()
		return err
	}
	if err := replaceResponseHeaders(ctx, tx, params.InterfaceID, params.ResponseHeaders); err != nil {
		_ = tx.Rollback()
		return err
	}
	if err := replaceResponseResults(ctx, tx, params.InterfaceID, params.ResponseResults); err != nil {
		_ = tx.Rollback()
		return err
	}
	if err := replaceResponseExamples(ctx, tx, params.InterfaceID, params.ResponseExamples); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
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

func (Service) Reorder(ctx context.Context, userID int64, params ReorderRequest) error {
	if _, err := workspace.Require(ctx, userID, params.WorkspaceID, workspace.ActionProjectUpdate); err != nil {
		return err
	}
	if err := ensureProject(ctx, params.WorkspaceID, params.ProjectID); err != nil {
		return err
	}
	if err := ensureFolder(ctx, params.ProjectID, params.FolderID); err != nil {
		return err
	}

	rows, err := database.DB.API.
		Query().
		Where(
			entapi.ProjectIDEQ(params.ProjectID),
			entapi.FolderIDEQ(params.FolderID),
		).
		All(ctx)
	if err != nil {
		return err
	}

	if len(params.InterfaceIDs) != len(rows) {
		return bizerr.New("接口列表不完整")
	}

	existingIDs := make(map[int64]struct{}, len(rows))
	for _, row := range rows {
		existingIDs[row.ID] = struct{}{}
	}

	seen := make(map[int64]struct{}, len(params.InterfaceIDs))
	for _, id := range params.InterfaceIDs {
		if _, ok := existingIDs[id]; !ok {
			return bizerr.New("接口不存在或不属于该分组")
		}
		if _, dup := seen[id]; dup {
			return bizerr.New("接口列表重复")
		}
		seen[id] = struct{}{}
	}

	tx, err := database.DB.Tx(ctx)
	if err != nil {
		return err
	}

	for index, id := range params.InterfaceIDs {
		if err := tx.API.UpdateOneID(id).SetSortOrder(index).Exec(ctx); err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func buildParamItems(rows []*ent.InterfaceRequestHeader) []ParamItem {
	items := make([]ParamItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, ParamItem{
			ID:          row.ID,
			Name:        row.Name,
			Type:        row.Type,
			Required:    row.Required,
			Description: row.Description,
			Example:     row.Example,
		})
	}
	return items
}

func buildQueryParamItems(rows []*ent.InterfaceQueryParam) []ParamItem {
	items := make([]ParamItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, ParamItem{
			ID:          row.ID,
			Name:        row.Name,
			Type:        row.Type,
			Required:    row.Required,
			Description: row.Description,
			Example:     row.Example,
		})
	}
	return items
}

func buildBodyFieldTree(rows []*ent.InterfaceBodyField) []BodyFieldItem {
	if len(rows) == 0 {
		return []BodyFieldItem{}
	}

	byParent := make(map[int64][]*ent.InterfaceBodyField)
	for _, row := range rows {
		byParent[row.ParentID] = append(byParent[row.ParentID], row)
	}

	var walk func(parentID int64) []BodyFieldItem
	walk = func(parentID int64) []BodyFieldItem {
		children := byParent[parentID]
		items := make([]BodyFieldItem, 0, len(children))
		for _, row := range children {
			items = append(items, BodyFieldItem{
				ID:          row.ID,
				ParentID:    row.ParentID,
				Name:        row.Name,
				Type:        row.Type,
				Required:    row.Required,
				Description: row.Description,
				Example:     row.Example,
				Children:    walk(row.ID),
			})
		}
		return items
	}

	return walk(0)
}

func buildResponseHeaders(rows []*ent.InterfaceHeader) []ResponseHeaderItem {
	items := make([]ResponseHeaderItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, ResponseHeaderItem{
			ID:          row.ID,
			Name:        row.Name,
			Type:        row.Type,
			Required:    row.Required,
			Description: row.Description,
			Example:     row.Example,
		})
	}
	return items
}

func buildResponseResults(rows []*ent.InterfaceResult) []ResponseResultItem {
	items := make([]ResponseResultItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, ResponseResultItem{
			ID:         row.ID,
			Name:       row.Name,
			StatusCode: row.StatusCode,
			Format:     row.Format,
			DataType:   row.DataType,
			Fields:     buildResponseFieldTree(row.Edges.Fields),
		})
	}
	return items
}

func buildResponseFieldTree(rows []*ent.InterfaceField) []ResponseFieldItem {
	if len(rows) == 0 {
		return []ResponseFieldItem{}
	}

	byParent := make(map[int64][]*ent.InterfaceField)
	for _, row := range rows {
		byParent[row.ParentID] = append(byParent[row.ParentID], row)
	}

	var walk func(parentID int64) []ResponseFieldItem
	walk = func(parentID int64) []ResponseFieldItem {
		children := byParent[parentID]
		items := make([]ResponseFieldItem, 0, len(children))
		for _, row := range children {
			items = append(items, ResponseFieldItem{
				ID:          row.ID,
				ParentID:    row.ParentID,
				Name:        row.Name,
				Type:        row.Type,
				Required:    row.Required,
				Description: row.Description,
				Mock:        row.Mock,
				Example:     row.Example,
				Children:    walk(row.ID),
			})
		}
		return items
	}

	return walk(0)
}

func buildResponseExamples(rows []*ent.InterfaceExample) []ResponseExampleItem {
	items := make([]ResponseExampleItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, ResponseExampleItem{
			ID:          row.ID,
			Name:        row.Name,
			StatusCode:  row.StatusCode,
			ContentType: row.ContentType,
			Raw:         row.Raw,
		})
	}
	return items
}

func replaceRequestHeaders(ctx context.Context, tx *ent.Tx, interfaceID int64, headers []ParamItem) error {
	if _, err := tx.InterfaceRequestHeader.Delete().Where(entreqheader.InterfaceIDEQ(interfaceID)).Exec(ctx); err != nil {
		return err
	}
	for index, header := range headers {
		if strings.TrimSpace(header.Name) == "" {
			continue
		}
		headerType := header.Type
		if headerType == "" {
			headerType = "string"
		}
		if _, err := tx.InterfaceRequestHeader.
			Create().
			SetInterfaceID(interfaceID).
			SetName(header.Name).
			SetType(headerType).
			SetRequired(header.Required).
			SetDescription(header.Description).
			SetExample(header.Example).
			SetSortOrder(index).
			Save(ctx); err != nil {
			return err
		}
	}
	return nil
}

func replaceQueryParams(ctx context.Context, tx *ent.Tx, interfaceID int64, params []ParamItem) error {
	if _, err := tx.InterfaceQueryParam.Delete().Where(entquery.InterfaceIDEQ(interfaceID)).Exec(ctx); err != nil {
		return err
	}
	for index, param := range params {
		if strings.TrimSpace(param.Name) == "" {
			continue
		}
		paramType := param.Type
		if paramType == "" {
			paramType = "string"
		}
		if _, err := tx.InterfaceQueryParam.
			Create().
			SetInterfaceID(interfaceID).
			SetName(param.Name).
			SetType(paramType).
			SetRequired(param.Required).
			SetDescription(param.Description).
			SetExample(param.Example).
			SetSortOrder(index).
			Save(ctx); err != nil {
			return err
		}
	}
	return nil
}

func replaceBodyFields(ctx context.Context, tx *ent.Tx, interfaceID int64, fields []BodyFieldItem) error {
	if _, err := tx.InterfaceBodyField.Delete().Where(entbodyfield.InterfaceIDEQ(interfaceID)).Exec(ctx); err != nil {
		return err
	}
	return saveBodyFields(ctx, tx, interfaceID, 0, fields)
}

func saveBodyFields(ctx context.Context, tx *ent.Tx, interfaceID, parentID int64, fields []BodyFieldItem) error {
	for index, fieldItem := range fields {
		if strings.TrimSpace(fieldItem.Name) == "" {
			continue
		}
		fieldType := fieldItem.Type
		if fieldType == "" {
			fieldType = "string"
		}

		saved, err := tx.InterfaceBodyField.
			Create().
			SetInterfaceID(interfaceID).
			SetParentID(parentID).
			SetName(fieldItem.Name).
			SetType(fieldType).
			SetRequired(fieldItem.Required).
			SetDescription(fieldItem.Description).
			SetExample(fieldItem.Example).
			SetSortOrder(index).
			Save(ctx)
		if err != nil {
			return err
		}

		if len(fieldItem.Children) > 0 {
			if err := saveBodyFields(ctx, tx, interfaceID, saved.ID, fieldItem.Children); err != nil {
				return err
			}
		}
	}
	return nil
}

func replaceResponseHeaders(ctx context.Context, tx *ent.Tx, interfaceID int64, headers []ResponseHeaderItem) error {
	if _, err := tx.InterfaceHeader.Delete().Where(entheader.InterfaceIDEQ(interfaceID)).Exec(ctx); err != nil {
		return err
	}

	for index, header := range headers {
		if strings.TrimSpace(header.Name) == "" {
			continue
		}
		headerType := header.Type
		if headerType == "" {
			headerType = "string"
		}
		if _, err := tx.InterfaceHeader.
			Create().
			SetInterfaceID(interfaceID).
			SetName(header.Name).
			SetType(headerType).
			SetRequired(header.Required).
			SetDescription(header.Description).
			SetExample(header.Example).
			SetSortOrder(index).
			Save(ctx); err != nil {
			return err
		}
	}
	return nil
}

func replaceResponseResults(ctx context.Context, tx *ent.Tx, interfaceID int64, results []ResponseResultItem) error {
	existingResults, err := tx.InterfaceResult.
		Query().
		Where(entresult.InterfaceIDEQ(interfaceID)).
		All(ctx)
	if err != nil {
		return err
	}
	for _, row := range existingResults {
		if _, err := tx.InterfaceField.Delete().Where(entfield.ResultIDEQ(row.ID)).Exec(ctx); err != nil {
			return err
		}
	}
	if _, err := tx.InterfaceResult.Delete().Where(entresult.InterfaceIDEQ(interfaceID)).Exec(ctx); err != nil {
		return err
	}

	for index, result := range results {
		format := result.Format
		if format == "" {
			format = "JSON"
		}
		dataType := result.DataType
		if dataType == "" {
			dataType = "Object"
		}
		statusCode := result.StatusCode
		if statusCode == 0 {
			statusCode = 200
		}

		saved, err := tx.InterfaceResult.
			Create().
			SetInterfaceID(interfaceID).
			SetName(result.Name).
			SetStatusCode(statusCode).
			SetFormat(format).
			SetDataType(dataType).
			SetSortOrder(index).
			Save(ctx)
		if err != nil {
			return err
		}

		if err := saveResponseFields(ctx, tx, saved.ID, 0, result.Fields); err != nil {
			return err
		}
	}
	return nil
}

func saveResponseFields(ctx context.Context, tx *ent.Tx, resultID, parentID int64, fields []ResponseFieldItem) error {
	for index, fieldItem := range fields {
		if strings.TrimSpace(fieldItem.Name) == "" {
			continue
		}
		fieldType := fieldItem.Type
		if fieldType == "" {
			fieldType = "string"
		}

		saved, err := tx.InterfaceField.
			Create().
			SetResultID(resultID).
			SetParentID(parentID).
			SetName(fieldItem.Name).
			SetType(fieldType).
			SetRequired(fieldItem.Required).
			SetDescription(fieldItem.Description).
			SetMock(fieldItem.Mock).
			SetExample(fieldItem.Example).
			SetSortOrder(index).
			Save(ctx)
		if err != nil {
			return err
		}

		if len(fieldItem.Children) > 0 {
			if err := saveResponseFields(ctx, tx, resultID, saved.ID, fieldItem.Children); err != nil {
				return err
			}
		}
	}
	return nil
}

func replaceResponseExamples(ctx context.Context, tx *ent.Tx, interfaceID int64, examples []ResponseExampleItem) error {
	if _, err := tx.InterfaceExample.Delete().Where(entexample.InterfaceIDEQ(interfaceID)).Exec(ctx); err != nil {
		return err
	}

	for index, example := range examples {
		statusCode := example.StatusCode
		if statusCode == 0 {
			statusCode = 200
		}
		contentType := example.ContentType
		if contentType == "" {
			contentType = "application/json"
		}
		name := example.Name
		if name == "" {
			name = "示例"
		}

		if _, err := tx.InterfaceExample.
			Create().
			SetInterfaceID(interfaceID).
			SetName(name).
			SetStatusCode(statusCode).
			SetContentType(contentType).
			SetRaw(example.Raw).
			SetSortOrder(index).
			Save(ctx); err != nil {
			return err
		}
	}
	return nil
}

func nextInterfaceSortOrder(ctx context.Context, projectID, folderID int64) (int, error) {
	rows, err := database.DB.API.
		Query().
		Where(
			entapi.ProjectIDEQ(projectID),
			entapi.FolderIDEQ(folderID),
		).
		Order(ent.Desc(entapi.FieldSortOrder)).
		Limit(1).
		All(ctx)
	if err != nil {
		return 0, err
	}
	if len(rows) == 0 {
		return 0, nil
	}
	return rows[0].SortOrder + 1, nil
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
