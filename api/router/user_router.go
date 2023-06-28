package router

import (
	"strconv"
	"voltaserve/errorpkg"
	"voltaserve/service"

	"github.com/gofiber/fiber/v2"
)

type UserRouter struct {
	userSvc *service.UserService
}

func NewUserRouter() *UserRouter {
	return &UserRouter{
		userSvc: service.NewUserService(),
	}
}

func (r *UserRouter) AppendRoutes(g fiber.Router) {
	g.Get("/", r.List)
}

// List godoc
//
//	@Summary		List
//	@Description	List
//	@Tags			Users
//	@Id				users_list
//	@Produce		json
//	@Param			query			query		string	false	"Query"
//	@Param			organization_id	query		string	false	"Organization ID"
//	@Param			group			query		string	false	"Group ID"
//	@Param			page			query		string	false	"Page"
//	@Param			size			query		string	false	"Size"
//	@Param			sort_by			query		string	false	"Sort By"
//	@Param			sort_order		query		string	false	"Sort Order"
//	@Success		200				{object}	service.UserList
//	@Failure		404				{object}	errorpkg.ErrorResponse
//	@Failure		500				{object}	errorpkg.ErrorResponse
//	@Router			/users [get]
func (r *UserRouter) List(c *fiber.Ctx) error {
	var err error
	var page int64
	if c.Query("page") == "" {
		page = 1
	} else {
		page, err = strconv.ParseInt(c.Query("page"), 10, 32)
		if err != nil {
			page = 1
		}
	}
	var size int64
	if c.Query("size") == "" {
		size = WorkspaceDefaultPageSize
	} else {
		size, err = strconv.ParseInt(c.Query("size"), 10, 32)
		if err != nil {
			return err
		}
	}
	sortBy := c.Query("sort_by")
	if !IsValidSortBy(sortBy) {
		return errorpkg.NewInvalidQueryParamError("sort_by")
	}
	sortOrder := c.Query("sort_order")
	if !IsValidSortOrder(sortOrder) {
		return errorpkg.NewInvalidQueryParamError("sort_order")
	}
	userID := GetUserID(c)
	var nonGroupMembersOnly bool
	if c.Query("non_group_members_only") != "" {
		nonGroupMembersOnly, err = strconv.ParseBool(c.Query("non_group_members_only"))
		if err != nil {
			return err
		}
	}
	res, err := r.userSvc.List(service.UserListOptions{
		Query:               c.Query("query"),
		OrganizationID:      c.Query("organization_id"),
		GroupID:             c.Query("group_id"),
		NonGroupMembersOnly: nonGroupMembersOnly,
		SortBy:              sortBy,
		SortOrder:           sortOrder,
		Page:                uint(page),
		Size:                uint(size),
	}, userID)
	if err != nil {
		return err
	}
	return c.JSON(res)
}
