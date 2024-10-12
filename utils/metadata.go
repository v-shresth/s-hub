package utils

import (
	"context"
	"google.golang.org/grpc/metadata"
	"strconv"
)

// Metadata struct for holding metadata for requests
type Metadata struct {
	AuthedUserId    uint
	AuthedSessionId uint
}

const (
	AuthedUserId    = "authUserId"
	AuthedSessionId = "authedSessionId"
)

func ExtractMetadata(ctx context.Context) *Metadata {
	m := &Metadata{}

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if userAgents := md.Get(AuthedUserId); len(userAgents) > 0 {
			userId, err := strconv.Atoi(userAgents[0])
			if err == nil {
				m.AuthedUserId = uint(userId)
			}
		}

		if sessionIds := md.Get(AuthedSessionId); len(sessionIds) > 0 {
			sessionId, err := strconv.Atoi(sessionIds[0])
			if err == nil {
				m.AuthedSessionId = uint(sessionId)
			}
		}
	}

	return m
}
