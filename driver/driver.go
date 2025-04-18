package driver

import (
	"context"
	"fmt"

	"github.com/atopos31/nsxbot/types"
	"github.com/tidwall/gjson"
)

type Driver interface {
	EmitterMux
	Listener
}

type Listener interface {
	Listen(ctx context.Context, eventChan chan<- types.Event) error
}

type EmitterMux interface {
	GetEmitter(selfId int64) (Emitter, error)
	AddEmitter(selfId int64, emitter Emitter)
}

type Emitter interface {
	SendPvtMsg(ctx context.Context, userId int64, msg types.MeaasgeChain) (*types.SendMsgRes, error)
	SendGrMsg(ctx context.Context, groupId int64, msg types.MeaasgeChain) (*types.SendMsgRes, error)
	GetMsg(ctx context.Context, msgId int) (*types.GetMsgRes, error)
	DelMsg(ctx context.Context, messageId int) error
	GetLoginInfo(ctx context.Context) (*types.LoginInfo, error)
	GetStrangerInfo(ctx context.Context, userId int64, noCache bool) (*types.StrangerInfo, error)
	GetStatus(ctx context.Context) (*types.Status, error)
	GetSelfId(ctx context.Context) (int64, error)
	Raw(ctx context.Context, action Action, params any) ([]byte, error)
}

type Request[T any] struct {
	Echo   string `json:"echo"`
	Action Action `json:"action"`
	Params T      `json:"params,omitempty"`
}

type Response[T any] struct {
	Status  string `json:"status"`
	RetCode int    `json:"retCode"`
	Data    T      `json:"data,omitempty"`
	Echo    string `json:"echo"`
}

func contentToEvent(content []byte) (types.Event, error) {
	strContent := string(content)
	postType := gjson.Get(strContent, "post_type")
	if !postType.Exists() {
		return types.Event{}, fmt.Errorf("invalid event, post_type: %v", postType.Exists())
	}

	Type := gjson.Get(strContent, postType.String()+"_type")
	if !Type.Exists() {
		return types.Event{}, fmt.Errorf("invalid event, %s_type: %v", postType.String(), Type.Exists())
	}

	time := gjson.Get(strContent, "time")
	selfID := gjson.Get(strContent, "self_id")
	if !time.Exists() || !selfID.Exists() {
		return types.Event{}, fmt.Errorf("invalid event, post_type: %v, time: %v, self_id: %v", postType.Exists(), time.Exists(), selfID.Exists())
	}

	return types.Event{
		Types:   []types.EventType{postType.String(), postType.String() + ":" + Type.String()},
		RawData: content,
		SelfID:  selfID.Int(),
		Time:    time.Int(),
	}, nil
}
