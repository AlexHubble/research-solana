package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/tencentyun/cos-go-sdk-v5"
)

func main() {
	var (
		operation = flag.String("op", "", "操作类型: upload 或 download")
		bucketURL = flag.String("url", "", "COS 存储桶 URL (例如: https://examplebucket-1250000000.cos.ap-guangzhou.myqcloud.com)")
		secretID  = flag.String("id", "", "腾讯云 SecretID")
		secretKey = flag.String("key", "", "腾讯云 SecretKey")
		objectKey = flag.String("object", "", "COS 对象键名（远程文件名）")
		localFile = flag.String("file", "", "本地文件路径")
	)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "腾讯云 COS 文件上传下载工具\n\n")
		fmt.Fprintf(os.Stderr, "使用方法:\n")
		fmt.Fprintf(os.Stderr, "  上传文件: %s -op=upload -url=<bucket_url> -id=<secret_id> -key=<secret_key> -object=<object_key> -file=<local_file>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  下载文件: %s -op=download -url=<bucket_url> -id=<secret_id> -key=<secret_key> -object=<object_key> -file=<local_file>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\n参数说明:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n示例:\n")
		fmt.Fprintf(os.Stderr, "  上传: %s -op=upload -url=https://mybucket-1250000000.cos.ap-guangzhou.myqcloud.com -id=AKIDxxxxx -key=xxxxx -object=test.txt -file=./test.txt\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  下载: %s -op=download -url=https://mybucket-1250000000.cos.ap-guangzhou.myqcloud.com -id=AKIDxxxxx -key=xxxxx -object=test.txt -file=./downloaded.txt\n", os.Args[0])
	}

	flag.Parse()

	// 验证必需参数
	if *operation == "" || *bucketURL == "" || *secretID == "" || *secretKey == "" || *objectKey == "" || *localFile == "" {
		fmt.Fprintf(os.Stderr, "错误: 所有参数都是必需的\n\n")
		flag.Usage()
		os.Exit(1)
	}

	if *operation != "upload" && *operation != "download" {
		fmt.Fprintf(os.Stderr, "错误: 操作类型必须是 'upload' 或 'download'\n\n")
		flag.Usage()
		os.Exit(1)
	}

	// 创建 COS 客户端
	client, err := createCOSClient(*bucketURL, *secretID, *secretKey)
	if err != nil {
		fmt.Fprintf(os.Stderr, "创建 COS 客户端失败: %v\n", err)
		os.Exit(1)
	}

	// 执行操作
	switch *operation {
	case "upload":
		err = uploadFile(client, *objectKey, *localFile)
	case "download":
		err = downloadFile(client, *objectKey, *localFile)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "操作失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("操作成功完成!\n")
}

// createCOSClient 创建 COS 客户端
func createCOSClient(bucketURL, secretID, secretKey string) (*cos.Client, error) {
	u, err := url.Parse(bucketURL)
	if err != nil {
		return nil, fmt.Errorf("解析存储桶 URL 失败: %v", err)
	}

	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  secretID,
			SecretKey: secretKey,
		},
	})

	return client, nil
}

// uploadFile 上传文件到 COS
func uploadFile(client *cos.Client, objectKey, localFilePath string) error {
	// 检查本地文件是否存在
	if _, err := os.Stat(localFilePath); os.IsNotExist(err) {
		return fmt.Errorf("本地文件不存在: %s", localFilePath)
	}

	fmt.Printf("开始上传文件: %s -> %s\n", localFilePath, objectKey)

	// 打开本地文件
	file, err := os.Open(localFilePath)
	if err != nil {
		return fmt.Errorf("打开本地文件失败: %v", err)
	}
	defer file.Close()

	// 获取文件信息
	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("获取文件信息失败: %v", err)
	}

	fmt.Printf("文件大小: %d 字节\n", fileInfo.Size())

	// 上传文件
	_, err = client.Object.Put(context.Background(), objectKey, file, nil)
	if err != nil {
		return fmt.Errorf("上传文件失败: %v", err)
	}

	fmt.Printf("文件上传成功: %s\n", objectKey)
	return nil
}

// downloadFile 从 COS 下载文件
func downloadFile(client *cos.Client, objectKey, localFilePath string) error {
	fmt.Printf("开始下载文件: %s -> %s\n", objectKey, localFilePath)

	// 检查对象是否存在
	_, err := client.Object.Head(context.Background(), objectKey, nil)
	if err != nil {
		return fmt.Errorf("远程文件不存在或无法访问: %v", err)
	}

	// 创建本地文件目录（如果不存在）
	localDir := filepath.Dir(localFilePath)
	if err := os.MkdirAll(localDir, 0755); err != nil {
		return fmt.Errorf("创建本地目录失败: %v", err)
	}

	// 下载文件
	resp, err := client.Object.Get(context.Background(), objectKey, nil)
	if err != nil {
		return fmt.Errorf("下载文件失败: %v", err)
	}
	defer resp.Body.Close()

	// 创建本地文件
	localFile, err := os.Create(localFilePath)
	if err != nil {
		return fmt.Errorf("创建本地文件失败: %v", err)
	}
	defer localFile.Close()

	// 复制数据
	written, err := io.Copy(localFile, resp.Body)
	if err != nil {
		return fmt.Errorf("写入本地文件失败: %v", err)
	}

	fmt.Printf("文件下载成功: %s (%d 字节)\n", localFilePath, written)
	return nil
}
