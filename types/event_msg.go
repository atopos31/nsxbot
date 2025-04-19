package types

import (
	"encoding/json"
	"errors"
)

type EventType = string

const (
	POST_TYPE_MESSAGE    = "message"
	POST_TYPE_NOTICE     = "notice"
	POST_TYPE_REQUEST    = "request"
	POST_TYPE_META_ENEVT = "meta_event"
)

type Event struct {
	Types   []EventType
	Time    int64
	SelfID  int64
	RawData []byte
	Replyer Replayer
}

type Eventer interface {
	Type() EventType
}

type EventMsg interface {
	Eventer
	TextFirst() (*Text, error)
	Texts() ([]Text, int)
	FaceFirst() (*Face, error)
	Faces() ([]Face, int)
	AtFirst() (*At, error)
	Ats() ([]At, int)
	ImageFirst() (*Image, error)
	Images() ([]Image, int)
}

type BaseMessage struct {
	SubType    string    `json:"sub_type"`
	MessageId  int       `json:"message_id"`
	UserId     int64     `json:"user_id"`
	Messages   []Message `json:"message"`
	RawMessage string    `json:"raw_message"`
	Font       int       `json:"font"`
	Sender     Sender    `json:"sender"`
}

type Sender struct {
	UserID   int64  `json:"user_id"`
	Nickname string `json:"nickname"`
	Sex      string `json:"sex"`
	Age      int    `json:"age"`
}

type EventPvtMsg struct {
	*BaseMessage
}

func (e EventPvtMsg) Type() EventType {
	return "message:private"
}

func (em *BaseMessage) Reply(replyer Replayer, text string) error {
	if replyer == nil {
		return ErrNoAvailable
	}
	data := struct {
		Reply string `json:"reply"`
	}{Reply: text}
	return replyer.Reply(data)
}

func (em *BaseMessage) Id() int {
	return em.MessageId
}

func (em *BaseMessage) TextFirst() (*Text, error) {
	return first[Text]("text", em.Messages)
}

func (em *BaseMessage) Texts() ([]Text, int) {
	return all[Text]("text", em.Messages)
}

func (em *BaseMessage) Faces() ([]Face, int) {
	return all[Face]("face", em.Messages)
}

func (em *BaseMessage) FaceFirst() (*Face, error) {
	return first[Face]("face", em.Messages)
}

func (em *BaseMessage) AtFirst() (*At, error) {
	return first[At]("at", em.Messages)
}

func (em *BaseMessage) Ats() ([]At, int) {
	return all[At]("at", em.Messages)
}

func (em *BaseMessage) ReplyFirst() (*Reply, error) {
	return first[Reply]("reply", em.Messages)
}

func (em *BaseMessage) ImageFirst() (*Image, error) {
	return first[Image]("image", em.Messages)
}

func (em *BaseMessage) Images() ([]Image, int) {
	return all[Image]("image", em.Messages)
}

var (
	ErrNotFound = errors.New("not found")
)

func first[T any](msgType string, msg []Message) (*T, error) {
	for _, msg := range msg {
		if msg.Type == msgType {
			var data T
			if err := json.Unmarshal(msg.Data, &data); err != nil {
				return nil, err
			}
			return &data, nil
		}
	}
	return nil, ErrNotFound
}

func all[T any](msgType string, msg []Message) ([]T, int) {
	var data []T
	for _, msg := range msg {
		if msg.Type == msgType {
			var d T
			if err := json.Unmarshal(msg.Data, &data); err != nil {
				continue
			}
			data = append(data, d)
		}
	}
	return data, len(data)
}

type EventGrMsg struct {
	*BaseMessage
	GroupId   int64      `json:"group_id"`
	Anonymous *Anonymous `json:"anonymous"`
}

type Anonymous struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
	Flag string `json:"flag"`
}

func (e EventGrMsg) Type() EventType {
	return "message:group"
}

type EventAllMsg struct {
	*EventGrMsg
}

func (em EventAllMsg) Type() EventType {
	return "message"
}
