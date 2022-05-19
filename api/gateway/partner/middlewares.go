package partner

import (
	"backend/api/gateway/client"
	pb "backend/proto/partner"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
	"time"
)

var whiteList = [...]string{".jpg", ".jpeg", ".png"}

const maxImageSize = 3 * 1024 * 1024

func CheckFile(ctx *gin.Context) {
	_, header, err := ctx.Request.FormFile("image")
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "file not found",
		})
		return
	}
	if header.Size > maxImageSize {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "too large file",
		})
		return
	}
	extension := filepath.Ext(header.Filename)
	isValidateType := false
	for idx := range whiteList {
		white := whiteList[idx]
		if extension == white {
			isValidateType = true
			break
		}
	}
	if !isValidateType {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"message": "invalidate file",
		})
		return
	}
	ctx.Next()
}

func Upload(ctx *gin.Context) {
	file, header, err := ctx.Request.FormFile("image")
	name := ctx.Request.FormValue("name")
	description := ctx.Request.FormValue("description")
	bytes := make([]byte, header.Size)
	_, err = file.Read(bytes)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "fail read file",
		})
	}
	timeout, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()
	req := pb.SaveRequest{
		File: &pb.ImageFile{
			Chunk: bytes,
		},
		Name:        name,
		Description: description,
	}
	res, err := client.Partner().SaveNFTInfo(timeout, &req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "internal server error",
		})
		return
	}
	if !res.Success {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"message": "fail write",
		})
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
		"id":      res.Id,
	})
}

func GetNFTInfo(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "has not id",
		})
		return
	}
	req := pb.LoadRequest{Id: id}
	timeout, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()
	res, err := client.Partner().LoadNFTInfo(timeout, &req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "internal server error",
		})
		return
	}
	switch int(res.Status.Code) {
	case http.StatusBadRequest:
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "wrong request",
		})
	case http.StatusNotFound:
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message": "not found",
		})
	case http.StatusOK:
		ctx.JSON(http.StatusOK, gin.H{
			"message":     "success",
			"name":        res.Name,
			"description": res.Description,
			"image":       res.Url,
		})
	default:
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "internal error",
		})
	}

}
