package actions

import (
	"github.com/davidalvarez305/soflo_go/server/database"
	"github.com/davidalvarez305/soflo_go/server/types"
)

func CreateUser(user types.User) (types.User, error) {
	result := database.DB.Create(user)

	if result.Error != nil {
		return user, result.Error
	}

	return user, nil
}
