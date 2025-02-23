package roleUsecase

import Entities "github.com/onosannnnt/bonbaan-BE/src/entities"

type RoleRepository interface {
	Insert(role *Entities.Role) error
	GetAll() (*[]Entities.Role, error)
	GetByName(name *string) (*Entities.Role, error)
	GetByID(id *string) (*Entities.Role, error)
	Update(role *Entities.Role) error
	Delete(id *string) error
}
