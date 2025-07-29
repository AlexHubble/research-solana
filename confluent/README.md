# Confluent Kafka Demo (Go版本)

这是一个用于实验Confluent Kafka集群的demo项目，使用Go语言实现，包含producer和consumer示例代码。

## 项目结构

```
confluent/
├── config/
│   └── kafka.yaml          # Kafka配置文件
├── pkg/
│   └── config/
│       └── config.go       # 配置管理包
├── producer/
│   └── main.go             # 消息生产者
├── consumer/
│   └── main.go             # 消息消费者
├── go.mod                  # Go模块文件
├── go.sum                  # Go依赖校验文件
└── README.md              # 项目说明
```

## 安装依赖

```bash
go mod tidy
```

## 配置

1. 编辑 `config/kafka.yaml` 文件，填入你的Confluent Cloud配置：
   - `sasl_username`: 你的API Key
   - `sasl_password`: 你的API Secret
   - `bootstrap_servers`: 你的集群地址
   - 其他配置根据需要调整

## 使用方法

### 启动Consumer（建议先启动）

```bash
cd consumer
go run main.go
```

可选参数：
- `-config`: 配置文件路径（默认: ../config/kafka.yaml）
- `-timeout`: 消息轮询超时时间（默认: 1s）

### 启动Producer

```bash
cd producer
go run main.go
```

可选参数：
- `-config`: 配置文件路径（默认: ../config/kafka.yaml）
- `-count`: 发送消息数量（默认: 10）
- `-interval`: 发送间隔时间（默认: 1s）

## 示例

### 发送10条消息，每秒1条
```bash
cd producer
go run main.go -count 10 -interval 1s
```

### 发送100条消息，每0.5秒1条
```bash
cd producer
go run main.go -count 100 -interval 500ms
```

### 使用自定义配置文件
```bash
cd consumer
go run main.go -config /path/to/your/config.yaml
```

### 编译后运行
```bash
# 编译producer
cd producer
go build -o producer main.go
./producer -count 5 -interval 2s

# 编译consumer
cd consumer
go build -o consumer main.go
./consumer -timeout 2s
```

## 功能特性

### Producer
- ✅ 从YAML配置文件读取Kafka连接信息
- ✅ 支持自定义发送消息数量和间隔
- ✅ JSON格式消息，包含时间戳和序列号
- ✅ 消息发送状态回调
- ✅ 错误处理和重试机制

### Consumer
- ✅ 从YAML配置文件读取Kafka连接信息
- ✅ 优雅退出处理（Ctrl+C）
- ✅ JSON消息解析和格式化输出
- ✅ 消息统计和错误处理
- ✅ 支持consumer group

## 注意事项

1. 确保你的Confluent Cloud集群已经创建并运行
2. 确保API Key有足够的权限读写指定的topic
3. 如果topic不存在，Confluent Cloud会自动创建（如果开启了自动创建功能）
4. Consumer使用`earliest`策略，会从最早的消息开始消费

## 故障排除

### 连接问题
- 检查网络连接
- 验证API Key和Secret是否正确
- 确认集群地址是否正确

### 权限问题
- 确保API Key有读写topic的权限
- 检查consumer group权限

### 配置问题
- 验证YAML文件格式是否正确
- 检查所有必需的配置项是否填写
