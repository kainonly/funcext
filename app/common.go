package app

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/weplanx/go/route"
	"github.com/weplanx/openapi/app/geo"
	"github.com/weplanx/openapi/app/index"
	"github.com/weplanx/openapi/common"
)

var Provides = wire.NewSet(
	index.Provides,
	geo.Provides,
	New,
)

func New(
	values *common.Values,
	index *index.Controller,
	geo *geo.Controller,
) *gin.Engine {
	r := globalMiddleware(gin.New(), values)
	r.GET("/", route.Use(index.Index))
	r.GET("/ip", route.Use(index.Ip))
	_geo := r.Group("/geo")
	{
		_geo.GET("/countries", route.Use(geo.Countries))
		_geo.GET("/states", route.Use(geo.States))
		_geo.GET("/cities", route.Use(geo.Cities))
	}
	return r
}