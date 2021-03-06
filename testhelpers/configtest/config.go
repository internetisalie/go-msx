// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package configtest

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
)

func NewInMemoryConfig(values map[string]string) *config.Config {
	provider := config.NewInMemoryProvider("testdata", values)
	cfg := config.NewConfig(provider)
	_ = cfg.Load(context.Background())
	return cfg
}

func ContextWithNewInMemoryConfig(ctx context.Context, values map[string]string) context.Context {
	cfg := NewInMemoryConfig(values)
	return config.ContextWithConfig(ctx, cfg)
}
