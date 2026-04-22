package routers

import "github.com/gin-gonic/gin"

func (r *APIRouter) registerTaskRoutes(group *gin.RouterGroup) {
	tasks := group.Group("/tasks")
	{
		tasks.POST("", r.taskHandler.CreateTask)
		tasks.GET("", r.taskHandler.ListTasks)
	}
}
