package mikrotik

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
)

type Service struct {
	client *Client
}

func NewService(client *Client) *Service {
	return &Service{client: client}
}

func (s *Service) GetSystemResource(ctx context.Context) (*SystemResource, error) {
	var res SystemResource
	if err := s.client.get(ctx, "/system/resource", &res); err != nil {
		return nil, fmt.Errorf("system/resource: %w", err)
	}
	return &res, nil
}

func (s *Service) GetInterfaces(ctx context.Context) ([]Interface, error) {
	var result []Interface
	if err := s.client.get(ctx, "/interface", &result); err != nil {
		return nil, fmt.Errorf("interface: %w", err)
	}
	return result, nil
}

func (s *Service) GetIPAddresses(ctx context.Context) ([]IPAddress, error) {
	var result []IPAddress
	if err := s.client.get(ctx, "/ip/address", &result); err != nil {
		return nil, fmt.Errorf("ip/address: %w", err)
	}
	return result, nil
}

func (s *Service) GetRoutes(ctx context.Context) ([]Route, error) {
	var result []Route
	if err := s.client.get(ctx, "/ip/route", &result); err != nil {
		return nil, fmt.Errorf("ip/route: %w", err)
	}
	return result, nil
}

func (s *Service) GetDHCPLeases(ctx context.Context) ([]DHCPLease, error) {
	var result []DHCPLease
	if err := s.client.get(ctx, "/ip/dhcp-server/lease", &result); err != nil {
		return nil, fmt.Errorf("dhcp-server/lease: %w", err)
	}
	return result, nil
}

type result[T any] struct {
	data T
	err  error
}

func (s *Service) CollectSnapshot(ctx context.Context) *Snapshot {
	snap := &Snapshot{
		Errors: make(map[string]string),
	}

	sysCh := make(chan result[*SystemResource], 1)
	ifCh := make(chan result[[]Interface], 1)
	addrCh := make(chan result[[]IPAddress], 1)
	routeCh := make(chan result[[]Route], 1)
	dhcpCh := make(chan result[[]DHCPLease], 1)

	var wg sync.WaitGroup
	wg.Add(5)

	go func() {
		defer wg.Done()
		data, err := s.GetSystemResource(ctx)
		sysCh <- result[*SystemResource]{data, err}
	}()

	go func() {
		defer wg.Done()
		data, err := s.GetInterfaces(ctx)
		ifCh <- result[[]Interface]{data, err}
	}()

	go func() {
		defer wg.Done()
		data, err := s.GetIPAddresses(ctx)
		addrCh <- result[[]IPAddress]{data, err}
	}()

	go func() {
		defer wg.Done()
		data, err := s.GetRoutes(ctx)
		routeCh <- result[[]Route]{data, err}
	}()

	go func() {
		defer wg.Done()
		data, err := s.GetDHCPLeases(ctx)
		dhcpCh <- result[[]DHCPLease]{data, err}
	}()

	go func() {
		wg.Wait()
		close(sysCh)
		close(ifCh)
		close(addrCh)
		close(routeCh)
		close(dhcpCh)
	}()

	if r := <-sysCh; r.err != nil {
		slog.Error("collection error system", "err", r.err)
		snap.Errors["system"] = r.err.Error()
	} else {
		snap.System = r.data
	}

	if r := <-ifCh; r.err != nil {
		slog.Error("collection error interfaces", "err", r.err)
		snap.Errors["interfaces"] = r.err.Error()
	} else {
		snap.Interfaces = r.data
	}

	if r := <-addrCh; r.err != nil {
		slog.Error("collection error addresses", "err", r.err)
		snap.Errors["addresses"] = r.err.Error()
	} else {
		snap.Addresses = r.data
	}

	if r := <-routeCh; r.err != nil {
		slog.Error("collection error routes", "err", r.err)
		snap.Errors["routes"] = r.err.Error()
	} else {
		snap.Routes = r.data
	}

	if r := <-dhcpCh; r.err != nil {
		slog.Error("collection error dhcp", "err", r.err)
		snap.Errors["dhcp_leases"] = r.err.Error()
	} else {
		snap.DHCPLeases = r.data
	}

	return snap
}
