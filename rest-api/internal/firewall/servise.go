package firewall

import (
	"context"
	"fmt"

	"github.com/dmytroyunyk/MikrotikApi/internal/mikrotik"
)

type Servise struct {
	client *mikrotik.Client
}

func NewServise(client *mikrotik.Client) *Servise {
	return &Servise{client: client}
}

func (s *Servise) GetFilterRules(ctx context.Context) ([]FilterRule, error) {
	var Result []FilterRule
	if err := s.client.Get(ctx, "/ip/firewall/filter", &Result); err != nil {
		return nil, fmt.Errorf("filter rules: %w", err)
	}
	return Result, nil
}

func (s *Servise) GetNATRules(ctx context.Context) ([]NATRule, error) {
	var result []NATRule
	if err := s.client.Get(ctx, "/ip/firewall/nat", &result); err != nil {
		return nil, fmt.Errorf("nat rules: %w", err)
	}
	return result, nil
}

func (s *Servise) GetAddressList(ctx context.Context) ([]AdressListEntry, error) {
	var result []AdressListEntry
	if err := s.client.Get(ctx, "/ip/firewall/adress-list", &result); err != nil {
		return nil, fmt.Errorf("adress list: %w", err)
	}
	return result, nil
}
