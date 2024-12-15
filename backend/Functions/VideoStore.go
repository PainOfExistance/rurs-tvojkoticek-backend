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

	// Connect to MongoDB and GridFS
	client := Mongo.GetMongoDB()
	bucket, err := gridfs.NewBucket(client.Database("Pametni-Paketnik-baza"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error creating GridFS bucket"})
		return
	}

	// Convert video_id to ObjectID
	objectID, err := primitive.ObjectIDFromHex(videoID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid video_id"})
		return
	}

	// Retrieve video metadata
	var videoMetadata Schemas.Video
	err = Mongo.GetCollection("videostore").FindOne(context.TODO(), bson.M{"video_id": videoID}).Decode(&videoMetadata)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Video metadata not found"})
		return
	}

	// Set response headers for a ZIP file
	c.Header("Content-Type", "application/zip")
	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.zip"`, videoMetadata.VideoName))

	// Create a ZIP writer
	zipWriter := zip.NewWriter(c.Writer)
	defer zipWriter.Close()

	// Add metadata.txt to the ZIP file
	metadataFile, err := zipWriter.Create("metadata.txt")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error creating metadata file"})
		return
	}

	// Write video metadata to metadata.txt
	metadataContent := fmt.Sprintf(
		"Video Name: %s\nUploader: %s\nDescription: %s\nTags: %v\nPosted At: %s\n",
		videoMetadata.VideoName,
		videoMetadata.Uploader,
		videoMetadata.Description,
		videoMetadata.Tags,
		videoMetadata.PostedAt.Format(time.RFC3339),
	)
	_, err = metadataFile.Write([]byte(metadataContent))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error writing metadata file"})
		return
	}

	// Add video file to the ZIP archive
	videoFile, err := zipWriter.Create(fmt.Sprintf("%s.mp4", videoMetadata.VideoName))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error creating video file entry in ZIP"})
		return
	}

	// Stream video data into the ZIP file
	_, err = bucket.DownloadToStream(objectID, videoFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error streaming video to ZIP", "error": err.Error()})
		return
	}

	// Finalize the ZIP file
	c.Status(http.StatusOK)
}

func GetAllVideos(c *gin.Context) {
	// Connect to MongoDB and GridFS
	client := Mongo.GetMongoDB()
	collection := Mongo.GetCollection("videostore")
	bucket, err := gridfs.NewBucket(client.Database("Pametni-Paketnik-baza"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error creating GridFS bucket"})
		return
	}

	// Query all videos
	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error querying video metadata"})
		return
	}
	defer cursor.Close(context.TODO())

	// Set response headers for ZIP file
	c.Header("Content-Type", "application/zip")
	c.Header("Content-Disposition", `attachment; filename="all_videos_with_metadata.zip"`)

	// Create a ZIP writer
	zipWriter := zip.NewWriter(c.Writer)
	defer zipWriter.Close()

	// Counter to track videos added
	var videoCount int

	// Loop through all videos
	for cursor.Next(context.TODO()) {
		var videoMetadata Schemas.Video
		if err := cursor.Decode(&videoMetadata); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error decoding video metadata"})
			return
		}

		// Convert video_id to ObjectID
		objectID, err := primitive.ObjectIDFromHex(videoMetadata.VideoID)
		if err != nil {
			fmt.Printf("Invalid ObjectID: %v\n", err)
			continue
		}

		// Add metadata.txt file for this video
		metadataFile, err := zipWriter.Create(fmt.Sprintf("%s_metadata.txt", videoMetadata.VideoName))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error creating metadata file"})
			return
		}

		// Write video metadata content
		metadataContent := fmt.Sprintf(
			"Video Name: %s\nUploader: %s\nDescription: %s\nTags: %v\nPosted At: %s\n",
			videoMetadata.VideoName,
			videoMetadata.Uploader,
			videoMetadata.Description,
			videoMetadata.Tags,
			videoMetadata.PostedAt.Format(time.RFC3339),
		)
		_, err = metadataFile.Write([]byte(metadataContent))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error writing metadata file"})
			return
		}

		// Add the video file itself
		videoFile, err := zipWriter.Create(fmt.Sprintf("%s.mp4", videoMetadata.VideoName))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error creating video file in ZIP"})
			return
		}

		// Stream video content into the ZIP file
		_, err = bucket.DownloadToStream(objectID, videoFile)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error streaming video to ZIP", "error": err.Error()})
			return
		}

		videoCount++
	}

	// Check for cursor errors
	if err := cursor.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Cursor error"})
		return
	}

	if videoCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No videos found in the database"})
		return
	}

	fmt.Printf("Total videos added to ZIP: %d\n", videoCount)
	c.Status(http.StatusOK)
}

func GetAllVideosByName(c *gin.Context) {
	videoName := c.Query("video_name")
	if videoName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "video_name is required"})
		return
	}

	// Connect to MongoDB and GridFS
	client := Mongo.GetMongoDB()
	collection := Mongo.GetCollection("videostore")
	bucket, err := gridfs.NewBucket(client.Database("Pametni-Paketnik-baza"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error creating GridFS bucket"})
		return
	}

	// Query for all videos matching the video_name
	cursor, err := collection.Find(context.TODO(), bson.M{"video_name": videoName})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error querying video metadata"})
		return
	}
	defer cursor.Close(context.TODO())

	// Set response headers for ZIP file
	c.Header("Content-Type", "application/zip")
	c.Header("Content-Disposition", `attachment; filename="videos_with_metadata.zip"`)

	// Create a ZIP writer
	zipWriter := zip.NewWriter(c.Writer)
	defer zipWriter.Close()

	// Counter to track files added
	var videoCount int

	// Loop through matching videos and add them to the ZIP
	for cursor.Next(context.TODO()) {
		var videoMetadata Schemas.Video
		if err := cursor.Decode(&videoMetadata); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error decoding video metadata"})
			return
		}

		// Convert video_id to ObjectID
		objectID, err := primitive.ObjectIDFromHex(videoMetadata.VideoID)
		if err != nil {
			fmt.Printf("Invalid ObjectID: %v\n", err)
			continue
		}

		// Add metadata.txt file for this video
		metadataFile, err := zipWriter.Create(fmt.Sprintf("%s_metadata.txt", videoMetadata.VideoName))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error creating metadata file"})
			return
		}

		metadataContent := fmt.Sprintf(
			"Video Name: %s\nUploader: %s\nDescription: %s\nTags: %v\nPosted At: %s\n",
			videoMetadata.VideoName,
			videoMetadata.Uploader,
			videoMetadata.Description,
			videoMetadata.Tags,
			videoMetadata.PostedAt.Format(time.RFC3339),
		)
		_, err = metadataFile.Write([]byte(metadataContent))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error writing metadata file"})
			return
		}

		// Add the video file itself
		videoFile, err := zipWriter.Create(fmt.Sprintf("%s.mp4", videoMetadata.VideoName))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error creating video file in ZIP"})
			return
		}

		// Stream video content into the ZIP file
		_, err = bucket.DownloadToStream(objectID, videoFile)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error streaming video", "error": err.Error()})
			return
		}

		videoCount++
	}

	if err := cursor.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Cursor error"})
		return
	}

	if videoCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No videos found with the given name"})
		return
	}

	fmt.Printf("Total videos added to ZIP: %d\n", videoCount)
	c.Status(http.StatusOK)
}
