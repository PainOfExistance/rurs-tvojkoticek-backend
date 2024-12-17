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
	"go.mongodb.org/mongo-driver/mongo"
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

	// Query the database for existing videos with the same name
	collection := Mongo.GetCollection("videostore")
	cursor, err := collection.Find(context.TODO(), bson.M{"video_name": videoName})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error querying video name count"})
		return
	}
	defer cursor.Close(context.TODO())

	// Count existing videos with the same name
	count := 0
	for cursor.Next(context.TODO()) {
		count++
	}

	// Generate a unique filename by appending index and timestamp
	uniqueFilename := fmt.Sprintf("%s_%d_%d", videoName, count+1, time.Now().UnixNano())
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
		PostedAt:    time.Now(),
		Flagged:     0,
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

	// Retrieve video metadata excluding flagged videos
	var videoMetadata Schemas.Video
	err = Mongo.GetCollection("videostore").FindOne(context.TODO(), bson.M{"video_id": videoID, "flagged": bson.M{"$lte": 3}}).Decode(&videoMetadata)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"message": "Video not found or flagged"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error retrieving video metadata"})
		}
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
		"Video Name: %s\nUploader: %s\nDescription: %s\nTags: %v\nPosted At: %s\nVideo ID: %s\nFlagged Count: %d\\n",
		videoMetadata.VideoName,
		videoMetadata.Uploader,
		videoMetadata.Description,
		videoMetadata.Tags,
		videoMetadata.PostedAt.Format(time.RFC3339),
		videoMetadata.VideoID,
		videoMetadata.Flagged,
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
	client := Mongo.GetMongoDB()
	collection := Mongo.GetCollection("videostore")
	bucket, err := gridfs.NewBucket(client.Database("Pametni-Paketnik-baza"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error creating GridFS bucket"})
		return
	}

	// Query all videos excluding flagged ones
	cursor, err := collection.Find(context.TODO(), bson.M{"flagged": bson.M{"$lte": 3}})
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

	index := 1
	videoCount := 0

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

		// Add metadata.txt file for this video with an index
		metadataFile, err := zipWriter.Create(fmt.Sprintf("%s_%d_metadata.txt", videoMetadata.VideoName, index))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error creating metadata file"})
			return
		}

		metadataContent := fmt.Sprintf(
			"Video Name: %s\nUploader: %s\nDescription: %s\nTags: %v\nPosted At: %s\nVideo ID: %s\nFlagged Count: %d\\n",
			videoMetadata.VideoName,
			videoMetadata.Uploader,
			videoMetadata.Description,
			videoMetadata.Tags,
			videoMetadata.PostedAt.Format(time.RFC3339),
			videoMetadata.VideoID, // Include the video_id here
			videoMetadata.Flagged,
		)

		_, err = metadataFile.Write([]byte(metadataContent))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error writing metadata file"})
			return
		}

		// Add the video file itself with an index
		videoFile, err := zipWriter.Create(fmt.Sprintf("%s_%d.mp4", videoMetadata.VideoName, index))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error creating video file in ZIP"})
			return
		}

		_, err = bucket.DownloadToStream(objectID, videoFile)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error streaming video to ZIP", "error": err.Error()})
			return
		}

		index++
		videoCount++
	}

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

	// Query for all videos matching the video_name and flagged count
	cursor, err := collection.Find(context.TODO(), bson.M{"video_name": videoName, "flagged": bson.M{"$lte": 3}})
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

	// Initialize an index for file differentiation
	index := 1
	videoCount := 0

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

		// Add metadata.txt file for this video with an index
		metadataFile, err := zipWriter.Create(fmt.Sprintf("%s_%d_metadata.txt", videoMetadata.VideoName, index))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error creating metadata file"})
			return
		}

		// Write metadata content
		metadataContent := fmt.Sprintf(
			"Video Name: %s\nUploader: %s\nDescription: %s\nTags: %v\nPosted At: %s\nVideo ID: %s\nFlagged Count: %d\n",
			videoMetadata.VideoName,
			videoMetadata.Uploader,
			videoMetadata.Description,
			videoMetadata.Tags,
			videoMetadata.PostedAt.Format(time.RFC3339),
			videoMetadata.VideoID,
			videoMetadata.Flagged, // Add flagged count
		)

		_, err = metadataFile.Write([]byte(metadataContent))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error writing metadata file"})
			return
		}

		// Add the video file itself with an index
		videoFile, err := zipWriter.Create(fmt.Sprintf("%s_%d.mp4", videoMetadata.VideoName, index))
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

		// Increment counters
		index++
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

func DeleteVideoByID(c *gin.Context) {
	// Get video_id from query parameters
	videoID := c.Query("video_id")
	if videoID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "video_id is required"})
		return
	}

	// Convert video_id to ObjectID
	objectID, err := primitive.ObjectIDFromHex(videoID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid video_id"})
		return
	}

	// Connect to MongoDB
	client := Mongo.GetMongoDB()
	collection := Mongo.GetCollection("videostore")
	bucket, err := gridfs.NewBucket(client.Database("Pametni-Paketnik-baza"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error creating GridFS bucket"})
		return
	}

	// Find and delete video metadata
	result := collection.FindOneAndDelete(context.TODO(), bson.M{"video_id": videoID})
	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"message": "Video not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error deleting video metadata", "error": result.Err().Error()})
		}
		return
	}

	// Delete video file from GridFS
	err = bucket.Delete(objectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error deleting video file from GridFS", "error": err.Error()})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{"message": "Video deleted successfully"})
}

func FlagVideo(c *gin.Context) {
	videoID := c.Query("video_id")
	userID := c.Query("user_id")
	if videoID == "" || userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "video_id and user_id are required"})
		return
	}

	// Connect to MongoDB
	collection := Mongo.GetCollection("videostore")
	flagsCollection := Mongo.GetCollection("flags")
	usersCollection := Mongo.GetCollection("users")

	// Verify the user exists in the users collection
	var user bson.M
	err := usersCollection.FindOne(context.TODO(), bson.M{"user_id": userID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"message": "Invalid user"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error verifying user", "error": err.Error()})
		}
		return
	}

	// Check if the user has already flagged this video
	var existingFlag bson.M
	err = flagsCollection.FindOne(context.TODO(), bson.M{"video_id": videoID, "user_id": userID}).Decode(&existingFlag)
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "You have already flagged this video"})
		return
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid video_id"})
		return
	}

	// Increment the flagged counter in the videostore collection
	_, err = collection.UpdateOne(
		context.TODO(),
		bson.M{"video_id": videoID},
		bson.M{"$inc": bson.M{"flagged": 1}},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error flagging the video", "error": err.Error()})
		return
	}

	// Insert the user_id and video_id into the flags collection
	_, err = flagsCollection.InsertOne(context.TODO(), bson.M{
		"video_id":   videoID,
		"user_id":    userID,
		"flagged_at": time.Now(), // Optional: track when the flag occurred
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error recording flag action", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Video flagged successfully"})
}
