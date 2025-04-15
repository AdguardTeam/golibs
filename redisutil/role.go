package redisutil

import (
	"fmt"

	"github.com/AdguardTeam/golibs/errors"
	"github.com/AdguardTeam/golibs/validate"
)

// Role is a valid Redis role.
type Role string

// Valid Redis roles and their string representations.
const (
	RoleStringMaster   = "master"
	RoleStringSentinel = "sentinel"
	RoleStringSlave    = "slave"

	RoleMaster   Role = RoleStringMaster
	RoleSentinel Role = RoleStringSentinel
	RoleSlave    Role = RoleStringSlave
)

// NewRole converts s into a role.
func NewRole(s string) (r Role, err error) {
	r = Role(s)
	err = r.Validate()
	if err != nil {
		// Don't wrap the error, because it's informative enough as is.
		return "", err
	}

	return r, nil
}

// type check
var _ validate.Interface = RoleMaster

// Validate implements the [validate.Interface] interface for Role.
func (r Role) Validate() (err error) {
	switch r {
	case RoleMaster, RoleSentinel, RoleSlave:
		return nil
	default:
		return fmt.Errorf("%w: %q", errors.ErrBadEnumValue, string(r))
	}
}
