// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package streamops

import (
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"github.com/pkg/errors"
)

var logger = log.NewPackageLogger()

const PeerNameContentType = "contentType"
const PeerNameContentEncoding = "contentEncoding"

type contextKey string

var ErrNotImplemented = errors.New("Not implemented")
