package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

var urlStore = make(map[string]string)

func generateShortID() string {
	b := make([]byte, 4)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

func shortenURL(c *gin.Context) {
	var json struct {
		Original string `json:"original"`
	}
	if err := c.BindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	id := generateShortID()
	urlStore[id] = json.Original

	// получаем домен из оригинального URL (если нужно — из запроса)
	shortDomain := c.Request.Host // домен, на который пришёл запрос
	if u, err := url.Parse(json.Original); err == nil {
		shortDomain = u.Host // подставляем домен исходного URL
	}

	// итоговая короткая ссылка будет на том же домене
	shortURL := fmt.Sprintf("https://%s/%s", shortDomain, id)
	c.JSON(http.StatusOK, gin.H{"short": shortURL})
}

func redirectURL(c *gin.Context) {
	id := c.Param("short")
	if original, ok := urlStore[id]; ok {
		c.Redirect(http.StatusMovedPermanently, original)
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
	}
}

func main() {
	r := gin.Default()
	r.POST("/shorten", shortenURL)
	r.GET("/:short", redirectURL)
	r.Run(":8080")
}
