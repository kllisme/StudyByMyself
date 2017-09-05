package util

import (
	"path/filepath"
	"os"
	"fmt"
	"io"
	"crypto/md5"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/spf13/viper"
	"mime/multipart"
	"time"
	"strings"
	"github.com/juju/errors"
)

//
func Upload(formFile *multipart.FileHeader, ossObject string) (string, error) {
	fileName := formFile.Filename
	fileExt := filepath.Ext(fileName)
	fileName = string(time.Now().Unix()) + fileExt

	extStr := viper.GetString("resource.oss.ext")
	extList := strings.Split(extStr, ",")
	supportExt := false
	for _, ext := range extList {
		if fileExt == ext {
			supportExt = true
		}
	}
	if !supportExt {
		return "", errors.New("不支持的格式")
	}
	file, err := formFile.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	absPath := viper.GetString("resource.oss.tmpPath") + fileName
	out, err := os.OpenFile(absPath, os.O_WRONLY | os.O_CREATE, 0666)
	fmt.Println(out, err)
	if err != nil {
		return "", err
	}
	defer out.Close()

	fileInfo, err := out.Stat()
	if err != nil {
		return "", err
	}

	if fileInfo.Size() > int64(viper.GetInt("resource.oss.maxSize")) {
		return "", errors.New("文件超长")
	}

	io.Copy(out, file)
	tmpFile, err := os.Open(absPath)
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()
	md5h := md5.New()
	io.Copy(md5h, tmpFile)
	hashName := fmt.Sprintf("%x", md5h.Sum(nil))
	objectName := ossObject + hashName + fileExt
	shortPath := hashName + fileExt
	tmpFile, err = os.Open(absPath)
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	client, err := oss.New(viper.GetString("resource.oss.endpoint"), viper.GetString("resource.oss.accessKeyId"), viper.GetString("resource.oss.accessKeySecret"))
	if err != nil {
		return "", err
	}

	bucket, err := client.Bucket(viper.GetString("resource.oss.bucketName"))
	if err != nil {
		return "", err
	}

	err = bucket.PutObject(objectName, tmpFile)
	os.Remove(absPath)
	if err != nil {
		return "", err
	}
	return shortPath, nil
}
