package api

import (
	"context"
	"github.com/bytedance/go-tagexpr/v2/binding"
	"github.com/bytedance/go-tagexpr/v2/validator"
	"github.com/bytedance/gopkg/util/logger"
	"github.com/bytedance/sonic/decoder"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/errors"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/google/wire"
	"github.com/weplanx/openapi/api/geo"
	"github.com/weplanx/openapi/api/index"
	"github.com/weplanx/openapi/common"
	"net/http"
)

var Provides = wire.NewSet(
	index.Provides,
	geo.Provides,
)

type API struct {
	*common.Inject

	Hertz *server.Hertz

	IndexController *index.Controller
	IndexService    *index.Service
	GeoController   *geo.Controller
	GeoService      *geo.Service
}

func (x *API) Routes() (h *server.Hertz, err error) {
	h = x.Hertz
	h.Use(x.ErrHandler())

	h.GET("/", x.IndexController.Index)
	h.GET("/ip", x.IndexController.GetIp)

	_geo := h.Group("/geo")
	{
		_geo.GET("/countries", x.GeoController.Countries)
		_geo.GET("/states", x.GeoController.States)
		_geo.GET("/cities", x.GeoController.Cities)
	}

	return
}

// ErrHandler 错误处理
func (x *API) ErrHandler() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		c.Next(ctx)
		err := c.Errors.Last()
		if err == nil {
			return
		}

		if err.IsType(errors.ErrorTypePublic) {
			statusCode := http.StatusBadRequest
			result := utils.H{"message": err.Error()}
			if meta, ok := err.Meta.(map[string]interface{}); ok {
				if meta["statusCode"] != nil {
					statusCode = meta["statusCode"].(int)
				}
				if meta["code"] != nil {
					result["code"] = meta["code"]
				}
			}
			c.JSON(statusCode, result)
			return
		}

		switch any := err.Err.(type) {
		case decoder.SyntaxError:
			c.JSON(http.StatusBadRequest, utils.H{
				"message": any.Description(),
			})
			break
		case *binding.Error:
			c.JSON(http.StatusBadRequest, utils.H{
				"message": any.Error(),
			})
			break
		case *validator.Error:
			c.JSON(http.StatusBadRequest, utils.H{
				"message": any.Error(),
			})
			break
		default:
			logger.Error(err)
			c.Status(http.StatusInternalServerError)
		}
	}
}