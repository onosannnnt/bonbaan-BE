package roleUsecase

import (
	Entities "github.com/onosannnnt/bonbaan-BE/src/entities"
)

type RoleUsecase interface {
	InsertRole(role *Entities.Role) error
	GetAll() (*[]Entities.Role, error)
}

type RoleService struct {
	roleRepo RoleRepository
}

func NewRoleService(repo RoleRepository) RoleUsecase {
	return &RoleService{
		roleRepo: repo,
	}
}

func (s *RoleService) InsertRole(role *Entities.Role) error {
	return s.roleRepo.Insert(role)
}

func (s *RoleService) GetAll() (*[]Entities.Role, error) {
	return s.roleRepo.GetAll()
}
