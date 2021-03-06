package user

import (
	"fmt"
	"omega/engine"
	"omega/internal/param"
	"omega/utils/password"
)

// Service for injecting auth repo
type Service struct {
	Repo   Repo
	Engine engine.Engine
}

// ProvideService for user is used in wire
func ProvideService(p Repo) Service {
	return Service{Repo: p, Engine: p.Engine}
}

// FindAll users
func (p *Service) FindAll() (users []User, err error) {
	users, err = p.Repo.FindAll()
	p.Engine.CheckError(err, "all users")
	return
}

// List of users
func (p *Service) List(params param.Param) (data map[string]interface{}, err error) {
	data = make(map[string]interface{})

	data["users"], err = p.Repo.List(params)
	p.Engine.CheckError(err, "users list")

	data["count"], err = p.Repo.Count(params)
	p.Engine.CheckError(err, "users count")

	return
}

// FindByID for user
func (p *Service) FindByID(id uint64) (user User, err error) {
	user, err = p.Repo.FindByID(id)
	p.Engine.CheckError(err, fmt.Sprintf("User with id %v", id))

	return
}

// Save user
func (p *Service) Save(user User) (createdUser User, err error) {
	user.Password, err = password.Hash(user.Password, p.Engine.Environments.Setting.PasswordSalt)
	p.Engine.CheckError(err, fmt.Sprintf("Hashing password failed for %+v", user))

	createdUser, err = p.Repo.Save(user)
	p.Engine.CheckInfo(err, fmt.Sprintf("Failed in saving user for %+v", user))

	createdUser.Password = ""

	return
}

// func (p *Service) SaveSimple(model interface{}) (createdUser interface{}, err error) {
// 	user := *(model.(*User))
// 	user.Password, err = password.Hash(user.Password, p.Engine.Environments.Setting.PasswordSalt)
// 	p.Engine.CheckError(err, fmt.Sprintf("Hashing password failed for %+v", user))

// 	createdUser, err = p.Repo.Save(user)
// 	p.Engine.CheckInfo(err, fmt.Sprintf("Failed in saving user for %+v", user))

// 	// createdUser.Password = ""

// 	return
// }

// Delete user
func (p *Service) Delete(user User) error {
	return p.Repo.Delete(user)
}
