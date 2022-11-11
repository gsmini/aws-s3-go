package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

const (
	maxPartSize   = int64(5 * 1024 * 1024) //每次最大上传size [5MiB ,5GiB]
	maxRetries    = 3                      //重试次数
	awsBucketName = "uniquebucket1233"
	filePath      = "xxxxx.mp4"
)

func main() {
	// Load the Shared AWS Configuration (~/.aws/config)
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	// Create an Amazon S3 service client
	client := s3.NewFromConfig(cfg)

	//<1>简单当前文件信息
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("打开文件失败: %s", err)
		return
	}
	defer file.Close()
	fileInfo, _ := file.Stat()
	size := fileInfo.Size() //文件总大小
	buffer := make([]byte, size)
	fileType := http.DetectContentType(buffer) //文件类型
	file.Read(buffer)

	key := "/media/" + file.Name() //目标上传目录
	//<2>申请创建分段上传 其中比较重要的是	resp.UploadId 表示此次分段上传操作的唯一id
	resp, err := client.CreateMultipartUpload(context.TODO(), &s3.CreateMultipartUploadInput{
		Bucket:      aws.String(awsBucketName),
		Key:         aws.String(key),
		ContentType: aws.String(fileType),
	})

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("开始分段上传文件....")

	var curr, partLength int64
	var remaining = size
	var completedParts []types.CompletedPart
	var partNumber int32 = 1
	for curr = 0; remaining != 0; curr += partLength {
		//如果剩下的文件大小 小于单次上传的文件大小，则本次上传的文件大小为remaining
		//否则为固定maxPartSize
		if remaining < maxPartSize {
			partLength = remaining
		} else {
			partLength = maxPartSize
		}
		//<3>上传某一部分文件 返回当前文件的上传信息
		completedPart, err := uploadPart(client, resp, buffer[curr:curr+partLength], partNumber)
		if err != nil {
			fmt.Printf("上传文件失败，错误内容:%s\n", err.Error())
			//上传失败要取消否则还是会收费
			err := abortMultipartUpload(client, resp)
			if err != nil {
				fmt.Println(err.Error())
			}
			return
		}
		fmt.Printf("上传文件成功:%d-%d\n", curr, curr+partLength)
		remaining -= partLength
		partNumber++
		completedParts = append(completedParts, completedPart)
	}

	_, err = completeMultipartUpload(client, resp, completedParts)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Printf("完成文件的全部上传")
}

//uploadPart 上传单个分片文件
//fileBytes当前文件字节
//partNumber 当前分片标记
//types.CompletedPart 当前分片完成上传后的分片信息，里面的etag信息到时候要收集上传到s3
//给s3做验证
func uploadPart(client *s3.Client, resp *s3.CreateMultipartUploadOutput, fileBytes []byte, partNumber int32) (types.CompletedPart, error) {
	tryNum := 1
	partInput := &s3.UploadPartInput{
		Body:          bytes.NewReader(fileBytes),
		Bucket:        resp.Bucket,
		Key:           resp.Key,
		PartNumber:    partNumber,
		UploadId:      resp.UploadId,
		ContentLength: int64(len(fileBytes)),
	}

	for tryNum <= maxRetries {
		uploadResult, err := client.UploadPart(context.TODO(), partInput)
		if err != nil {
			if tryNum == maxRetries {
				return types.CompletedPart{}, err
			}
			fmt.Printf("Retrying to upload part #%v\n", partNumber)
			tryNum++
		} else {
			fmt.Printf("Uploaded part #%v\n", partNumber)
			//CompletedPart中etag和PartNumber必须要传，因为s3要做校验
			return types.CompletedPart{
				ETag: uploadResult.ETag,
				//ChecksumCRC32:  uploadResult.ChecksumCRC32,
				//ChecksumCRC32C: uploadResult.ChecksumCRC32C,
				//ChecksumSHA1:   uploadResult.ChecksumSHA1,
				//ChecksumSHA256: uploadResult.ChecksumSHA256,
				PartNumber: partNumber,
			}, nil
		}
	}
	return types.CompletedPart{}, nil
}

//completeMultipartUpload 全部分片完成上传后调用次函数 显式表面当前操作完成
func completeMultipartUpload(client *s3.Client, resp *s3.CreateMultipartUploadOutput, completedParts []types.CompletedPart) (*s3.CompleteMultipartUploadOutput, error) {
	completeInput := &s3.CompleteMultipartUploadInput{
		Bucket:   resp.Bucket,
		Key:      resp.Key,
		UploadId: resp.UploadId,
		MultipartUpload: &types.CompletedMultipartUpload{
			Parts: completedParts,
		},
	}
	fmt.Printf("完成上传")
	return client.CompleteMultipartUpload(context.TODO(), completeInput)
}

//删除此次分段上传的全部文件，如果不删除还是会产生费用
func abortMultipartUpload(svc *s3.Client, resp *s3.CreateMultipartUploadOutput) error {
	fmt.Println("Aborting multipart upload for UploadId#" + *resp.UploadId)
	abortInput := &s3.AbortMultipartUploadInput{
		Bucket:   resp.Bucket,
		Key:      resp.Key,
		UploadId: resp.UploadId,
	}
	_, err := svc.AbortMultipartUpload(context.TODO(), abortInput)
	return err
}
