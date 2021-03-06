// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package auditing

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/streamtest"
	"github.com/pkg/errors"
	"testing"
)

func TestNewMessageListener(t *testing.T) {
	type args struct {
		fn      MessageHandler
		filters []MessageFilter
	}
	tests := []struct {
		name string
		args args
		test *streamtest.TopicReceiveTest
	}{
		{
			name: "NoFilter",
			args: args{
				fn: func(ctx context.Context, message Message) error {
					streamtest.TopicReceiveTestFromContext(ctx).Received()
					return nil
				},
				filters: []MessageFilter{},
			},
			test: streamtest.NewTopicReceiveTest().
				WithWantReceive(true),
		},
		{
			name: "MatchingFilter",
			args: args{
				fn: func(ctx context.Context, message Message) error {
					streamtest.TopicReceiveTestFromContext(ctx).Received()
					return nil
				},
				filters: []MessageFilter{
					FilterByAction("message-action"),
				},
			},
			test: streamtest.NewTopicReceiveTest().
				WithPayload([]byte(`{"action":"message-action"}`)).
				WithWantReceive(true),
		},
		{
			name: "NonMatchingFilter",
			args: args{
				fn: func(ctx context.Context, message Message) error {
					streamtest.TopicReceiveTestFromContext(ctx).Received()
					return nil
				},
				filters: []MessageFilter{
					FilterByAction("message-action"),
					FilterByService("some-service"),
				},
			},
			test: streamtest.NewTopicReceiveTest().
				WithPayload([]byte(`{"action":"message-action","service":"other-service"}`)).
				WithWantReceive(false),
		},
		{
			name: "PayloadError",
			args: args{
				fn: func(ctx context.Context, message Message) error {
					streamtest.TopicReceiveTestFromContext(ctx).Received()
					return nil
				},
			},
			test: streamtest.NewTopicReceiveTest().
				WithPayload([]byte("[")).
				WithWantReceive(false).
				WithWantError(true),
		},
		{
			name: "ListenerError",
			args: args{
				fn: func(ctx context.Context, message Message) error {
					streamtest.TopicReceiveTestFromContext(ctx).Received()
					return errors.New("listener-error")
				},
			},
			test: streamtest.NewTopicReceiveTest().
				WithWantReceive(true).
				WithWantError(true),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, tt.test.
			WithTopic(TopicName).
			WithAction(NewMessageListener(tt.args.fn, tt.args.filters)).
			Test)
	}
}
