package HTTP

import (
	"backend/Functions"

	"github.com/gin-gonic/gin"
)

func Router(router *gin.Engine) {
	router.POST("/register", Functions.Register)
	router.POST("/login", Functions.Login)
	router.GET("/profile", Functions.GetProfile)
	router.POST("/changePassword", Functions.ChangePassword)

	router.GET("/post", Functions.GetPost)
	router.GET("/posts", Functions.GetAllPosts)
	router.POST("/post", Functions.CreatePost)
	router.DELETE("/post", Functions.DeletePost)

	router.POST("/comment", Functions.CreateComment)
	router.DELETE("/comment", Functions.DeleteComment)

	router.POST("/videostore/upload", Functions.UploadVideo)
	router.GET("/videostore/video:id", Functions.GetVideo)
	router.GET("/videostore/all", Functions.GetAllVideos)
	router.GET("/videostore/videos/name", Functions.GetAllVideosByName)
	router.DELETE("/videostore/video:video_id", Functions.DeleteVideoByID)
	router.POST("/videostore/flag", Functions.FlagVideo)
	router.GET("/videostore/flagged", Functions.GetFlaggedVideos)
	router.POST("/videostore/reset-flagged", Functions.ResetFlaggedCounter)

}
