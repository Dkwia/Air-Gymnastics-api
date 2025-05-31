package routes

import (
	"github.com/gin-gonic/gin"
	"go-whatsapp-api/handlers"
)

func SetupRoutes(r *gin.Engine) {
	api := r.Group("/api/v1")
	{
		users := api.Group("/user")
		{
			users.POST("/registration/", handlers.CreateUser)
			users.GET("/login/", handlers.LoginUser)
			users.GET("/data/", handlers.GetUser)
			users.POST("/data/", handlers.UpdateUser)
		}

		main := api.Group("/main")
		{
			main.POST("/news/", handlers.UpdateNews)
			main.GET("/news/", handlers.GetNews)
			main.POST("/competitions/", handlers.UpdateCompetition)
			main.GET("/competitions/", handlers.GetCompetition)
		}

		schedule := api.Group("/schedule")
		{
			schedule.POST("/", handlers.UpdateSchedule)
			schedule.GET("/", handlers.GetSchedule)
		}
	}
}
