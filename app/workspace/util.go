package workspace

import "nest-api/internal/ent"

func UserDisplayName(u *ent.User) string {
	if u == nil {
		return ""
	}
	if u.Name != "" {
		return u.Name
	}
	return u.Account
}
