package roles

const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)

var Roles = []string{
	RoleAdmin,
	RoleUser,
}

func IsRole(role string) bool {
	for _, r := range Roles {
		if r == role {
			return true
		}
	}

	return false
}
