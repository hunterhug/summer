package summer

import (
	"context"
	"google.golang.org/grpc/metadata"
)

func GetIncomingContext(ctx context.Context, name string) []string {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		return md.Get(name)
	}

	return nil
}

func GetIncomingContextOne(ctx context.Context, name string) string {
	data := GetIncomingContext(ctx, name)
	if len(data) > 0 {
		return data[0]
	}
	return ""
}

func GetOutgoingContext(ctx context.Context, name string) []string {
	md, ok := metadata.FromOutgoingContext(ctx)
	if ok {
		return md.Get(name)
	}

	return nil
}

func GetOutgoingContextOne(ctx context.Context, name string) string {
	data := GetOutgoingContext(ctx, name)
	if len(data) > 0 {
		return data[0]
	}
	return ""
}

func SetOutgoingContext(ctx context.Context, kv []string) context.Context {
	ctx = metadata.AppendToOutgoingContext(ctx, kv...)
	return ctx
}
