package Functions

import (
	"archive/zip"
	"backend/Mongo"
	"backend/Schemas"
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
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

	// Get video title from form data (renamed from video_name to title)
	videoTitle := c.PostForm("title")
	if videoTitle == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "title is required"})
		return
	}

	// Query the database for existing videos with the same title
	collection := Mongo.GetCollection("videostore")
	cursor, err := collection.Find(context.TODO(), bson.M{"video_name": videoTitle})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error querying video title count"})
		return
	}
	defer cursor.Close(context.TODO())

	// Count existing videos with the same title
	count := 0
	for cursor.Next(context.TODO()) {
		count++
	}

	// Generate a unique filename by appending index and timestamp
	uniqueFilename := fmt.Sprintf("%s_%d_%d", videoTitle, count+1, time.Now().UnixNano())
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
		VideoName:   videoTitle, // Passed from the user
		Uploader:    c.PostForm("uploader_username"),
		Description: c.PostForm("description"),
		Tags:        c.PostFormArray("flags"), // Expecting 'flags' from the frontend
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
	videoID := c.Param("id")
	if videoID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "video_id is required"})
		return
	}
	fmt.Printf("Received video_id: %s\n", videoID)

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
		c.JSON(http.StatusBadRequest, gin.H{"message": "Video not found or flagged"})
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
		"Video Name: %s\nUploader: %s\nDescription: %s\nTags: %v\nPosted At: %s\nVideo ID: %s\nFlagged Count: %d\n",
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
	// Get video_id and uploader_username from query parameters
	videoID := c.Query("video_id")
	uploaderUsername := c.Query("uploader_username")

	if videoID == "" || uploaderUsername == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "video_id and uploader_username are required"})
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
	usersCollection := Mongo.GetCollection("users")
	bucket, err := gridfs.NewBucket(client.Database("Pametni-Paketnik-baza"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error creating GridFS bucket"})
		return
	}

	// Verify the user exists and check admin status
	var user bson.M
	err = usersCollection.FindOne(context.TODO(), bson.M{"username": uploaderUsername}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error retrieving user data", "error": err.Error()})
		return
	}

	// Check if the user is an admin
	isAdmin, _ := user["admin"].(bool)

	// Find the video and verify uploader
	var video bson.M
	err = collection.FindOne(context.TODO(), bson.M{"video_id": videoID}).Decode(&video)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"message": "Video not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error retrieving video metadata", "error": err.Error()})
		return
	}

	// Check if the uploader_username matches the uploader or if the user is an admin
	uploader, ok := video["uploader_username"].(string)
	if !ok || (uploader != uploaderUsername && !isAdmin) {
		c.JSON(http.StatusForbidden, gin.H{"message": "You are not authorized to delete this video"})
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

type VideoFlagRequest struct {
	VideoID string `json:"video_id"`
	UserID  string `json:"user_id"`
}

func FlagVideo(c *gin.Context) {

	var request VideoFlagRequest
	if err := c.BindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	videoID := request.VideoID
	userID := request.UserID

	fmt.Printf("Received video_id: %s, user_id: %s\n", videoID, userID) // Log incoming parameters

	if videoID == "" || userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "video_id and user_id are required"})
		return
	}

	// Convert userID to ObjectId
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	fmt.Printf("userObjectID: %s\n", userObjectID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid user_id format"})
		return
	}

	// Connect to MongoDB
	collection := Mongo.GetCollection("videostore")
	usersCollection := Mongo.GetCollection("users")

	// Verify the user exists in the users collection
	var user bson.M
	err = usersCollection.FindOne(context.TODO(), bson.M{"_id": userObjectID}).Decode(&user)

	if errors.Is(err, mongo.ErrNoDocuments) {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"message": "Invalid user"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error verifying user", "error": err.Error()})
		return
	}

	// Atomically check if the user_id is already in the flagged_by array and add it if not
	updateResult, err := collection.UpdateOne(
		context.TODO(),
		bson.M{
			"video_id":   videoID,
			"flagged_by": bson.M{"$ne": userID}, // Ensure the user is not already in the flagged_by array
		},
		bson.M{
			"$inc":      bson.M{"flagged": 1},         // Increment the flagged count
			"$addToSet": bson.M{"flagged_by": userID}, // Add userID to the flagged_by array
		},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error flagging the video", "error": err.Error()})
		return
	}

	// Check if the user was already in the flagged_by array
	if updateResult.MatchedCount == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "You have already flagged this video"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Video flagged successfully"})
}

func GetFlaggedVideos(c *gin.Context) {
	client := Mongo.GetMongoDB()
	collection := Mongo.GetCollection("videostore")
	bucket, err := gridfs.NewBucket(client.Database("Pametni-Paketnik-baza"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error creating GridFS bucket"})
		return
	}

	// Query videos with flagged count greater than 3
	cursor, err := collection.Find(context.TODO(), bson.M{"flagged": bson.M{"$gt": 3}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error querying video metadata"})
		return
	}
	defer cursor.Close(context.TODO())

	// Set response headers for ZIP file
	c.Header("Content-Type", "application/zip")
	c.Header("Content-Disposition", `attachment; filename="flagged_videos.zip"`)

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
			"Video Name: %s\nUploader: %s\nDescription: %s\nTags: %v\nPosted At: %s\nVideo ID: %s\nFlagged Count: %d\n",
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
		c.JSON(http.StatusNotFound, gin.H{"message": "No flagged videos found in the database"})
		return
	}

	fmt.Printf("Total flagged videos added to ZIP: %d\n", videoCount)
	c.Status(http.StatusOK)
}

func ResetFlaggedCounter(c *gin.Context) {
	// Get video_id and admin_username from query parameters
	videoID := c.Query("video_id")
	adminUsername := c.Query("admin_username")

	if videoID == "" || adminUsername == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "video_id and admin_username are required"})
		return
	}

	// Connect to MongoDB
	collection := Mongo.GetCollection("videostore")
	usersCollection := Mongo.GetCollection("users")

	// Verify the user exists and is an admin
	var user bson.M
	err := usersCollection.FindOne(context.TODO(), bson.M{"username": adminUsername}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"message": "Admin user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error retrieving user data", "error": err.Error()})
		return
	}

	// Check if the user is an admin
	isAdmin, _ := user["admin"].(bool)
	if !isAdmin {
		c.JSON(http.StatusForbidden, gin.H{"message": "You are not authorized to reset the flagged counter"})
		return
	}

	// Reset the flagged counter of the video to 0
	updateResult, err := collection.UpdateOne(
		context.TODO(),
		bson.M{"video_id": videoID},
		bson.M{"$set": bson.M{"flagged": 0}},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error resetting flagged counter", "error": err.Error()})
		return
	}

	// Check if the video was found and updated
	if updateResult.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Video not found"})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{"message": "Flagged counter reset successfully"})
}
