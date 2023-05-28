package http_tools

import (
	"api-gateway/core"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"net/http"
)

const (
	AuthHeader       = "Authorization"
	ClientIdCtxtKey  = "clientId"
	ClientRoleCtxKey = "clientRole"
)

func BindRequestBody(c *gin.Context, obj any) bool {
	if err := c.BindJSON(&obj); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, core.CodeIncorrectBody, core.ErrIncorrectBody.Error())
		return false
	}
	return true
}

func NewErrorResponse(c *gin.Context, statusCode int64, code int, message string) {
	logrus.Infof("%d: %s", statusCode, message)
	c.AbortWithStatusJSON(int(statusCode), core.ErrorBody{Message: message, Code: code})
}

func ErrorResponse(c *gin.Context, statusCode int64, message string) {
	logrus.Infof("%d: %s", statusCode, message)
	c.AbortWithStatusJSON(int(statusCode), core.ErrorBody{Message: message})
}

func GetClientId(c *gin.Context) (uuid.UUID, error) {
	idStr := c.GetString(ClientIdCtxtKey)
	return uuid.Parse(idStr)
}

// ParseUUIDFromParam returns Error response if it couldn't parse token
func ParseUUIDFromParam(c *gin.Context) uuid.UUID {
	id := c.Param("id")
	itemUUID, err := uuid.Parse(id)
	if err != nil {
		NewErrorResponse(c, http.StatusForbidden, core.CodeIncorrectBody, "could not parse uuid from id parameter")
		return uuid.Nil
	}
	return itemUUID
}
