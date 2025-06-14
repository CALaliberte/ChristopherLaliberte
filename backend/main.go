package main

import (
    "encoding/json"
    "fmt"
    "io"
    "log"
    "net/http"
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
    r.Static("/static", "../")
    r.GET("/", func(c *gin.Context) {
        c.File("../index.html")
    })

    // YouTube latest video
    r.GET("/youtube/latest", func(c *gin.Context) {
        channelUrl := "https://www.googleapis.com/youtube/v3/channels?part=contentDetails&channelId=UCTPxOafLB3PA6Twtl18dw2g&key=" + apiKey
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
        filePath := filepath.Join("../SeniorJury", filename)

        if _, err := os.Stat(filePath); os.IsNotExist(err) {
            c.Data(http.StatusNotFound, "text/html", []byte("<p class='text-gray-500 text-center'>PDF not found</p>"))
            return
        }

        c.Data(http.StatusOK, "text/html", []byte(fmt.Sprintf(`
            <iframe src="/static/SeniorJury/%s" class="w-full h-full" frameborder="0" allow="fullscreen"></iframe>
        `, filename)))
    })

    r.Run(":" + port)
}
