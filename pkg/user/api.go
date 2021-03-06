package user

import (
	"net/http"
	"omega/engine"
	"omega/internal/param"
	"omega/internal/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

// API for injecting user service
type API struct {
	Service Service
	Engine  engine.Engine
}

// ProvideAPI for user is used in wire
func ProvideAPI(p Service) API {
	return API{Service: p, Engine: p.Engine}
}

// FindAll users
func (p *API) FindAll(c *gin.Context) {
	if p.Engine.CheckAccess(c, "users:names") {
		response.NoPermission(c)
		return
	}
	users, err := p.Service.FindAll()

	if err != nil {
		response.RecordNotFound(c, err, "users")
		return
	}

	response.Success(c, users)
}

// List of users
func (p *API) List(c *gin.Context) {
	if p.Engine.CheckAccess(c, "users:read") {
		response.NoPermissionRecord(c, p.Engine, "user-list-forbidden")
		return
	}

	params := param.Get(c)

	p.Engine.Debug(params)
	data, err := p.Service.List(params)
	if err != nil {
		response.RecordNotFound(c, err, "users")
		return
	}

	p.Engine.Record(c, "user-list")
	response.Success(c, data)
}

// FindByID is used for fetch a user by his id
func (p *API) FindByID(c *gin.Context) {
	if p.Engine.CheckAccess(c, "users:read") {
		response.NoPermissionRecord(c, p.Engine, "user-view-forbidden")
		return
	}
	id, err := strconv.ParseUint(c.Param("id"), 10, 16)
	if err != nil {
		response.InvalidID(c, err)
		return
	}

	user, err := p.Service.FindByID(id)

	if err != nil {
		response.RecordNotFound(c, err, "user")
		return
	}

	p.Engine.Record(c, "user-view")
	response.Success(c, user)
}

// Create user
func (p *API) Create(c *gin.Context) {
	var user User

	err := c.BindJSON(&user)
	if err != nil {
		response.ErrorInBinding(c, err, "create user")
		return
	}

	if p.Engine.CheckAccess(c, "users:write") {
		response.NoPermissionRecord(c, p.Engine, "user-create-forbidden", nil, user)
		return
	}

	createdUser, err := p.Service.Save(user)
	if err != nil {
		response.ErrorOnSave(c, err, "user")
		return
	}

	p.Engine.Record(c, "user-create", nil, user)
	response.SuccessSave(c, createdUser, "user/create")
}

// Update user
func (p *API) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 16)
	if err != nil {
		response.InvalidID(c, err)
		return
	}

	var user User

	if err = c.BindJSON(&user); err != nil {
		response.ErrorInBinding(c, err, "update user")
		return
	}
	user.ID = id
	if p.Engine.CheckAccess(c, "users:write") {
		response.NoPermissionRecord(c, p.Engine, "user-update-forbidden", nil, user)
		return
	}

	userBefore, err := p.Service.FindByID(id)
	if err != nil {
		response.RecordNotFound(c, err, "update user")
		return
	}

	updatedUser, err := p.Service.Save(user)
	if err != nil {
		response.ErrorOnSave(c, err, "update user")
		return
	}

	userBefore.Password = ""
	updatedUser.Password = ""

	p.Engine.Record(c, "user-update", userBefore, updatedUser)
	response.SuccessSave(c, updatedUser, "user updated")
}

// Delete user
func (p *API) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 16)
	if err != nil {
		response.InvalidID(c, err)
		return
	}
	if p.Engine.CheckAccess(c, "users:write") {
		response.NoPermissionRecord(c, p.Engine, "user-delete-forbidden", nil, id)
		return
	}

	var user User

	user, err = p.Service.FindByID(id)
	if err != nil {
		response.RecordNotFound(c, err, "delete user")
		return
	}

	err = p.Service.Delete(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &response.Result{
			Message: "Something went wrong, cannot delete this user",
			Code:    1500,
		})
		return
	}

	user.Password = ""
	p.Engine.Record(c, "user-delete", user)
	response.SuccessSave(c, user, "user successfully deleted")
}
