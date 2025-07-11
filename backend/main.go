package main

import (
    "encoding/json"
    "fmt"
    "io"
    "log"
    "net/http"
    "net/url"
    "os"
    "path/filepath"

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
    apiKey := os.Getenv("YOUTUBE_API_KEY")

    r := gin.Default()

    // Serve static files
    r.Static("/SeniorJury", "./SeniorJury")
    r.GET("/", func(c *gin.Context) {
        c.File("../index.html")
    })
    r.GET("/main.js", func(c *gin.Context) {
        c.File("../main.js")
    })

    // YouTube latest video
    r.GET("/youtube/latest", func(c *gin.Context) {
        if apiKey == "" {
            c.JSON(http.StatusInternalServerError, gin.H{
                "error": "YouTube API key is missing. Please configure the API key to fetch the latest video.",
            })
            return
        }

        channelUrl := "https://www.googleapis.com/youtube/v3/channels?part=contentDetails&id=UCTPxOafLB3PA6Twtl18dw2g&key=" + apiKey
        resp, err := http.Get(channelUrl)
        if err != nil {
            c.Data(http.StatusInternalServerError, "text/html", []byte("<p class='text-gray-500 text-center'>Error fetching channel</p>"))
            return
        }
        defer resp.Body.Close()

        body, err := io.ReadAll(resp.Body)
        if err != nil {
            c.Data(http.StatusInternalServerError, "text/html", []byte("<p class='text-gray-500 text-center'>Error reading channel</p>"))
            return
        }

        var channelResp ChannelResponse
        if err := json.Unmarshal(body, &channelResp); err != nil {
            c.Data(http.StatusInternalServerError, "text/html", []byte("<p class='text-gray-500 text-center'>Error parsing channel</p>"))
            return
        }
        if len(channelResp.Items) == 0 {
            c.Data(http.StatusNotFound, "text/html", []byte("<p class='text-gray-500 text-center'>No channel data</p>"))
            return
        }
        uploadsPlaylistId := channelResp.Items[0].ContentDetails.RelatedPlaylists.Uploads

        playlistUrl := fmt.Sprintf("https://www.googleapis.com/youtube/v3/playlistItems?part=snippet&playlistId=%s&maxResults=1&key=%s", uploadsPlaylistId, apiKey)
        resp, err = http.Get(playlistUrl)
        if err != nil {
            c.Data(http.StatusInternalServerError, "text/html", []byte("<p class='text-gray-500 text-center'>Error fetching playlist</p>"))
            return
        }
        defer resp.Body.Close()

        body, err = io.ReadAll(resp.Body)
        if err != nil {
            c.Data(http.StatusInternalServerError, "text/html", []byte("<p class='text-gray-500 text-center'>Error reading playlist</p>"))
            return
        }

        var playlistResp PlaylistResponse
        if err := json.Unmarshal(body, &playlistResp); err != nil {
            c.Data(http.StatusInternalServerError, "text/html", []byte("<p class='text-gray-500 text-center'>Error parsing playlist</p>"))
            return
        }
        if len(playlistResp.Items) == 0 {
            c.Data(http.StatusNotFound, "text/html", []byte("<p class='text-gray-500 text-center'>No videos found</p>"))
            return
        }
        videoId := playlistResp.Items[0].Snippet.ResourceId.VideoId

        c.Data(http.StatusOK, "text/html", []byte(fmt.Sprintf(`
            <iframe id="latest-video" title="Latest YouTube Video" src="https://www.youtube.com/embed/%s?enablejsapi=1&rel=0" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share" allowfullscreen></iframe>
        `, videoId)))
    })

    // PDF endpoint
    r.GET("/pdf/:filename", func(c *gin.Context) {
    filename := c.Param("filename")
    baseDir, err := os.Getwd()
    if err != nil {
        log.Printf("Error getting working directory: %v", err)
        c.Data(http.StatusInternalServerError, "text/html", []byte("<p class='text-gray-500 text-center'>Internal server error</p>"))
        return
    }
    filePath := filepath.Join("SeniorJury", filename)
    cleanPath := filepath.Join(baseDir, filePath)
    juryDir := filepath.Join(baseDir, "SeniorJury")
    log.Printf("Requested file: %s, Resolved path: %s, JuryDir: %s", filename, cleanPath, juryDir)

    fileInfo, err := os.Stat(cleanPath)
    if os.IsNotExist(err) {
        log.Printf("File not found: %s", cleanPath)
        c.Data(http.StatusNotFound, "text/html", []byte("<p class='text-gray-500 text-center'>PDF not found</p>"))
        return
    }
    if err != nil {
        log.Printf("Error checking file %s: %v", cleanPath, err)
        c.Data(http.StatusInternalServerError, "text/html", []byte("<p class='text-gray-500 text-center'>Internal server error</p>"))
        return
    }
    log.Printf("File found: %s, Size: %d bytes, Permissions: %s", cleanPath, fileInfo.Size(), fileInfo.Mode())

    if !filepath.HasPrefix(cleanPath, juryDir) {
        log.Printf("Invalid path: %s does not start with %s", cleanPath, juryDir)
        c.Data(http.StatusBadRequest, "text/html", []byte("<p class='text-gray-500 text-center'>Invalid path</p>"))
        return
    }

    encodedFilename := url.PathEscape(filename)
    c.Redirect(http.StatusFound, fmt.Sprintf("/SeniorJury/%s", encodedFilename))
})    
  r.Run(":" + port)
}
