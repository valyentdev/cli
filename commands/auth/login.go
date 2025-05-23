package auth

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	stdHTTP "net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"github.com/spf13/cobra"
	"github.com/valyentdev/cli/config"
	"github.com/valyentdev/cli/http"
	"github.com/valyentdev/cli/pkg/env"
	"github.com/valyentdev/cli/pkg/exit"
	"github.com/valyentdev/valyent.go"
)

func newLoginCmd() *cobra.Command {
	loginCmd := &cobra.Command{
		Use: "login",
		RunE: func(cmd *cobra.Command, args []string) error {
			apiKey, err := cmd.Flags().GetString("api-key")
			if err != nil {
				return err
			}
			organization, err := cmd.Flags().GetString("organization")
			if err != nil {
				return err
			}

			if err := runLoginCmd(apiKey, organization); err != nil {
				exit.WithError(err)
			}

			return nil
		},
	}
	loginCmd.Flags().String("api-key", "", "API key (optional)")
	loginCmd.Flags().String("organization", "", "Organization (optional)")

	return loginCmd
}

func runLoginCmd(apiKey, organization string) (err error) {
	if apiKey != "" {
		goto authenticate
	}

	{
		manual := false

		err = huh.NewConfirm().
			Title("Do you want to manually copy/paste an API Key created from Valyent's dashboard?").
			Negative("No.").
			Affirmative("Yes!").
			Value(&manual).
			Run()
		if err != nil {
			return
		}

		if manual {
			err = huh.NewInput().
				Title("Copy/paste your API key below").
				Value(&apiKey).
				Run()
			if err != nil {
				return
			}
		} else {
			organization, apiKey, err = retrieveAPIKeyFromTheBrowser()
			if err != nil {
				return
			}
		}
	}

authenticate:
	// Store the API Key on the user's machine.
	err = config.StoreAPIKey(organization, apiKey)
	if err != nil {
		return
	}

	// Authenticate Docker (allowing to interact directly with Valyent's registry).
	err = authenticateDockerRegistry(apiKey)
	if err != nil {
		exit.WithError(err)
	}

	fmt.Println("🎉 Successfully authenticated.")

	return
}

func retrieveAPIKeyFromTheBrowser() (namespace, apiKey string, err error) {
	err = spinner.New().
		Title("Waiting for authentication...").
		Action(func() {
			// Initialize new Valyent API HTTP client.
			var client *valyent.Client
			client, err = http.NewClient()
			if err != nil {
				err = fmt.Errorf("failed to initialize Valyent API HTTP client: %v", err)
				return
			}

			// Retrieve an authentication session.
			type response struct {
				SessionID string `json:"sessionID"`
			}
			res := &response{}

			err = client.PerformRequest(stdHTTP.MethodGet, "/auth/cli/session", nil, res)
			if err != nil {
				return
			}

			// Open the authentication page in the browser.
			baseURL := env.GetVar("VALYENT_API_URL", valyent.DEFAULT_BASE_URL)
			url := baseURL + "/auth/cli/" + res.SessionID
			err = openInBrowser(url)
			if err != nil {
				return
			}

			// Wait for the user to authenticate his session.
			namespace, apiKey, err = waitForLogin(res.SessionID)
			if err != nil {
				return
			}
		}).
		Run()

	return
}

func waitForLogin(sessionId string) (namespace, apiKey string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return "", "", fmt.Errorf("login timed out after 2 minutes")
		case <-ticker.C:
			// Initialize new Valyent API HTTP client.
			client, err := http.NewClient()
			if err != nil {
				return "", "", fmt.Errorf("failed to initialize Valyent API HTTP client: %v", err)
			}

			type waitForLoginResponse struct {
				Status    string `json:"status"`
				APIKey    string `json:"apiKey"`
				Namespace string `json:"namespace"`
			}
			res := &waitForLoginResponse{}

			path := "/auth/cli/" + sessionId + "/wait"
			err = client.PerformRequest(stdHTTP.MethodGet, path, nil, res)
			if err != nil {
				return "", "", fmt.Errorf("authentication check failed: %w", err)
			}

			if res.Status != "pending" {
				return res.Namespace, res.APIKey, nil
			}
		}
	}
}

