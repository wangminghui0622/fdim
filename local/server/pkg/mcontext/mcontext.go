// Copyright © 2023 OpenIM. All rights reserved.
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

package mcontext

import (
	"context"
	"crypto/rand"
	"encoding/hex"

	"fdim/pkg/errs"
	"fdim/protocol/constant"
)

var mapper = []string{constant.OperationID, constant.OpUserID, constant.OpUserPlatform, constant.ConnID}

func WithOpUserIDContext(ctx context.Context, opUserID string) context.Context {
	return context.WithValue(ctx, constant.OpUserID, opUserID)
}

func WithOpUserPlatformContext(ctx context.Context, platform string) context.Context {
	return context.WithValue(ctx, constant.OpUserPlatform, platform)
}

func WithTriggerIDContext(ctx context.Context, triggerID string) context.Context {
	return context.WithValue(ctx, constant.TriggerID, triggerID)
}

func NewCtx(operationID string) context.Context {
	c := context.Background()
	ctx := context.WithValue(c, constant.OperationID, operationID)
	return SetOperationID(ctx, operationID)
}

func SetOperationID(ctx context.Context, operationID string) context.Context {
	return context.WithValue(ctx, constant.OperationID, operationID)
}

func SetOpUserID(ctx context.Context, opUserID string) context.Context {
	return context.WithValue(ctx, constant.OpUserID, opUserID)
}

func SetConnID(ctx context.Context, connID string) context.Context {
	return context.WithValue(ctx, constant.ConnID, connID)
}

func GetOperationID(ctx context.Context) string {
	s, _ := ctx.Value(constant.OperationID).(string)
	return s
}

func GetOpUserID(ctx context.Context) string {
	s, _ := ctx.Value(constant.OpUserID).(string)
	return s
}

func GetConnID(ctx context.Context) string {
	s, _ := ctx.Value(constant.ConnID).(string)
	return s
}

func GetTriggerID(ctx context.Context) string {
	s, _ := ctx.Value(constant.TriggerID).(string)
	return s
}

func GetOpUserPlatform(ctx context.Context) string {
	s, _ := ctx.Value(constant.OpUserPlatform).(string)
	return s
}

func GetRemoteAddr(ctx context.Context) string {
	s, _ := ctx.Value(constant.RemoteAddr).(string)
	return s
}

func GetMustCtxInfo(ctx context.Context) (operationID, opUserID, platform, connID string, err error) {
	operationID, ok := ctx.Value(constant.OperationID).(string)
	if !ok {
		err = errs.ErrArgs.WrapMsg("ctx missing operationID")
		return
	}
	opUserID, ok1 := ctx.Value(constant.OpUserID).(string)
	if !ok1 {
		err = errs.ErrArgs.WrapMsg("ctx missing opUserID")
		return
	}
	platform, ok2 := ctx.Value(constant.OpUserPlatform).(string)
	if !ok2 {
		err = errs.ErrArgs.WrapMsg("ctx missing platform")
		return
	}
	connID, _ = ctx.Value(constant.ConnID).(string)
	return
}

func GetCtxInfos(ctx context.Context) (operationID, opUserID, platform, connID string, err error) {
	operationID, ok := ctx.Value(constant.OperationID).(string)
	if !ok {
		err = errs.ErrArgs.WrapMsg("ctx missing operationID")
		return
	}
	opUserID, _ = ctx.Value(constant.OpUserID).(string)
	platform, _ = ctx.Value(constant.OpUserPlatform).(string)
	connID, _ = ctx.Value(constant.ConnID).(string)
	return
}

func WithMustInfoCtx(values []string) context.Context {
	ctx := context.Background()
	for i, v := range values {
		ctx = context.WithValue(ctx, mapper[i], v)
	}
	return ctx
}

// GenGroupID generates a unique group ID with existence check callback
func GenGroupID(ctx context.Context, existsCheck func(id string) (bool, error)) (string, error) {
	for i := 0; i < 10; i++ {
		b := make([]byte, 16)
		rand.Read(b)
		id := hex.EncodeToString(b)
		
		if existsCheck != nil {
			exists, err := existsCheck(id)
			if err != nil {
				return "", err
			}
			if exists {
				continue // ID exists, try again
			}
		}
		return id, nil
	}
	return "", errs.ErrInternalServer.WrapMsg("failed to generate unique group ID after 10 attempts")
}
