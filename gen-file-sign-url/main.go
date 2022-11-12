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

func PresignPutUrl(client *s3.Client, bucketName, Key string) string {
	defer func() (bool, string) {
		if err := recover(); err != nil {
			reason := fmt.Sprintf("Runtime panic caught: %v\n", err)
			return false, reason

		}

		return true, ""
	}()
	presignClient := s3.NewPresignClient(client)
	putObjectInput := &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(Key),
	}
	// 操作过期时间
	// loger.Debug("BrowserPostPreSign()", zap.Int64("now", now))
	//expire := time.Now().Unix() + 5*60
	//expiration := time.Unix(expire, 0).In(time.UTC).Format("2006-01-02T15:04:05.000Z")
	presignDuration := func(po *s3.PresignOptions) {
		po.Expires = 5 * time.Minute //不设置默认900s

	}

	presignResult, err := presignClient.PresignPutObject(context.TODO(), putObjectInput, presignDuration)
	fmt.Printf("%s", presignResult.URL)
	if err != nil {
		panic("Couldn't get presigned URL for GetObject")
	}

	url := fmt.Sprintf("%s\n", presignResult.URL)
	return url

}

func main() {
	// Load the Shared AWS Configuration (~/.aws/config)
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	// Create an Amazon S3 service client
	client := s3.NewFromConfig(cfg)

	presignUrl := PresignPutUrl(client, "uniquebucket1233", "xxxx.jpeg")

	fmt.Printf("\n%s", presignUrl)
	//生成的链接采用requests.put访问上传文件 python demo
	/*
		import pprint
		import requests

		url = ""
		with open("iu.jpeg", 'rb') as f:
		    files = {'file': ("xxxx.jpeg", f)}
		    http_response = requests.put(url,  files=files)
		    pprint.pprint(http_response.text)
	*/
}
