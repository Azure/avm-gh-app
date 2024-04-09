package pkg

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/alexellis/go-execute/v2"
	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/google/go-github/v60/github"
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

	fullName := event.Repo.GetFullName()
	if !strings.Contains(fullName, "terraform-azurerm-avm-") {
		logger.Debug().Msg("non-avm repo, ignore")
		return nil
	}

	if !strings.HasSuffix(event.GetRef(), fmt.Sprintf("/%s", event.Repo.GetDefaultBranch())) {
		logger.Debug().Msg("push to non-default branch, ignore")
		return nil
	}

	if event.Pusher.GetName() == Cfg.ExpectedPusherName {
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
	return postPush(ctx, event, token)
}

func postPush(ctx context.Context, event github.PushEvent, token string) error {
	fullName := event.Repo.GetFullName()
	tmp, err := os.MkdirTemp("", fmt.Sprintf("%s*", strings.ReplaceAll(fullName, "/", "-")))
	if err != nil {
		return fmt.Errorf("cannot create temp folder")
	}
	defer func() {
		_ = os.RemoveAll(tmp)
	}()
	owner := strings.TrimSuffix(fullName, fmt.Sprintf("/%s", event.Repo.GetName()))
	scriptFolder := "scripts"
	if strings.Contains(fullName, "-avm-") {
		scriptFolder = "avm_scripts"
	}
	task := execute.ExecTask{
		Command: "curl -H 'Cache-Control: no-cache, no-store' -sSL \"https://raw.githubusercontent.com/Azure/tfmod-scaffold/main/scripts/post-push-starter.sh\" | bash",
		Shell:   true,
		Env: []string{
			fmt.Sprintf("GITHUB_REPOSITORY=%s", fullName),
			fmt.Sprintf("GITHUB_REPOSITORY_OWNER=%s", owner),
			fmt.Sprintf("GITHUB_TOKEN=%s", token),
			fmt.Sprintf("SCRIPT_FOLDER=%s", scriptFolder),
		},
		Cwd:         tmp,
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
