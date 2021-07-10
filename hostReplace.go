package httpclient

import (
	"context"
)

type (
	hostReplaceType string
)

const (
	hostReplaceKey hostReplaceType = "ctx"
)

// HostReplaceCtx replace request host
func HostReplaceCtx(ctx context.Context, host string) context.Context {
	return context.WithValue(ctx, hostReplaceKey, host)
}

func getHostReplaceCtx(ctx context.Context) string {
	var v = ctx.Value(hostReplaceKey)
	if v == nil {
		return ""
	}

	var host, ok = v.(string)
	if !ok {
		return ""
	}

	return host
}
