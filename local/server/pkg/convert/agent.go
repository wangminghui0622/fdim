package convert

import (
	"time"

	"fdim/pkg/db/table/bot"
	"fdim/pkg/util/datautil"
	pbbot "fdim/protocol/bot"
)

func DB2PBAgent(a *bot.Agent) *pbbot.Agent {
	return &pbbot.Agent{
		UserID:     a.UserID,
		Nickname:   a.NickName,
		FaceURL:    a.FaceURL,
		Url:        a.Url,
		Key:        a.Key,
		Identity:   a.Identity,
		Model:      a.Model,
		Prompts:    a.Prompts,
		CreateTime: a.CreateTime.UnixMilli(),
	}
}

func PB2DBAgent(a *pbbot.Agent) *bot.Agent {
	return &bot.Agent{
		UserID:     a.UserID,
		NickName:   a.Nickname,
		FaceURL:    a.FaceURL,
		Key:        a.Key,
		Url:        a.Url,
		Identity:   a.Identity,
		Model:      a.Model,
		Prompts:    a.Prompts,
		CreateTime: time.UnixMilli(a.CreateTime),
	}
}

func BatchDB2PBAgent(a []*bot.Agent) []*pbbot.Agent {
	return datautil.Batch(DB2PBAgent, a)
}
