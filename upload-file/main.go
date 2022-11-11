//配置文件方式调用sdk
//

package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"log"
	"os"
)

//UploadFile 单个文件上传函数
//bucketName 目标上传的bucket名字
//filePath 本地文件地址
//Key 上传到s3的目标地址
func UploadFile(client *s3.Client, bucketName, filePath, Key string) (bool, string) {
	defer func() (bool, string) {
		if err := recover(); err != nil {
			reason := fmt.Sprintf("Runtime panic caught: %v\n", err)
			return false, reason

		}

		return true, ""
	}()
	// Place an object in a bucket.
	fmt.Println("Upload an object to the bucket")
	// Get the object body to upload.
	// Image credit: https://unsplash.com/photos/iz58d89q3ss
	stat, err := os.Stat(filePath)
	if err != nil {
		return false, fmt.Sprintf("%s", err.Error())
	}
	file, err := os.Open(filePath)

	if err != nil {
		return false, fmt.Sprintf("%s", err.Error())

	}

	_, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:        aws.String(bucketName),
		Key:           aws.String(Key),
		Body:          file,
		ContentLength: stat.Size(),
	})

	file.Close()

	if err != nil {
		return false, fmt.Sprintf("%s", err.Error())
	}

	return true, ""

}
func main() {
	// Load the Shared AWS Configuration (~/.aws/config)
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	// Create an Amazon S3 service client
	client := s3.NewFromConfig(cfg)
	status, reason := UploadFile(client, "uniquebucket1233", "iu.jpeg", "uploads/user/pic.jpeg")

	if !status {
		fmt.Printf("上传文件失败，失败原因是: %s\n\n", reason)
	} else {
		fmt.Println("上传文件成功")
	}
}
