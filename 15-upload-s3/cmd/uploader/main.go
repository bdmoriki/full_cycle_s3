package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var (
	s3Client *s3.Client
	s3Bucket string
)

func init() {
	ctx := context.Background()
	sdkConfig, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		fmt.Println("Couldn't load default configuration. Have you set up your AWS account?")
		fmt.Println(err)
		return
	}
	s3Client = s3.NewFromConfig(sdkConfig)
	s3Bucket = "goexpert-bucket-exemplo"
}

func main() {
	dir, err := os.Open("./tmp")
	if err != nil {
		panic(err)
	}
	defer dir.Close()

	for {
		files, err := dir.ReadDir(1)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Printf("Erros reading directory %s\n:", err)
			continue
		}

		uploadFile(files[0].Name())
	}
}

func uploadFile(fileName string) {
	ctx := context.Background()
	completeFileName := fmt.Sprintf("./tmp/%s", fileName)
	fmt.Printf("Uploading file %s to bucket %s\n:", fileName, s3Bucket)
	f, err := os.Open(completeFileName)
	if err != nil {
		fmt.Printf("Error opening file %s:\n", fileName)
		return
	}
	defer f.Close()

	_, err = s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s3Bucket),
		Key:    aws.String(fileName),
		Body:   f,
	})
	if err != nil {
		fmt.Printf("Error uploading file %s:\n", fileName)
		return
	}
	fmt.Printf("File %s upload successfully\n:", fileName)
}
