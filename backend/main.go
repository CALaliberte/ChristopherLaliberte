package main

import (
	"crypto/subtle"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// Config holds environment variables
type Config struct {
	YouTubeAPIKey string
	Port          string
}

// LoadConfig from .env file or environment
func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("error loading .env: %w", err)
	}
	return &Config{
		YouTubeAPIKey: os.Getenv("YOUTUBE_API_KEY"),
		Port:          os.Getenv("PORT"),
	}, nil
}

// APIKeyMiddleware validates API key using constant-time comparison
func APIKeyMiddleware(config *Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, "Bearer ")
		if len(parts) != 2 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header"})
			c.Abort()
			return
		}

		providedKey := parts[1]
		if subtle.ConstantTimeCompare([]byte(providedKey), []byte(config.YouTubeAPIKey)) != 1 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid API key"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func main() {
	// Load configuration
	config, err := LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize Gin router
	r := gin.Default()

	// Enable CORS for GitHub Pages and local testing
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://calaliberte.github.io", "http://localhost:8000"},
		AllowMethods:     []string{"GET"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	// Health check endpoint for Render
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	// Stream PDF without allowing download
	r.GET("/pdf/:filename", func(c *gin.Context) {
		filename := c.Param("filename")
		filePath := fmt.Sprintf("./SeniorJury/%s", filename)

		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
			return
		}

		c.Header("Content-Type", "application/pdf")
		c.Header("Content-Disposition", "inline; filename="+filename)
		c.Header("Cache-Control", "no-store, no-cache, must-revalidate, proxy-revalidate")
		c.Header("Pragma", "no-cache")
		c.Header("Expires", "0")
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("Access-Control-Allow-Origin", "http://localhost:8000")
		c.Header("Access-Control-Allow-Methods", "GET")
		c.Header("Access-Control-Allow-Headers", "Content-Type")

		c.File(filePath)
	})

	// Proxy for YouTube API channels
	r.GET("/youtube/channels", APIKeyMiddleware(config), func(c *gin.Context) {
		channelID := "UCTPxOafLB3PA6Twtl18dw2g"
		url := fmt.Sprintf("https://www.googleapis.com/youtube/v3/channels?part=contentDetails&id=%s&key=%s", channelID, config.YouTubeAPIKey)
		log.Printf("Fetching YouTube channels: %s", url)

		resp, err := http.Get(url)
		if err != nil {
			log.Printf("Failed to fetch YouTube data: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch YouTube data"})
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			log.Printf("YouTube API error: %s", string(body))
			c.JSON(resp.StatusCode, gin.H{"error": string(body)})
			return
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Failed to read YouTube response: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read YouTube response"})
			return
		}

		if len(body) == 0 {
			log.Printf("Empty YouTube response for channelID: %s", channelID)
			c.JSON(http.StatusNotFound, gin.H{"error": "No channel data found"})
			return
		}

		c.Data(http.StatusOK, "application/json", body)
	})

	// Proxy for YouTube API playlist items
	r.GET("/youtube/playlist", APIKeyMiddleware(config), func(c *gin.Context) {
		playlistID := c.Query("playlistId")
		if playlistID == "" {
			log.Printf("Missing playlistId parameter")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing playlistId parameter"})
			return
		}
		url := fmt.Sprintf("https://www.googleapis.com/youtube/v3/playlistItems?part=snippet&maxResults=1&playlistId=%s&key=%s", playlistID, config.YouTubeAPIKey)
		log.Printf("Fetching YouTube playlist: %s", url)

		resp, err := http.Get(url)
		if err != nil {
			log.Printf("Failed to fetch YouTube playlist: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch YouTube playlist"})
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			log.Printf("YouTube playlist API error (status %d): %s", resp.StatusCode, string(body))
			c.JSON(resp.StatusCode, gin.H{"error": string(body)})
			return
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Failed to read YouTube playlist response: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read YouTube playlist response"})
			return
		}

		if len(body) == 0 {
			log.Printf("Empty YouTube playlist response for playlistID: %s", playlistID)
			c.JSON(http.StatusNoContent, gin.H{"error": "No playlist items found"})
			return
		}

		c.Data(http.StatusOK, "application/json", body)
	})

	// Start server
	port := config.Port
	if port == "" {
		port = "8080"
	}
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
