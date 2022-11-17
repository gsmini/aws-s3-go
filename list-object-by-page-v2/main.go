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

// https://docs.aws.amazon.com/zh_cn/AmazonS3/latest/API/API_ListObjectsV2.html
func listObjectPage(client *s3.Client, pageSize int32, startAfter string, prefix string) (data []string, count int32) {
	start := time.Now() // 获取当前时间
	bucket := "uniquebucket1233"
	//s3 的页不是和平时的api page page-size一样，这里是利用StartAfter方式分页的，也就是"从这个文件开始查找后面的"
	outPut, err := client.ListObjectsV2(context.TODO(),
		&s3.ListObjectsV2Input{
			Bucket:     aws.String(bucket),
			Prefix:     aws.String(prefix),
			MaxKeys:    pageSize,
			StartAfter: aws.String(startAfter)},
	)
	if err != nil {
		fmt.Printf("查询错误%s", err.Error())
	}

	//go build && ./main
	fmt.Println("开始时间：", start)
	fmt.Println("该函数执行完成耗时：", time.Since(start))
	fmt.Println("两天获取的数据大小为", len(outPut.Contents))

	for _, item := range outPut.Contents {
		data = append(data, *item.Key)
		//fmt.Println(*item.Key)
	}
	return data, int32(len(outPut.Contents))
}
func main() {
	var pageSize int32 = 500
	var count int32 = 0
	var result []string
	startAfter := ""
	prefix := "iot/xdasj12dxj/2022/11/13/"
	// Load the Shared AWS Configuration (~/.aws/config)
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	// Create an Amazon S3 service client
	client := s3.NewFromConfig(cfg)

	//假设这是前端访问：每页pageSize个，如果返回的结果小于pageSize，说明已经是最后一页
	for {
		data, dataCount := listObjectPage(client, pageSize, startAfter, prefix)
		count = dataCount
		startAfter = data[count-1] //取最后一个作为筛选过滤的s3 key传递给startAfter
		fmt.Printf("当前获取的数据长度为:%d", dataCount)
		result = append(result, data...)
		if count < pageSize {
			break
		}
	}
	fmt.Println(len(result))

}
