// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package kafka

import (
	"context"
	"github.com/Shopify/sarama"
	"github.com/pkg/errors"
)

func TopicMap(conn *Connection) (map[string]struct{}, error) {
	client := conn.Client()
	existingTopics, err := client.Topics()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to retrieve Kafka Topics")
	}

	topicMap := make(map[string]struct{})
	for _, topic := range existingTopics {
		topicMap[topic] = struct{}{}
	}

	return topicMap, nil
}

func CreateTopics(ctx context.Context, conn *Connection, topics ...string) (err error) {
	logger.WithContext(ctx).Info("Ensuring topics exist")

	topicMap, err := TopicMap(conn)
	if err != nil {
		return errors.Wrap(err, "Failed to create topic map")
	}

	detail := sarama.TopicDetail{
		NumPartitions:     int32(conn.cfg.DefaultPartitions),
		ReplicationFactor: int16(conn.cfg.ReplicationFactor),
	}

	saramaClusterAdmin, err := conn.ClusterAdmin()
	if err != nil {
		return errors.Wrap(err, "Failed to create cluster admin")
	}

	for _, topic := range topics {
		_, exists := topicMap[topic]
		if exists {
			continue
		}

		logger.Infof("Creating kafka topic: %s", topic)
		err = saramaClusterAdmin.CreateTopic(topic, &detail, false)
		if err != nil {
			// Ignore topic creation failure due to existing topic
			var topicError *sarama.TopicError
			ok := errors.As(err, &topicError)
			if ok && topicError.Err == sarama.ErrTopicAlreadyExists {
				return ErrTopicAlreadyExists
			}
		}

		if err != nil {
			return errors.Wrap(err, "Failed to create Kafka Topic")
		}
	}

	return nil
}
