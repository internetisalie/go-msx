package notification

import (
	"context"
	"errors"
	"time"

	"cto-github.cisco.com/NFV-BU/go-msx/security"
	"cto-github.cisco.com/NFV-BU/go-msx/stream/topics"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
)

type Message struct {
	Timestamp topics.Time            `json:"timestamp"`
	Version   *string                `json:"version"`
	Context   Context                `json:"context"`
	EventName string                 `json:"event"`
	Payload   map[string]interface{} `json:"payload"`
}

type Context struct {
	User            Identifier   `json:"user"`
	Provider        Identifier   `json:"provider"`
	Tenant          Identifier   `json:"tenant"`
	RecipientEmails []string     `json:"recipientEmails"`
	EmailLocale     string       `json:"emailLocale"`
	Recipients      []Identifier `json:"recipients"`
	ServiceType     ServiceType  `json:"serviceType"`
	Initiator       Identifier   `json:"initiator"`
}

type Identifier struct {
	Id   types.UUID `json:"id"`
	Name string     `json:"name"`
}

type ServiceType struct {
	LogicName              string `json:"logicName"`
	DisplayNameResourceKey string `json:"diplayNameResourceKey"`
}

type MessageProducer interface {
	Message(context.Context) (Message, error)
}

func NewMessage(ctx context.Context) (Message, error) {
	userContextDetails, err := security.NewUserContextDetails(ctx)
	if err != nil {
		return Message{}, err
	} else if userContextDetails.ProviderId == nil {
		return Message{}, errors.New("provider id not identified")
	}

	return Message{
		Timestamp: topics.Time(time.Now().UTC()),
		Context: Context{
			User: Identifier{
				Id:   userContextDetails.UserId,
				Name: *userContextDetails.Username,
			},
			Provider: Identifier{
				Id:   userContextDetails.ProviderId,
				Name: *userContextDetails.ProviderName,
			},
			Tenant: Identifier{
				Id:   userContextDetails.TenantId,
				Name: *userContextDetails.TenantName,
			},
		},
		Payload: make(map[string]interface{}),
	}, nil
}
