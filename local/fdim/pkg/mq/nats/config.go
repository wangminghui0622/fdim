// Copyright © 2024 OpenIM open source community. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package nats

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"os"
	"time"

	"fdim/pkg/errs"
	"github.com/nats-io/nats.go"
)

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
	Token                     string    `yaml:"token"`
	MaxReconnects             int       `yaml:"maxReconnects"`
	ReconnectWait             int       `yaml:"reconnectWait"` // seconds
	Timeout                   int       `yaml:"timeout"`        // seconds
	MaxMessageBytes           int       `yaml:"maxMessageBytes"`
	ConsumerFetchDefaultBytes int       `yaml:"consumerFetchDefaultBytes"`
	ConsumerFetchMaxBytes     int       `yaml:"consumerFetchMaxBytes"`
	Addr                      []string  `yaml:"addr"`
	TLS                       TLSConfig `yaml:"tls"`
}

// DefaultConfig returns a default NATS configuration
func DefaultConfig() *Config {
	return &Config{
		MaxReconnects: 60,
		ReconnectWait: 2,
		Timeout:       5,
	}
}

// BuildNatsOptions builds NATS connection options from config
func BuildNatsOptions(conf *Config) ([]nats.Option, error) {
	var opts []nats.Option

	// Set reconnection options
	if conf.MaxReconnects > 0 {
		opts = append(opts, nats.MaxReconnects(conf.MaxReconnects))
	}
	if conf.ReconnectWait > 0 {
		opts = append(opts, nats.ReconnectWait(time.Duration(conf.ReconnectWait)*time.Second))
	}
	if conf.Timeout > 0 {
		opts = append(opts, nats.Timeout(time.Duration(conf.Timeout)*time.Second))
	}

	// Set authentication
	if conf.Token != "" {
		opts = append(opts, nats.Token(conf.Token))
	} else if conf.Username != "" || conf.Password != "" {
		opts = append(opts, nats.UserInfo(conf.Username, conf.Password))
	}

	// Set TLS if enabled
	if conf.TLS.EnableTLS {
		tlsConfig, err := newTLSConfig(conf.TLS.ClientCrt, conf.TLS.ClientKey, conf.TLS.CACrt, []byte(conf.TLS.ClientKeyPwd), conf.TLS.InsecureSkipVerify)
		if err != nil {
			return nil, err
		}
		opts = append(opts, nats.Secure(tlsConfig))
	}

	return opts, nil
}

// decryptPEM decrypts a PEM block using a password.
func decryptPEM(data []byte, passphrase []byte) ([]byte, error) {
	if len(passphrase) == 0 {
		return data, nil
	}
	b, _ := pem.Decode(data)
	d, err := x509.DecryptPEMBlock(b, passphrase)
	if err != nil {
		return nil, errs.WrapMsg(err, "DecryptPEMBlock failed")
	}
	return pem.EncodeToMemory(&pem.Block{
		Type:  b.Type,
		Bytes: d,
	}), nil
}

func readEncryptablePEMBlock(path string, pwd []byte) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, errs.WrapMsg(err, "ReadFile failed, path=%s", path)
	}
	return decryptPEM(data, pwd)
}

// newTLSConfig setup the TLS config from general config file.
func newTLSConfig(clientCertFile, clientKeyFile, caCertFile string, keyPwd []byte, insecureSkipVerify bool) (*tls.Config, error) {
	var tlsConfig tls.Config
	if clientCertFile != "" && clientKeyFile != "" {
		certPEMBlock, err := os.ReadFile(clientCertFile)
		if err != nil {
			return nil, errs.WrapMsg(err, "ReadFile failed, clientCertFile=%s", clientCertFile)
		}
		keyPEMBlock, err := readEncryptablePEMBlock(clientKeyFile, keyPwd)
		if err != nil {
			return nil, err
		}

		cert, err := tls.X509KeyPair(certPEMBlock, keyPEMBlock)
		if err != nil {
			return nil, errs.WrapMsg(err, "X509KeyPair failed")
		}
		tlsConfig.Certificates = []tls.Certificate{cert}
	}

	if caCertFile != "" {
		caCert, err := os.ReadFile(caCertFile)
		if err != nil {
			return nil, errs.WrapMsg(err, "ReadFile failed, caCertFile=%s", caCertFile)
		}
		caCertPool := x509.NewCertPool()
		if ok := caCertPool.AppendCertsFromPEM(caCert); !ok {
			return nil, errs.New("AppendCertsFromPEM failed")
		}
		tlsConfig.RootCAs = caCertPool
	}
	tlsConfig.InsecureSkipVerify = insecureSkipVerify
	return &tlsConfig, nil
}
