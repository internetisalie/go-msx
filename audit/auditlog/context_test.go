// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package auditlog

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

var mockRequestDetails = &RequestDetails{
	Source:   "10.10.10.10",
	Protocol: "https",
	Host:     "192.168.2.1",
	Port:     "8080",
}

func TestContextWithRequestDetails(t *testing.T) {
	ctx := ContextWithRequestDetails(context.Background(), mockRequestDetails)
	assert.Equal(t, ctx.Value(contextKeyRequestAudit), mockRequestDetails)
}

func TestRequestAuditFromContext(t *testing.T) {
	ctx := ContextWithRequestDetails(context.Background(), mockRequestDetails)
	assert.Equal(t, mockRequestDetails, RequestAuditFromContext(ctx))
}
