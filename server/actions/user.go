package actions

import (
	"github.com/davidalvarez305/soflo_go/server/database"
	"github.com/davidalvarez305/soflo_go/server/models"
)

func CreateUser(user models.User) (models.User, error) {
	result := database.DB.Create(user)

	if result.Error != nil {
		return user, result.Error
	}

	return user, nil
}
