package pkg

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/google/go-github/v58/github"
	"github.com/palantir/go-githubapp/githubapp"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"net/http"
)

var _ githubapp.EventHandler = &PushHandler{}

type PushHandler struct {
	githubapp.ClientCreator
}

func (p PushHandler) Handles() []string {
	return []string{"push"}
}

func (p PushHandler) Handle(ctx context.Context, eventType, deliveryID string, payload []byte) error {
	var event github.PushEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return errors.Wrap(err, "failed to parse issue comment event payload")
	}
	logger := zerolog.Ctx(ctx)

	if event.Pusher != nil && event.Pusher.GetName() == config.ExpectedPusherName {
		logger.Debug().Msg("push event triggered by this app, ignore")
		return nil
	}

	installationID := githubapp.GetInstallationIDFromEvent(&event)
	itr, err := ghinstallation.New(http.DefaultTransport, config.Github.App.IntegrationID, installationID, []byte(config.Github.App.PrivateKey))
	if err != nil {
		return err
	}

	token, err := itr.Token(context.Background())
	if err != nil {
		return err
	}
	logger.Debug().Msg(fmt.Sprintf("token generated for event %+v", event))
	return postPush(event, token)
}

func postPush(event github.PushEvent, token string) error {
	return nil
}
