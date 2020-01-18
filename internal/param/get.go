package param

import (
	"omega/internal/glog"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// Get is a function for filling param.Model
func Get(c *gin.Context) (param Param) {
	var err error
	var page uint64
	orderBy := "id"
	direction := "desc"

	if c.Query("order_by") != "" {
		orderBy = c.Query("order_by")
	}

	if c.Query("direction") != "" {
		direction = c.Query("direction")
	}

	param.Order = orderBy + " " + direction
	param.Select = "*"
	if c.Query("select") != "" {
		param.Select = c.Query("select")
	}

	if c.Query("page_size") != "" {
		param.Limit, err = strconv.ParseUint(c.Query("page_size"), 10, 16)
		if err != nil {
			// TODO: get path from gin.Context
			glog.CheckError(err, "Limit is not number")
			param.Limit = 10
		}
	}

	if c.Query("page") != "" {
		page, err = strconv.ParseUint(c.Query("page"), 10, 16)
		if err != nil {
			// TODO: get path from gin.Context
			glog.CheckError(err, "Offset is not a positive number")
			page = 1
		}
	}

	param.Search = strings.TrimSpace(c.Query("search"))

	userID, ok := c.Get("USER_ID")
	if ok {
		glog.CheckInfo(err, "User ID is not exist")
		param.UserID = userID.(uint64)
	}

	param.Offset = param.Limit * (page - 1)

	return param

}
