package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/phamtrunghieu/tinder-clone-backend/controllers"
	"github.com/phamtrunghieu/tinder-clone-backend/middleware"
	"net/http"
)

func RouteInit(engine *gin.Engine) {

	userCtrl := new(controllers.UserController)
	userActionCtrl := new(controllers.UserActionController)

	engine.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Tinder server")
	})
	engine.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})
	engine.Use(middleware.CORS())

	apiUser := engine.Group("/user")
	apiUser.GET("", userCtrl.GetList)
	apiUser.GET("/:uuid", userCtrl.GetDetail)
	apiUser.GET("/random", userCtrl.GetUserRandom)
	apiUser.GET("/available-user/:uuid", userCtrl.GetUserAvailable)
	apiUser.GET("/matches-user/:uuid", userCtrl.GetMatchesUser)
	apiUser.GET("/generate-data", userCtrl.DumpData)

	apiUserAction := engine.Group("/user-action")
	apiUserAction.PUT("/like/:user_uuid/:guest_uuid", userActionCtrl.LikeUser)
	apiUserAction.PUT("/pass/:user_uuid/:guest_uuid", userActionCtrl.PassUser)
	apiUserAction.GET("/liked/:uuid", userActionCtrl.GetUserLiked)

}
