package nsxbot

import (
	"context"
	"errors"
	"math"

	"github.com/atopos31/nsxbot/driver"
	"github.com/atopos31/nsxbot/types"
)

const abortIndex int8 = math.MaxInt8 >> 1

type Context[T any] struct {
	context.Context
	driver.Emitter

	replayer *types.Replyer
	Time     int64
	SelfId   int64
	index    int8
	Msg      T
	handlers HandlersChain[T]
}

func NewContext[T any](ctx context.Context, emitter driver.Emitter, selfId int64, time int64, data T, replayer *types.Replyer) Context[T] {
	return Context[T]{
		Context:  ctx,
		Emitter:  emitter,
		SelfId:   selfId,
		Time:     time,
		Msg:      data,
		replayer: replayer,
		index:    -1,
	}
}

var (
	ErrNoAvailable = errors.New("no replayer available")
)

func (c *Context[T]) Reply(text string) error {
	if c.replayer != nil {
		return c.replayer.Reply(text)
	}
	return ErrNoAvailable
}

func (c *Context[T]) Next() {
	c.index++
	for c.index < int8(len(c.handlers)) {
		if c.handlers[c.index] != nil {
			c.handlers[c.index](c)
		}
		c.index++
	}
}

func (c *Context[T]) Abort() {
	c.index = abortIndex
}
