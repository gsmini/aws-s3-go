package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"time"
)

func listObjectPage(pageSize int64, startAfter string, prefix string) (data []string, count int64) {
	start := time.Now() // 获取当前时间
	bucket := ""
	region := ""
	awsAccessKeyId := " "
	awsSecretAccessKey := ""

	sess, _ := session.NewSession(&aws.Config{
		Region:      aws.String(region),                                                       //桶所在的区域
		Credentials: credentials.NewStaticCredentials(awsAccessKeyId, awsSecretAccessKey, ""), //sts的临时凭证
	})

	svc := s3.New(sess)

	StartAfter := aws.String(startAfter)

	svc.ListObjectsV2Pages(&s3.ListObjectsV2Input{Bucket: aws.String(bucket), Prefix: &prefix, MaxKeys: &pageSize, StartAfter: StartAfter},
		func(page *s3.ListObjectsV2Output, lastPage bool) bool {
			fmt.Println(*page.KeyCount)
			for _, item := range page.Contents {
				//fmt.Println("Name:         ", *item.Key)
				data = append(data, *item.Key)

			}

			return false //翻页
		})

	elapsed := time.Since(start)
	//go build && ./main
	fmt.Println("开始时间：", start)
	fmt.Println("该函数执行完成耗时：", elapsed)
	fmt.Println("两天获取的数据大小为", len(data))
	return data, int64(len(data))
}
func main() {
	var pageSize int64 = 500
	var count int64 = 0
	var result []string
	startAfter := ""
	prefix := "iot/xdasj12dxj/2022/11/13/"

	//假设这是前段访问：每页pageSize个，如果返回的结果小于pageSize，说明已经是最后一页
	for {
		data, dataCount := listObjectPage(pageSize, startAfter, prefix)
		count = dataCount
		startAfter = data[count-1] //取最后一个作为筛选过滤的s3 key
		fmt.Printf("当前获取的数据长度为:%d", dataCount)
		result = append(result, data...)
		if count < pageSize {
			break
		}

	}
	fmt.Println(len(result))

}
