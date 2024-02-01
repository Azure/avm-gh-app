package pkg

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/go-github/v58/github"
	"github.com/palantir/go-githubapp/githubapp"
	"github.com/pkg/errors"
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

	fmt.Printf("After: %s\n", *event.After)
	fmt.Printf("Repo Full Name: %s\n", *event.Repo.FullName)

	installationID := githubapp.GetInstallationIDFromEvent(&event)
	client, err := p.NewInstallationClient(installationID)
	if err != nil {
		return err
	}
	token, _, err := client.Apps.CreateInstallationToken(ctx, installationID, nil)
	if err != nil {
		return err
	}
	fmt.Printf("installation token: %+v", token)
	return nil
}
