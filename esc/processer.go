package esc

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"

	"github.com/RainFallsSilent/activation-statistics/common"
)

func Process(ctx context.Context, days, startHour uint32) *common.Activation {
	service, err := New(days, int64(startHour))
	if err != nil {
		g.Log().Fatal(ctx, "create esc service failed", "error", err)
	}

	service.Start()

	return service.activation
}
