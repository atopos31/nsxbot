package main

import (
	"context"
	"log/slog"
	"os"
	"strconv"

	"github.com/nsxdevx/nsxbot"
	"github.com/nsxdevx/nsxbot/driver"
	"github.com/nsxdevx/nsxbot/filter"
	"github.com/nsxdevx/nsxbot/types"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	bot := nsxbot.Default(driver.NewDriverHttp(":8080", "http://localhost:4000"))
	groupId, _ := strconv.ParseInt(os.Getenv("TEST_GROUP"), 10, 64)
	gr := nsxbot.OnEvent[types.EventGrMsg](bot)
	gr.Handle(func(ctx *nsxbot.Context[types.EventGrMsg]) {
		text, err := ctx.Msg.TextFirst()
		if err != nil {
			slog.Error("Error parsing message", "error", err)
			return
		}
		slog.Info("Group Message", "message", text.Text)
	}, filter.OnlyGroups(groupId))

	// Run
	bot.Run(ctx)
}
