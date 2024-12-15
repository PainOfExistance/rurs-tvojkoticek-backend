package Functions

import (
	"archive/zip"
	"backend/Mongo"
	"backend/Schemas"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"net/http"
	"time"
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

	// Get video_name from form data
	videoName := c.PostForm("video_name")
	if videoName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "video_name is required"})
		return
	}

	// Generate a unique filename by appending a timestamp
	uniqueFilename := fmt.Sprintf("%s_%d", videoName, time.Now().UnixNano())

	// Connect to GridFS bucket
	client := Mongo.GetMongoDB()
	bucket, err := gridfs.NewBucket(client.Database("Pametni-Paketnik-baza"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error creating GridFS bucket"})
		return
	}

	// Upload to GridFS
	videoID, err := bucket.UploadFromStream(uniqueFilename, fileStream)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error uploading to GridFS", "error": err.Error()})
		return
	}

	// Create video metadata using the Video schema
	videoMetadata := Schemas.Video{
		ID:          primitive.NewObjectID(),
		VideoName:   videoName, // Passed from the user
		Uploader:    c.PostForm("uploader_username"),
		Description: c.PostForm("description"),
		Tags:        c.PostFormArray("tags"),
		VideoID:     videoID.Hex(),
		Comments:    []Schemas.Comment{}, // Initialize empty comments
		PostedAt:    time.Now(),          // Timestamp
	}

	// Save metadata to the 'videostore' collection
	_, err = Mongo.GetCollection("videostore").InsertOne(context.TODO(), videoMetadata)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error saving video metadata", "error": err.Error()})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"message":   "Video uploaded successfully",
		"video_id":  videoID.Hex(),
		"posted_at": videoMetadata.PostedAt.Format(time.RFC3339),
	})
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
	// Connect to MongoDB and query all videos
	client := Mongo.GetMongoDB()
	collection := Mongo.GetCollection("videostore")

	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error retrieving video metadata"})
		return
	}
	defer cursor.Close(context.TODO())

	// Initialize GridFS bucket
	bucket, err := gridfs.NewBucket(client.Database("Pametni-Paketnik-baza"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error creating GridFS bucket"})
		return
	}

	// Set response headers for ZIP file
	c.Header("Content-Type", "application/zip")
	c.Header("Content-Disposition", `attachment; filename="all_videos.zip"`)

	// Create a ZIP writer
	zipWriter := zip.NewWriter(c.Writer)
	defer zipWriter.Close()

	// Loop through all videos and add them to the ZIP
	for cursor.Next(context.TODO()) {
		var video Schemas.Video
		if err := cursor.Decode(&video); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error decoding video metadata"})
			return
		}

		// Convert video_id to ObjectID
		objectID, err := primitive.ObjectIDFromHex(video.VideoID)
		if err != nil {
			fmt.Printf("Invalid ObjectID: %v\n", err)
			continue
		}

		// Create a file entry in the ZIP archive
		zipFile, err := zipWriter.Create(fmt.Sprintf("%s.mp4", video.VideoName))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error creating ZIP file entry"})
			return
		}

		// Stream the video data from GridFS into the ZIP file
		if _, err := bucket.DownloadToStream(objectID, zipFile); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error streaming video", "error": err.Error()})
			return
		}
	}

	if err := cursor.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Cursor error"})
		return
	}

	c.Status(http.StatusOK)
}

func GetAllVideosByName(c *gin.Context) {
	videoName := c.Query("video_name")
	if videoName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "video_name is required"})
		return
	}

	// Connect to MongoDB and query 'videostore' for all matching videos
	client := Mongo.GetMongoDB()
	collection := Mongo.GetCollection("videostore")

	cursor, err := collection.Find(context.TODO(), bson.M{"video_name": videoName})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error querying video metadata"})
		return
	}
	defer cursor.Close(context.TODO())

	// Initialize GridFS bucket
	bucket, err := gridfs.NewBucket(client.Database("Pametni-Paketnik-baza"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error creating GridFS bucket"})
		return
	}

	// Set response headers for a ZIP file
	c.Header("Content-Type", "application/zip")
	c.Header("Content-Disposition", `attachment; filename="videos.zip"`)

	// Create a ZIP writer
	zipWriter := zip.NewWriter(c.Writer)
	defer zipWriter.Close()

	// Loop through all matching video metadata and add to ZIP
	for cursor.Next(context.TODO()) {
		var video Schemas.Video
		if err := cursor.Decode(&video); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error decoding video metadata"})
			return
		}

		// Convert video_id to ObjectID
		objectID, err := primitive.ObjectIDFromHex(video.VideoID)
		if err != nil {
			fmt.Printf("Invalid ObjectID: %v\n", err)
			continue
		}

		// Create a file in the ZIP archive
		zipFile, err := zipWriter.Create(fmt.Sprintf("%s.mp4", video.VideoName))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error creating ZIP file"})
			return
		}

		// Stream the video data from GridFS into the ZIP file
		if _, err := bucket.DownloadToStream(objectID, zipFile); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error streaming video to ZIP", "error": err.Error()})
			return
		}
	}

	if err := cursor.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Cursor error"})
		return
	}

	// Finalize the ZIP file
	c.Status(http.StatusOK)
}
