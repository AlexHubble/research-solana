package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
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
		count      = flag.Int("count", 10, "发送消息数量")
		interval   = flag.Duration("interval", time.Second, "发送间隔")
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

	// 创建Producer
	fmt.Println("🚀 创建Kafka Producer...")
	producer, err := kafka.NewProducer(cfg.GetProducerConfigMap())
	if err != nil {
		log.Fatalf("❌ 创建Producer失败: %v", err)
	}
	defer producer.Close()

	// 启动事件处理goroutine
	go handleEvents(producer)

	topic := cfg.Producer.Topic
	fmt.Printf("📤 开始发送消息到topic: %s\n", topic)

	// 发送消息
	for i := 0; i < *count; i++ {
		message := Message{
			ID:        i + 1,
			Timestamp: time.Now().Format(time.RFC3339),
			Message:   fmt.Sprintf("这是第 %d 条测试消息", i+1),
			Producer:  "confluent-demo-producer-go",
			Data: map[string]interface{}{
				"sequence":   i + 1,
				"batch_size": *count,
			},
		}

		// 序列化消息
		messageBytes, err := json.Marshal(message)
		if err != nil {
			log.Printf("❌ 序列化消息失败: %v", err)
			continue
		}

		// 发送消息
		fmt.Printf("📨 发送消息 %d/%d\n", i+1, *count)
		err = producer.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{
				Topic:     &topic,
				Partition: kafka.PartitionAny,
			},
			Value: messageBytes,
		}, nil)

		if err != nil {
			log.Printf("❌ 发送消息失败: %v", err)
			continue
		}

		// 等待间隔（最后一条消息不需要等待）
		if i < *count-1 {
			time.Sleep(*interval)
		}
	}

	// 等待所有消息发送完成
	fmt.Println("⏳ 等待所有消息发送完成...")
	producer.Flush(15 * 1000) // 15秒超时

	fmt.Println("✅ 所有消息发送完成!")
}

// handleEvents 处理Producer事件
func handleEvents(producer *kafka.Producer) {
	for e := range producer.Events() {
		switch ev := e.(type) {
		case *kafka.Message:
			if ev.TopicPartition.Error != nil {
				fmt.Printf("❌ 消息发送失败: %v\n", ev.TopicPartition.Error)
			} else {
				fmt.Printf("✅ 消息发送成功: topic=%s, partition=%d, offset=%v\n",
					*ev.TopicPartition.Topic,
					ev.TopicPartition.Partition,
					ev.TopicPartition.Offset)
			}
		case kafka.Error:
			fmt.Printf("❌ Kafka错误: %v\n", ev)
		default:
			// 忽略其他事件类型
		}
	}
}