func openInBrowser(url string) (err error) {
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	return
}

func authenticateDockerRegistry(apiKey string) (err error) {
	binary, err := exec.LookPath("docker")
	if err != nil {
		// If we can't find the Docker's binary, try configuring the JSON directly.
		if err := configureDockerJSON(apiKey); err == nil {
			return nil
		}
		return fmt.Errorf("docker cli not found - make sure it's installed and try again: %w", err)
	}

	// Compute Valyent's registry host.
	registryHost := env.GetVar("VALYENT_REGISTRY_HOST", "registry.valyent.cloud")

	// Prepare `docker login` command.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, binary, "login", "--username=x", "--password", apiKey, registryHost)

	// Start the command (without stopping it, yet).
	if err = cmd.Start(); err != nil {
		return
	}

	// Wait for the `docker login` command to be completed
	if err = cmd.Wait(); err != nil {
		return fmt.Errorf("failed authenticating with %s", registryHost)
	}

	return
}

// configureDockerJSON adds Valyent's registry auth stuff to Docker's config.json file.
func configureDockerJSON(apiKey string) error {
	if runtime.GOOS == "windows" {
		return errors.New("unsuppported")
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	if err := ensureDockerConfigDir(home); err != nil {
		return err
	}

	configPath := filepath.Join(home, ".docker", "config.json")
	configJSON, err := os.ReadFile(configPath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	updatedJSON, err := addValyentAuthToDockerConfig(apiKey, configJSON)
	if err != nil {
		return err
	}
	// It needs to be readable by Docker, if it gets installed in the future.
	return os.WriteFile(configPath, updatedJSON, 0o644)
}

// ensureDockerConfigDir checks to see if the "${HOME}/.docker" directory exists,
// it creates the dir if it doesn't.
func ensureDockerConfigDir(home string) error {
	dockerDir := filepath.Join(home, ".docker")
	fi, err := os.Stat(dockerDir)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		// It needs to be readable by Docker, if it gets installed in the
		// future.
		// The permission is 700 as like Docker itself.
		// https://github.com/docker/cli/blob/v23.0.5/cli/config/configfile/file.go#L142
		if err := os.Mkdir(dockerDir, 0o700); err != nil {
			return err
		}
	} else if !fi.IsDir() {
		return errors.New("~/.docker is not a dir")
	}
	return nil
}

// addValyentAuthToDockerConfig adds Valyent's registry to the provided JSON object
// and returns the updated JSON.
//
// The config.json is structured as follows:
//
//	{
//	  "auths": {
//	    "registry.valyent.cloud": {
//	      "auth": "x:..."
//	    }
//	  }
//	}
func addValyentAuthToDockerConfig(apiKey string, configJSON []byte) ([]byte, error) {
	// Compute Valyent's registry host.
	registryHost := env.GetVar("VALYENT_REGISTRY_HOST", "registry.valyent.cloud")

	var dockerConfig map[string]json.RawMessage
	if len(configJSON) == 0 {
		dockerConfig = make(map[string]json.RawMessage)
	} else if err := json.Unmarshal(configJSON, &dockerConfig); err != nil {
		return nil, err
	}

	var dockerAuthProviders map[string]json.RawMessage
	if a, ok := dockerConfig["auths"]; ok {
		if err := json.Unmarshal(a, &dockerAuthProviders); err != nil {
			return nil, err
		}
	} else {
		dockerAuthProviders = make(map[string]json.RawMessage)
	}

	var valyentAuth map[string]interface{}
	if a, ok := dockerAuthProviders[registryHost]; ok {
		if err := json.Unmarshal(a, &valyentAuth); err != nil {
			return nil, err
		}
	} else {
		valyentAuth = make(map[string]interface{})
	}
	valyentAuth["auth"] = base64.URLEncoding.EncodeToString([]byte("x:" + apiKey))

	b, err := json.Marshal(valyentAuth)
	if err != nil {
		return nil, err
	}

	dockerAuthProviders[registryHost] = b

	b, err = json.Marshal(dockerAuthProviders)
	if err != nil {
		return nil, err
	}

	dockerConfig["auths"] = b

	return json.Marshal(dockerConfig)
}
