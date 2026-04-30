package logic

import "strings"

const (
    RoleClient         = "client"
    RoleEmploye        = "employe"
    RoleReceptionniste = "receptionniste"
    RoleManager        = "manager"
)

func (u User) HasPermission(requiredRoles ...string) bool {
    for _, role := range requiredRoles {
        if u.Role == strings.TrimSpace(role) {
            return true
        }
    }
    return false
}