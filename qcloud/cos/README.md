# 腾讯云 COS 文件上传下载工具

这是一个使用腾讯云 COS (Cloud Object Storage) SDK 的示例程序，支持通过命令行进行文件上传和下载操作。

## 功能特性

- 支持文件上传到腾讯云 COS
- 支持从腾讯云 COS 下载文件
- 所有参数通过命令行传入
- 完整的错误处理和状态提示
- 自动创建本地目录（下载时）

## 安装依赖

```bash
go mod init cos-client
go mod tidy
```

## 使用方法

### 编译程序

```bash
go build -o cos_client cos_client.go
```

### 上传文件

```bash
./cos_client -op=upload \
  -url=https://your-bucket-1250000000.cos.ap-guangzhou.myqcloud.com \
  -id=your_secret_id \
  -key=your_secret_key \
  -object=remote/path/filename.txt \
  -file=./local/file.txt
```

### 下载文件

```bash
./cos_client -op=download \
  -url=https://your-bucket-1250000000.cos.ap-guangzhou.myqcloud.com \
  -id=your_secret_id \
  -key=your_secret_key \
  -object=remote/path/filename.txt \
  -file=./downloaded/file.txt
```

## 参数说明

| 参数 | 说明 | 示例 |
|------|------|------|
| `-op` | 操作类型，`upload` 或 `download` | `upload` |
| `-url` | COS 存储桶 URL | `https://mybucket-1250000000.cos.ap-guangzhou.myqcloud.com` |
| `-id` | 腾讯云 SecretID | `AKIDxxxxxxxxxxxxxxxxxxxxx` |
| `-key` | 腾讯云 SecretKey | `xxxxxxxxxxxxxxxxxxxxxxxx` |
| `-object` | COS 对象键名（远程文件路径） | `folder/file.txt` |
| `-file` | 本地文件路径 | `./local/file.txt` |

## 获取腾讯云密钥

1. 登录 [腾讯云控制台](https://console.cloud.tencent.com/)
2. 进入 [访问管理 - API密钥管理](https://console.cloud.tencent.com/cam/capi)
3. 创建或查看现有的 SecretID 和 SecretKey

## 存储桶 URL 格式

存储桶 URL 格式为：`https://<BucketName-APPID>.cos.<Region>.myqcloud.com`

- `BucketName`: 存储桶名称
- `APPID`: 腾讯云账户的 APPID
- `Region`: 存储桶所在地域，如 `ap-guangzhou`、`ap-beijing` 等

## 示例

### 上传示例

```bash
# 上传本地文件 test.txt 到 COS 的 documents/test.txt
./cos_client -op=upload \
  -url=https://mybucket-1250000000.cos.ap-guangzhou.myqcloud.com \
  -id=AKIDxxxxxxxxxxxxxxxxxxxxx \
  -key=xxxxxxxxxxxxxxxxxxxxxxxx \
  -object=documents/test.txt \
  -file=./test.txt
```

### 下载示例

```bash
# 从 COS 下载 documents/test.txt 到本地 downloads/test.txt
./cos_client -op=download \
  -url=https://mybucket-1250000000.cos.ap-guangzhou.myqcloud.com \
  -id=AKIDxxxxxxxxxxxxxxxxxxxxx \
  -key=xxxxxxxxxxxxxxxxxxxxxxxx \
  -object=documents/test.txt \
  -file=./downloads/test.txt
```

## 错误处理

程序包含完整的错误处理：

- 参数验证
- 文件存在性检查
- 网络连接错误
- 权限验证错误
- 文件读写错误

## 注意事项

1. 确保 SecretID 和 SecretKey 有足够的权限访问指定的存储桶
2. 存储桶 URL 必须正确，包含正确的地域信息
3. 下载时会自动创建本地目录
4. 上传时会覆盖同名的远程文件
5. 建议使用子账号密钥，遵循最小权限原则

## 相关链接

- [腾讯云 COS 官方文档](https://cloud.tencent.com/document/product/436)
- [COS Go SDK 文档](https://cloud.tencent.com/document/product/436/31215)
- [腾讯云控制台](https://console.cloud.tencent.com/cos5/bucket)
