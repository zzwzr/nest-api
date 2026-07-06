package workspace

const (
	RoleOwner  uint8 = 1
	RoleAdmin  uint8 = 2
	RoleEditor uint8 = 3
	RoleViewer uint8 = 4
)

func IsValidMemberRole(role uint8) bool {
	return role == RoleAdmin || role == RoleEditor || role == RoleViewer
}
