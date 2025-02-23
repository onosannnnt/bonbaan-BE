package roleUsecase

import (
	Entities "github.com/onosannnnt/bonbaan-BE/src/entities"
)

type RoleUsecase interface {
	InsertRole(role *Entities.Role) error
	GetAll() (*[]Entities.Role, error)
	GetByName(name *string) (*Entities.Role, error)
	GetByID(id *string) (*Entities.Role, error)
	Update(role *Entities.Role) error
	Delete(id *string) error
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

func (s *RoleService) GetByName(name *string) (*Entities.Role, error) {
	return s.roleRepo.GetByName(name)
}

func (s *RoleService) GetByID(id *string) (*Entities.Role, error) {
	return s.roleRepo.GetByID(id)
}

func (s *RoleService) Update(role *Entities.Role) error {
	return s.roleRepo.Update(role)
}

func (s *RoleService) Delete(id *string) error {
	return s.roleRepo.Delete(id)
}
