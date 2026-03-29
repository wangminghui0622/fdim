package kafka

import (
	"bytes"
	"strings"

	"fdim/pkg/errs"
	"github.com/IBM/sarama"
)

func BuildConsumerGroupConfig(conf *Config, initial int64, autoCommitEnable bool) (*sarama.Config, error) {
	kfk := sarama.NewConfig()
	kfk.Version = sarama.V2_0_0_0
	kfk.Consumer.Offsets.Initial = initial
	kfk.Consumer.Offsets.AutoCommit.Enable = autoCommitEnable
	kfk.Consumer.Return.Errors = false
	if conf.ConsumerFetchDefaultBytes > 0 {
		kfk.Consumer.Fetch.Default = int32(conf.ConsumerFetchDefaultBytes)
	}
	if conf.ConsumerFetchMaxBytes > 0 {
		kfk.Consumer.Fetch.Max = int32(conf.ConsumerFetchMaxBytes)
	}
	if conf.Username != "" || conf.Password != "" {
		kfk.Net.SASL.Enable = true
		kfk.Net.SASL.User = conf.Username
		kfk.Net.SASL.Password = conf.Password
	}
	if conf.TLS.EnableTLS {
		tls, err := newTLSConfig(conf.TLS.ClientCrt, conf.TLS.ClientKey, conf.TLS.CACrt, []byte(conf.TLS.ClientKeyPwd), conf.TLS.InsecureSkipVerify)
		if err != nil {
			return nil, err
		}
		kfk.Net.TLS.Config = tls
		kfk.Net.TLS.Enable = true
	}
	return kfk, nil
}

func NewConsumerGroup(conf *sarama.Config, addr []string, groupID string) (sarama.ConsumerGroup, error) {
	cg, err := sarama.NewConsumerGroup(addr, groupID, conf)
	if err != nil {
		return nil, errs.WrapMsg(err, "NewConsumerGroup failed, addr=%v, groupID=%s", addr, groupID)
	}
	return cg, nil
}

func BuildProducerConfig(conf Config) (*sarama.Config, error) {
	kfk := sarama.NewConfig()
	kfk.Producer.Return.Successes = true
	kfk.Producer.Return.Errors = true
	kfk.Producer.Partitioner = sarama.NewHashPartitioner
	if conf.Username != "" || conf.Password != "" {
		kfk.Net.SASL.Enable = true
		kfk.Net.SASL.User = conf.Username
		kfk.Net.SASL.Password = conf.Password
	}
	switch strings.ToLower(conf.ProducerAck) {
	case "no_response":
		kfk.Producer.RequiredAcks = sarama.NoResponse
	case "wait_for_local":
		kfk.Producer.RequiredAcks = sarama.WaitForLocal
	case "wait_for_all":
		kfk.Producer.RequiredAcks = sarama.WaitForAll
	default:
		kfk.Producer.RequiredAcks = sarama.WaitForAll
	}
	if conf.CompressType == "" {
		kfk.Producer.Compression = sarama.CompressionNone
	} else {
		if err := kfk.Producer.Compression.UnmarshalText(bytes.ToLower([]byte(conf.CompressType))); err != nil {
			return nil, errs.WrapMsg(err, "UnmarshalText failed, compressType=%s", conf.CompressType)
		}
	}
	if conf.MaxMessageBytes > 0 {
		kfk.Producer.MaxMessageBytes = conf.MaxMessageBytes
	}
	if conf.TLS.EnableTLS {
		tls, err := newTLSConfig(conf.TLS.ClientCrt, conf.TLS.ClientKey, conf.TLS.CACrt, []byte(conf.TLS.ClientKeyPwd), conf.TLS.InsecureSkipVerify)
		if err != nil {
			return nil, err
		}
		kfk.Net.TLS.Config = tls
		kfk.Net.TLS.Enable = true
	}
	return kfk, nil
}

func NewProducer(conf *sarama.Config, addr []string) (sarama.SyncProducer, error) {
	producer, err := sarama.NewSyncProducer(addr, conf)
	if err != nil {
		return nil, errs.WrapMsg(err, "NewSyncProducer failed, addr=%v", addr)
	}
	return producer, nil
}

type TLSConfig struct {
	EnableTLS          bool   `yaml:"enableTLS"`
	CACrt              string `yaml:"caCrt"`
	ClientCrt          string `yaml:"clientCrt"`
	ClientKey          string `yaml:"clientKey"`
	ClientKeyPwd       string `yaml:"clientKeyPwd"`
	InsecureSkipVerify bool   `yaml:"insecureSkipVerify"`
}

type Config struct {
	Username                  string    `yaml:"username"`
	Password                  string    `yaml:"password"`
	ProducerAck               string    `yaml:"producerAck"`
	CompressType              string    `yaml:"compressType"`
	MaxMessageBytes           int       `yaml:"maxMessageBytes"`
	ConsumerFetchDefaultBytes int       `yaml:"consumerFetchDefaultBytes"`
	ConsumerFetchMaxBytes     int       `yaml:"consumerFetchMaxBytes"`
	Addr                      []string  `yaml:"addr"`
	TLS                       TLSConfig `yaml:"tls"`
}

// DefaultConfig returns a default Kafka configuration
func DefaultConfig() *Config {
	return &Config{
		ProducerAck: "wait_for_all",
	}
}
