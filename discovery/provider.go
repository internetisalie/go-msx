// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

//go:generate mockery --inpackage --name=DiscoveryProvider --structname=MockDiscoveryProvider --case underscore
//go:generate mockery --inpackage --name=RegistrationProvider --structname=MockRegistrationProvider --case underscore
package discovery

import (
	"context"
	"github.com/pkg/errors"
)

type RegistrationProvider interface {
	Register(ctx context.Context) error
	Deregister(ctx context.Context) error
}

type DiscoveryProvider interface {
	Discover(ctx context.Context, name string, healthyOnly bool, tags ...string) (ServiceInstances, error)
	DiscoverAll(ctx context.Context, healthyOnly bool, tags ...string) (ServiceInstances, error)
}

var (
	registrationProvider RegistrationProvider
	discoveryProvider    DiscoveryProvider

	ErrRegistrationProviderNotDefined = errors.New("Registration provider not registered")
	ErrDiscoveryProviderNotDefined    = errors.New("Discovery provider not registered")
)

func RegisterDiscoveryProvider(provider DiscoveryProvider) {
	if provider != nil {
		discoveryProvider = provider
	}
}

func IsDiscoveryProviderRegistered() bool {
	return discoveryProvider != nil
}

func Discover(ctx context.Context, name string, healthyOnly bool, tags ...string) (ServiceInstances, error) {
	if discoveryProvider == nil {
		return nil, ErrDiscoveryProviderNotDefined
	}

	return discoveryProvider.Discover(ctx, name, healthyOnly, tags...)
}

func DiscoverAll(ctx context.Context, healthyOnly bool, tags ...string) (ServiceInstances, error) {
	if discoveryProvider == nil {
		return nil, ErrDiscoveryProviderNotDefined
	}

	return discoveryProvider.DiscoverAll(ctx, healthyOnly, tags...)
}

func RegisterRegistrationProvider(provider RegistrationProvider) {
	if provider != nil {
		registrationProvider = provider
	}
}

func Register(ctx context.Context) error {
	if registrationProvider == nil {
		return ErrRegistrationProviderNotDefined
	}
	return registrationProvider.Register(ctx)
}

func Deregister(ctx context.Context) error {
	if registrationProvider == nil {
		return ErrRegistrationProviderNotDefined
	}
	return registrationProvider.Deregister(ctx)
}
