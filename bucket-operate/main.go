package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"log"
)

//MakeBucket 创建bucket
//Bucket 创建的bucket名字
//ACL bucket权限 私有
//CreateBucketConfiguration bucket的区域 ap-northeast-1
func MakeBucket(client *s3.Client, name string) (bool, string) {
	defer func() (bool, string) {
		if err := recover(); err != nil {
			reason := fmt.Sprintf("Runtime panic caught: %v\n", err)
			return false, reason

		}

		return true, ""
	}()
	_, err := client.CreateBucket(context.TODO(), &s3.CreateBucketInput{
		Bucket:                    aws.String(name),
		ACL:                       types.BucketCannedACLPrivate,
		CreateBucketConfiguration: &types.CreateBucketConfiguration{LocationConstraint: types.BucketLocationConstraintApNortheast1},
	})

	if err != nil {
		reason := fmt.Sprintf("Runtime panic caught: %v\n", err)
		return false, reason
	}
	return true, ""

}

//DeleteBucket 删除bucket
func DeleteBucket(client *s3.Client, name string) (bool, string) {
	defer func() (bool, string) {
		if err := recover(); err != nil {
			reason := fmt.Sprintf("Runtime panic caught: %v\n", err)
			return false, reason

		}

		return true, ""
	}()
	_, err := client.DeleteBucket(context.TODO(), &s3.DeleteBucketInput{
		Bucket: aws.String(name),
	})

	if err != nil {
		reason := fmt.Sprintf("Runtime panic caught: %v\n", err)
		return false, reason
	}
	return true, ""

}
func main() {
	// Load the Shared AWS Configuration (~/.aws/config-file)
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	// Create an Amazon S3 service client
	client := s3.NewFromConfig(cfg)
	//这里的bucket需要全球唯一
	status, reason := MakeBucket(client, "my-bucketxxx00011122")
	if !status {
		fmt.Printf("创建bucket失败，失败原因是: %s\n\n", reason)
	} else {
		fmt.Println("创建bucket成功")
	}

	//这里的bucket需要全球唯一
	status, reason = DeleteBucket(client, "my-bucketxxx00011122")
	if !status {
		fmt.Printf("删除bucket失败，失败原因是: %s\n\n", reason)
	} else {
		fmt.Println("删除bucket成功")
	}

}
