package pkg

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/alexellis/go-execute/v2"
	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/google/go-github/v58/github"
	"github.com/palantir/go-githubapp/githubapp"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
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

	if event.Pusher != nil && event.Pusher.GetName() == Cfg.ExpectedPusherName {
		logger.Debug().Msg("push event triggered by this app, ignore")
		return nil
	}

	installationID := githubapp.GetInstallationIDFromEvent(&event)
	itr, err := ghinstallation.New(http.DefaultTransport, Cfg.Github.App.IntegrationID, installationID, []byte(Cfg.Github.App.PrivateKey))
	if err != nil {
		return err
	}

	token, err := itr.Token(context.Background())
	if err != nil {
		return err
	}
	logger.Debug().Msg(fmt.Sprintf("token generated for event %+v, %s", event, token))
	//return postPush(ctx, event, token)
	return nil
}

func postPush(ctx context.Context, event github.PushEvent, token string) error {
	task := execute.ExecTask{
		Command: "curl -H 'Cache-Control: no-cache, no-store' -sSL \"https://raw.githubusercontent.com/Azure/tfmod-scaffold/main/scripts/post-push-starter.sh\" | sh -s",
		Shell:   true,
		Env: []string{
			fmt.Sprintf("GITHUB_REPOSITORY="),
			fmt.Sprintf("GITHUB_TOKEN=%s", token),
		},
		StreamStdio: true,
	}
	result, err := task.Execute(ctx)
	if err != nil {
		return err
	}
	if result.ExitCode != 0 {
		return fmt.Errorf("unexpected return code: %d", result.ExitCode)
	}
	return nil
}
