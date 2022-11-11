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
	"time"
)

func SignUrl(client *s3.Client, bucketName, Key string) string {
	defer func() (bool, string) {
		if err := recover(); err != nil {
			reason := fmt.Sprintf("Runtime panic caught: %v\n", err)
			return false, reason

		}

		return true, ""
	}()
	fmt.Println("Create Presign client")
	presignClient := s3.NewPresignClient(client)

	presignParams := &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(Key),
	}

	// Apply an expiration via an option function
	presignDuration := func(po *s3.PresignOptions) {
		po.Expires = 5 * time.Minute
	}

	presignResult, err := presignClient.PresignGetObject(context.TODO(), presignParams, presignDuration)

	if err != nil {
		panic("Couldn't get presigned URL for GetObject")
	}

	fmt.Printf("Presigned URL For object: %s\n", presignResult.URL)
	return presignResult.URL

}
func main() {
	// Load the Shared AWS Configuration (~/.aws/config)
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	// Create an Amazon S3 service client
	client := s3.NewFromConfig(cfg)

	signUrl := SignUrl(client, "uniquebucket1233", "iu.jpeg")

	fmt.Printf("签名后的url地址为:%s", signUrl)
}
