package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"gopkg.in/yaml.v3"
)

// KafkaConfig 包含所有Kafka相关配置
type KafkaConfig struct {
	Kafka    KafkaConnection `yaml:"kafka"`
	Producer ProducerConfig  `yaml:"producer"`
	Consumer ConsumerConfig  `yaml:"consumer"`
}

// KafkaConnection Kafka连接配置
type KafkaConnection struct {
	BootstrapServers string `yaml:"bootstrap_servers"`
	SecurityProtocol string `yaml:"security_protocol"`
	SASLMechanism    string `yaml:"sasl_mechanism"`
	SASLUsername     string `yaml:"sasl_username"`
	SASLPassword     string `yaml:"sasl_password"`
}

// ProducerConfig Producer配置
type ProducerConfig struct {
	Topic     string `yaml:"topic"`
	Acks      string `yaml:"acks"`
	Retries   int    `yaml:"retries"`
	BatchSize int    `yaml:"batch_size"`
	LingerMs  int    `yaml:"linger_ms"`
}

// ConsumerConfig Consumer配置
type ConsumerConfig struct {
	Topic                 string `yaml:"topic"`
	GroupID               string `yaml:"group_id"`
	AutoOffsetReset       string `yaml:"auto_offset_reset"`
	EnableAutoCommit      bool   `yaml:"enable_auto_commit"`
	AutoCommitIntervalMs  int    `yaml:"auto_commit_interval_ms"`
	SessionTimeoutMs      int    `yaml:"session_timeout_ms"`
}

// LoadConfig 从YAML文件加载配置
func LoadConfig(configPath string) (*KafkaConfig, error) {
	// 如果是相对路径，转换为绝对路径
	if !filepath.IsAbs(configPath) {
		wd, err := os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("获取工作目录失败: %w", err)
		}
		configPath = filepath.Join(wd, configPath)
	}

	// 读取配置文件
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败 %s: %w", configPath, err)
	}

	// 解析YAML
	var config KafkaConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	return &config, nil
}

// GetProducerConfigMap 获取Producer的配置映射
func (c *KafkaConfig) GetProducerConfigMap() *kafka.ConfigMap {
	return &kafka.ConfigMap{
		"bootstrap.servers":  c.Kafka.BootstrapServers,
		"security.protocol":  c.Kafka.SecurityProtocol,
		"sasl.mechanism":     c.Kafka.SASLMechanism,
		"sasl.username":      c.Kafka.SASLUsername,
		"sasl.password":      c.Kafka.SASLPassword,
		"acks":               c.Producer.Acks,
		"retries":            c.Producer.Retries,
		"batch.size":         c.Producer.BatchSize,
		"linger.ms":          c.Producer.LingerMs,
	}
}

// GetConsumerConfigMap 获取Consumer的配置映射
func (c *KafkaConfig) GetConsumerConfigMap() *kafka.ConfigMap {
	return &kafka.ConfigMap{
		"bootstrap.servers":        c.Kafka.BootstrapServers,
		"security.protocol":        c.Kafka.SecurityProtocol,
		"sasl.mechanism":           c.Kafka.SASLMechanism,
		"sasl.username":            c.Kafka.SASLUsername,
		"sasl.password":            c.Kafka.SASLPassword,
		"group.id":                 c.Consumer.GroupID,
		"auto.offset.reset":        c.Consumer.AutoOffsetReset,
		"enable.auto.commit":       c.Consumer.EnableAutoCommit,
		"auto.commit.interval.ms":  c.Consumer.AutoCommitIntervalMs,
		"session.timeout.ms":       c.Consumer.SessionTimeoutMs,
	}
}
