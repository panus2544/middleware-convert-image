package main

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"net/http"
	"net/url"
	"time"

	"github.com/chai2010/webp"
	"github.com/gin-gonic/gin"
	"github.com/nfnt/resize"
)

type ImgParams struct {
	FileFormat string `uri:"format" binding:"required"`
	Size       int    `uri:"size" binding:"required"`
}

func IsValidFileFormat(fileFormat string) bool {
	switch fileFormat {
	case
		"jpeg",
		"webp":
		return true
	}
	return false
}

func IsValidSize(size int) bool {
	switch size {
	case
		150,
		300,
		500,
		800,
		1024,
		1200,
		2048:
		return true
	}
	return false
}

func main() {
	r := gin.Default()
	r.GET("/:format/:size", func(c *gin.Context) {
		now := time.Now()
		fmt.Println(now)

		var params ImgParams
		if err := c.ShouldBindUri(&params); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"messages": err})
			return
		}

		if !IsValidFileFormat(params.FileFormat) {
			c.JSON(http.StatusBadRequest, gin.H{"messages": "Invalid image format"})
			return
		}

		if !IsValidSize(params.Size) {
			c.JSON(http.StatusBadRequest, gin.H{"messages": "Invalid image size"})
			return
		}

		originalUrl := c.DefaultQuery("url", "")
		if originalUrl == "" {
			c.JSON(http.StatusNotFound, gin.H{"messages": "File not found"})
			return
		}

		_, err := url.Parse(originalUrl)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"messages": "Incorrect URL"})
			return
		}

		resp, err := http.Get(originalUrl)
		if err != nil || resp.StatusCode != http.StatusOK {
			c.JSON(http.StatusNotFound, gin.H{"messages": "File not found"})
			return
		}

		reader := resp.Body
		defer reader.Close()
		contentType := resp.Header.Get("Content-Type")

		var img image.Image
		if contentType == "image/png" {
			img, err = png.Decode(resp.Body)
		} else {
			img, err = jpeg.Decode(resp.Body)
		}

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"messages": err})
			return
		}

		m := resize.Resize(uint(params.Size), 0, img, resize.Lanczos3)
		out := bytes.NewBuffer([]byte{})

		var respnseContentType string
		if params.FileFormat == "webp" {
			webp.Encode(out, m, &webp.Options{Quality: 90})
			respnseContentType = "image/webp"
		} else {
			jpeg.Encode(out, m, &jpeg.Options{Quality: 90})
			respnseContentType = "image/jpeg"
		}

		c.Data(http.StatusOK, respnseContentType, out.Bytes())
	})

	// listen and serve on 0.0.0.0:3000
	r.Run(":3000")
}
