package roleUsecase

import (
	"github.com/onosannnnt/bonbaan-BE/src/entities"
)

type RoleRepository interface {
	Insert(role *Entities.Role) error
	GetAll() (*[]Entities.Role, error)
}
