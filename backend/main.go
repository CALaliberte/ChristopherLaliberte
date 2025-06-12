package main

import (
    "encoding/json"
    "fmt"
    "io"
    "log"
    "net/http"
    "os"
    "path/filepath"
    "strings"

    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
)

type ChannelResponse struct {
    Items []struct {
        ContentDetails struct {
            RelatedPlaylists struct {
                Uploads string `json:"uploads"`
            } `json:"relatedPlaylists"`
        } `json:"contentDetails"`
    } `json:"items"`
}

type PlaylistResponse struct {
    Items []struct {
        Snippet struct {
            ResourceId struct {
                VideoId string `json:"videoId"`
            } `json:"resourceId"`
        } `json:"snippet"`
    } `json:"items"`
}

func main() {
    err := godotenv.Load()
    if err != nil {
        log.Printf("Error loading .env file: %v", err)
    }

    port := os.Getenv("PORT")
    if port == "" {
        port = "5000"
    }

    r := gin.Default()

    // Serve static files from the parent directory
    r.Static("/static", "../")
    // Redirect root to index.html
    r.GET("/", func(c *gin.Context) {
        c.File("../index.html")
    })

    r.GET("/youtube/channels", func(c *gin.Context) {
        apiKey := c.GetHeader("Authorization")
        if !strings.HasPrefix(apiKey, "Bearer ") {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header"})
            return
        }
        apiKey = strings.TrimPrefix(apiKey, "Bearer ")

        url := "https://www.googleapis.com/youtube/v3/channels?part=contentDetails&mine=true"
        req, err := http.NewRequest("GET", url, nil)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        req.Header.Set("Authorization", "Bearer "+apiKey)

        client := &http.Client{}
        resp, err := client.Do(req)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        defer resp.Body.Close()

        body, err := io.ReadAll(resp.Body)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        var channelResp ChannelResponse
        if err := json.Unmarshal(body, &channelResp); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        c.Data(resp.StatusCode, "application/json", body)
    })

    r.GET("/youtube/playlist", func(c *gin.Context) {
        apiKey := c.GetHeader("Authorization")
        if !strings.HasPrefix(apiKey, "Bearer ") {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header"})
            return
        }
        apiKey = strings.TrimPrefix(apiKey, "Bearer ")

        playlistId := c.Query("playlistId")
        if playlistId == "" {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Missing playlistId parameter"})
            return
        }

        url := fmt.Sprintf("https://www.googleapis.com/youtube/v3/playlistItems?part=snippet&playlistId=%s&maxResults=1", playlistId)
        req, err := http.NewRequest("GET", url, nil)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        req.Header.Set("Authorization", "Bearer "+apiKey)

        client := &http.Client{}
        resp, err := client.Do(req)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        defer resp.Body.Close()

        body, err := io.ReadAll(resp.Body)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        var playlistResp PlaylistResponse
        if err := json.Unmarshal(body, &playlistResp); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        c.Data(resp.StatusCode, "application/json", body)
    })

    r.GET("/pdf/:filename", func(c *gin.Context) {
        filename := c.Param("filename")
        filePath := filepath.Join("../SeniorJury", filename)

        if _, err := os.Stat(filePath); os.IsNotExist(err) {
            c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
            return
        }

        c.File(filePath)
    })

    r.Run(":" + port)
}
