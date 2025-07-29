package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/AlexHubble/research-solana/confluent/pkg/config"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

// Message 消息结构体
type Message struct {
	ID        int                    `json:"id"`
	Timestamp string                 `json:"timestamp"`
	Message   string                 `json:"message"`
	Producer  string                 `json:"producer"`
	Data      map[string]interface{} `json:"data"`
}

func main() {
	// 命令行参数
	var (
		configPath = flag.String("config", "../config/kafka.yaml", "配置文件路径")
		timeout    = flag.Duration("timeout", time.Second, "消息轮询超时时间")
	)
	flag.Parse()

	// 获取配置文件绝对路径
	if !filepath.IsAbs(*configPath) {
		wd, err := os.Getwd()
		if err != nil {
			log.Fatalf("❌ 获取工作目录失败: %v", err)
		}
		*configPath = filepath.Join(wd, *configPath)
	}

	// 加载配置
	fmt.Printf("📖 加载配置文件: %s\n", *configPath)
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("❌ 加载配置失败: %v", err)
	}

	// 创建Consumer
	fmt.Println("🚀 创建Kafka Consumer...")
	consumer, err := kafka.NewConsumer(cfg.GetConsumerConfigMap())
	if err != nil {
		log.Fatalf("❌ 创建Consumer失败: %v", err)
	}
	defer consumer.Close()

	// 订阅topic
	topic := cfg.Consumer.Topic
	fmt.Printf("📥 订阅topic: %s\n", topic)
	err = consumer.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		log.Fatalf("❌ 订阅topic失败: %v", err)
	}

	// 设置优雅退出处理
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("🔄 开始消费消息... (按Ctrl+C退出)")
	fmt.Println("=" + string(make([]byte, 49)))

	messageCount := 0
	running := true

	// 消费消息循环
	for running {
		select {
		case sig := <-sigChan:
			fmt.Printf("\n🛑 收到退出信号 %v，正在优雅关闭...\n", sig)
			running = false
		default:
			// 轮询消息
			msg, err := consumer.ReadMessage(*timeout)
			if err != nil {
				// 超时不是错误，继续轮询
				if err.(kafka.Error).Code() == kafka.ErrTimedOut {
					continue
				}
				log.Printf("❌ Consumer错误: %v", err)
				continue
			}

			// 处理消息
			if processMessage(msg) {
				messageCount++
				fmt.Printf("✅ 已处理消息总数: %d\n", messageCount)
			}
		}
	}

	fmt.Printf("\n📊 消费统计: 总共处理了 %d 条消息\n", messageCount)
	fmt.Println("🔒 关闭Consumer...")
	fmt.Println("👋 程序已退出")
}

// processMessage 处理接收到的消息
func processMessage(msg *kafka.Message) bool {
	fmt.Println("📨 收到消息:")
	fmt.Printf("   Topic: %s\n", *msg.TopicPartition.Topic)
	fmt.Printf("   Partition: %d\n", msg.TopicPartition.Partition)
	fmt.Printf("   Offset: %v\n", msg.TopicPartition.Offset)
	
	if msg.Key != nil {
		fmt.Printf("   Key: %s\n", string(msg.Key))
	} else {
		fmt.Printf("   Key: <nil>\n")
	}
	
	if msg.Timestamp.IsZero() {
		fmt.Printf("   Timestamp: <nil>\n")
	} else {
		fmt.Printf("   Timestamp: %v\n", msg.Timestamp)
	}

	// 尝试解析JSON消息
	var message Message
	if err := json.Unmarshal(msg.Value, &message); err != nil {
		fmt.Printf("❌ JSON解析错误: %v\n", err)
		fmt.Printf("   原始消息: %s\n", string(msg.Value))
		fmt.Println("-" + string(make([]byte, 49)))
		return false
	}

	// 格式化输出消息内容
	messageJSON, err := json.MarshalIndent(message, "   ", "  ")
	if err != nil {
		fmt.Printf("❌ 格式化消息失败: %v\n", err)
		fmt.Printf("   原始消息: %s\n", string(msg.Value))
	} else {
		fmt.Printf("   消息内容:\n%s\n", string(messageJSON))
	}

	fmt.Println("-" + string(make([]byte, 49)))
	return true
}
