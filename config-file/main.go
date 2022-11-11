//配置文件方式调用sdk
//

package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"log"
)

/*
https://docs.aws.amazon.com/zh_cn/sdkref/latest/guide/file-format.html#file-format-creds
在你的mac/linux下面~/.aws/config写入下面环境变量

[default]
aws_access_key_id="xxx"
aws_secret_access_key="xxx"
region="xxx

*/
func main() {
	bucket := ""
	// Load the Shared AWS Configuration (~/.aws/config-file)
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	// Create an Amazon S3 service client
	client := s3.NewFromConfig(cfg)

	// Get the first page of results for ListObjectsV2 for a bucket
	output, err := client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("first page results:")
	for _, object := range output.Contents {
		log.Printf("key=%s size=%d", aws.ToString(object.Key), object.Size)
	}
}
