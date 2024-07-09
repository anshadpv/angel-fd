package context

import (
	"context"

	"github.com/angel-one/go-utils/constants"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type key string
type KeyHeader string
type KeyId string
type KeyPath string

var userKey key

const (
	headers = "Headers"
	idKey   = "id"
	pathKey = "path"

	HeaderRequestID = "X-Request-ID"
	UserID          = "userID"
	CorrelationId   = "correlationId"
)

type Info struct {
	XRequestID    string `header:"X-Request-Id"`
	Authorization string `header:"Authorization"`
	IPAddress     string `header:"ip-address"`
	AppVersion    string `header:"app-version"`
	OsType        string `header:"os-type"`
	UserID        string `header:"userId"`
	UserType      string `header:"userType"`
	Platform      string `header:"X-Platform"`
	AppSource     string `header:"X-Source"`
	CorrelationID string ``
}

func Background(id string) context.Context {
	var info = Info{
		CorrelationID: uuid.NewString(),
	}
	sgbctx := context.WithValue(context.Background(), KeyHeader(headers), info)
	return context.WithValue(sgbctx, KeyId(idKey), id)
}

func Get(ctx context.Context) Info {
	if ctx.Value(KeyHeader(headers)) != nil {
		return ctx.Value(KeyHeader(headers)).(Info)
	}
	return Info{}
}

func Build(ctx *gin.Context) context.Context {
	var info Info
	_ = ctx.BindHeader(&info)
	info.CorrelationID = uuid.NewString()
	id := ctx.Value(constants.IDLogParam)
	path := ctx.Value(constants.PathLogParam)
	nbuctx := context.WithValue(ctx.Request.Context(), KeyHeader(headers), info)
	nbuctx = context.WithValue(nbuctx, KeyId(idKey), id)
	return context.WithValue(nbuctx, KeyPath(pathKey), path)
}
