package story

import (
	// "time"
	// "context"
	// "database/sql"

	// "gitlab.com/pedrokoblitz/opus-go/internal/adapters"
	// "gitlab.com/pedrokoblitz/opus-go/internal/payloads"
	// "gitlab.com/pedrokoblitz/opus-go/internal/quality"
)

// type PermissionChecker struct {
// 	adapter *sql.DB
// }

// func NewPermissionChecker(adapter *sql.DB) *PermissionChecker {
// 	return &PermissionChecker{
// 		adapter: adapter,
// 	}
// }

// func (pm *PermissionChecker) IsSuperAdmin(ctx context.Context, userID uint) error {
// 	pr := &PermissionRequirement{
// 		UserID: userID,
// 		RoleSlug: "super-admin",
// 	}
// 	_, err := pm.adapter.ExecContext(
// 		ctx,
// 		CheckSuperAdminMySQLQuery,
// 		pr.RoleSlug,
// 	)
// 	if err != nil {
// 		return quality.ErrExecution.WithDetail(err.Error())
// 	}
// 	return nil

// }

// func (pm *PermissionChecker) IsAdmin(ctx context.Context, story string, userID uint) error {
// 	pr := &PermissionRequirement{
// 		UserID: userID,
// 		StorySlug: story,
// 		RoleSlug: "admin",
// 	}
// 	_, err := pm.adapter.ExecContext(
// 		ctx,
// 		CheckAdminMySQLQuery,
// 		pr.StorySlug, pr.RoleSlug,
// 	)
// 	if err != nil {
// 		return quality.ErrExecution.WithDetail(err.Error())
// 	}
// 	return nil
// }

// func (pm *PermissionChecker) Check(ctx context.Context, pr PermissionRequirement) error {
// 	_, err := pm.adapter.ExecContext(
// 		ctx,
// 		CheckPermissionMySQLQuery,
// 		pr.StoryID, pr.RoleID, pr.DefinitionID, pr.Action,
// 	)
// 	if err != nil {
// 		return quality.ErrExecution.WithDetail(err.Error())
// 	}
// 	return nil
// }
