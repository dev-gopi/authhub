package service

import (
	"context"
	"fmt"
	"time"

	"github.com/dev-gopi/authhub/internal/modules/rootauth/entity"
	"github.com/dev-gopi/authhub/internal/shared/security"
)

const (
	loginAttemptWindow = 24 * time.Hour
	maxLoginBlock      = 15 * time.Minute
)

func (s *Service) checkRateLimit(
	ctx context.Context,
	identifier string,
	ipAddress string,
) error {
	identifierHash := security.HashIdentifier(
		identifier,
	)

	identifierBlockKey := fmt.Sprintf(
		"authhub:rate:root-login:identifier:block:%s",
		identifierHash,
	)

	blocked, err := s.redis.Client.Exists(
		ctx,
		identifierBlockKey,
	).Result()
	if err != nil {
		return fmt.Errorf(
			"check identifier login limit: %w",
			err,
		)
	}

	if blocked > 0 {
		return entity.ErrRateLimited
	}

	if ipAddress != "" {
		ipHash := security.HashIdentifier(ipAddress)

		ipBlockKey := fmt.Sprintf(
			"authhub:rate:root-login:ip:block:%s",
			ipHash,
		)

		blocked, err = s.redis.Client.Exists(
			ctx,
			ipBlockKey,
		).Result()
		if err != nil {
			return fmt.Errorf(
				"check ip login limit: %w",
				err,
			)
		}

		if blocked > 0 {
			return entity.ErrRateLimited
		}
	}

	return nil
}

func (s *Service) recordRateLimitFailure(
	ctx context.Context,
	identifier string,
	ipAddress string,
) {
	s.recordFailureKey(
		ctx,
		"identifier",
		security.HashIdentifier(identifier),
	)

	if ipAddress != "" {
		s.recordFailureKey(
			ctx,
			"ip",
			security.HashIdentifier(ipAddress),
		)
	}
}

func (s *Service) recordFailureKey(
	ctx context.Context,
	kind string,
	valueHash string,
) {
	counterKey := fmt.Sprintf(
		"authhub:rate:root-login:%s:fail:%s",
		kind,
		valueHash,
	)

	count, err := s.redis.Client.Incr(
		ctx,
		counterKey,
	).Result()

	if err != nil {
		return
	}

	if count == 1 {
		_ = s.redis.Client.Expire(
			ctx,
			counterKey,
			loginAttemptWindow,
		).Err()
	}

	if count%5 != 0 {
		return
	}

	blockLevel := (count / 5) - 1

	blockDuration := 30 * time.Second

	for i := int64(0); i < blockLevel; i++ {
		blockDuration *= 2

		if blockDuration >= maxLoginBlock {
			blockDuration = maxLoginBlock
			break
		}
	}

	blockKey := fmt.Sprintf(
		"authhub:rate:root-login:%s:block:%s",
		kind,
		valueHash,
	)

	_ = s.redis.Client.Set(
		ctx,
		blockKey,
		"1",
		blockDuration,
	).Err()
}

func (s *Service) clearRateLimit(
	ctx context.Context,
	identifier string,
	ipAddress string,
) {
	keys := []string{
		fmt.Sprintf(
			"authhub:rate:root-login:identifier:fail:%s",
			security.HashIdentifier(identifier),
		),
	}

	if ipAddress != "" {
		keys = append(
			keys,
			fmt.Sprintf(
				"authhub:rate:root-login:ip:fail:%s",
				security.HashIdentifier(ipAddress),
			),
		)
	}

	_ = s.redis.Client.Del(
		ctx,
		keys...,
	).Err()
}
