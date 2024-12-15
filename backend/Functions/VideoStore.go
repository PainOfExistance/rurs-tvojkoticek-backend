package Functions

import (
	"backend/Mongo"
	"backend/Schemas"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"net/http"
)

func UploadVideo(c *gin.Context) {
	// Get file from the request
	file, err := c.FormFile("video")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Error getting video file"})
		return
	}

	// Open the file stream
	fileStream, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error opening file"})
		return
	}
	defer fileStream.Close()

	// Connect to GridFS bucket
	client := Mongo.GetMongoDB()
	bucket, err := gridfs.NewBucket(client.Database("Pametni-Paketnik-baza"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error creating GridFS bucket"})
		return
	}

	// Upload to GridFS
	videoID, err := bucket.UploadFromStream(file.Filename, fileStream)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error uploading to GridFS", "error": err.Error()})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{"message": "Video uploaded successfully", "video_id": videoID.Hex()})
}

func GetVideo(c *gin.Context) {
	videoID := c.Query("video_id")
	if videoID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "video_id is required"})
		return
	}

	client := Mongo.GetMongoDB()
	bucket, err := gridfs.NewBucket(client.Database("Pametni-Paketnik-baza"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error creating GridFS bucket"})
		return
	}

	objectID, err := primitive.ObjectIDFromHex(videoID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid video_id"})
		return
	}

	c.Header("Content-Type", "video/mp4")
	_, err = bucket.DownloadToStream(objectID, c.Writer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error streaming video"})
		return
	}
}

func GetAllVideos(c *gin.Context) {
	var videos []Schemas.Video

	cursor, err := Mongo.GetCollection("videostore").Find(c, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error retrieving videos"})
		return
	}
	defer cursor.Close(c)

	for cursor.Next(c) {
		var video Schemas.Video
		if err := cursor.Decode(&video); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error decoding video"})
			return
		}
		videos = append(videos, video)
	}

	if err := cursor.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Cursor error"})
		return
	}

	c.JSON(http.StatusOK, videos)
}
