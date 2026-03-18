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
	"context"
	"strings"

	"fdim/pkg/errs"
	"github.com/nats-io/nats.go"
)

func CheckSubjects(ctx context.Context, conf *Config, subjects []string) error {
	opts, err := BuildNatsOptions(conf)
	if err != nil {
		return err
	}

	var conn *nats.Conn
	if len(conf.Addr) > 0 {
		url := conf.Addr[0]
		if strings.HasPrefix(url, "nats://") {
			url = strings.TrimPrefix(url, "nats://")
		}
		conn, err = nats.Connect(url, opts...)
	} else {
		conn, err = nats.Connect(nats.DefaultURL, opts...)
	}
	if err != nil {
		return errs.WrapMsg(err, "Connect failed, config=%+v", conf)
	}
	defer conn.Close()

	// For NATS, subjects don't need to exist beforehand (they're created on publish/subscribe)
	// But we can verify the connection is working
	if !conn.IsConnected() {
		return errs.New("not connected to NATS server").Wrap()
	}

	return nil
}

func CheckHealth(ctx context.Context, conf *Config) error {
	opts, err := BuildNatsOptions(conf)
	if err != nil {
		return err
	}

	var conn *nats.Conn
	if len(conf.Addr) > 0 {
		url := conf.Addr[0]
		if strings.HasPrefix(url, "nats://") {
			url = strings.TrimPrefix(url, "nats://")
		}
		conn, err = nats.Connect(url, opts...)
	} else {
		conn, err = nats.Connect(nats.DefaultURL, opts...)
	}
	if err != nil {
		return errs.WrapMsg(err, "Connect failed, config=%+v", conf)
	}
	defer conn.Close()

	if !conn.IsConnected() {
		return errs.New("not connected to NATS server").Wrap()
	}

	// Check server status
	status := conn.Status()
	if status != nats.CONNECTED {
		return errs.New("NATS connection status is not CONNECTED").Wrap()
	}

	return nil
}
