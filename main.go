package main

import (
	"context"
	"errors"
	"flag"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/api/auth/approle"
	"github.com/hashicorp/vault/api/auth/kubernetes"
	"go.uber.org/multierr"
)

type TokenSource interface {
	Receive(ctx context.Context) (string, error)
	Cleanup(ctx context.Context) error
}

type TokenWriter interface {
	Write(ctx context.Context, data []byte) error
}

type App struct {
	source TokenSource
	dest   TokenWriter
}

func main() {
	var cfg Config
	opts := env.Options{
		Prefix: "VAULT_LOGIN_",
	}
	if err := env.ParseWithOptions(&cfg, opts); err != nil {
		slog.Error("could not parse config", "err", err)
		os.Exit(1)
	}

	parseFlags(&cfg)

	if err := cfg.Validate(); err != nil {
		slog.Error("invalid config", "err", err)
		os.Exit(1)
	}

	app, err := buildApp(cfg)
	if err != nil {
		slog.Error("could not build app", "err", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.TODO(), 15*time.Second)
	defer cancel()

	token, err := app.source.Receive(ctx)
	if err != nil {
		slog.Error("could not get token", "err", err)
		os.Exit(1)
	}
	slog.Info("Token received")

	if err := app.dest.Write(ctx, []byte(token)); err != nil {
		slog.Error("could not write token, trying to cleanup", "err", err)
		if err := app.source.Cleanup(ctx); err != nil {
			slog.Error("error while cleaning up token", "err", err)
		}
		os.Exit(1)
	}
	slog.Info("Wrote received token to configured storage")
}

type Config struct {
	AuthType                string `env:"AUTH_TYPE"`
	AuthRole                string `env:"AUTH_ROLE"`
	AuthMount               string `env:"AUTH_MOUNT"`
	AuthApproleSecretId     string `env:"AUTH_APPROLE_SECRET_ID"`
	AuthApproleSecretIdFile string `env:"AUTH_APPROLE_SECRET_ID_FILE"`

	OutputType            string `env:"OUTPUT_TYPE"`
	OutputSecretName      string `env:"OUTPUT_SECRET_NAME"`
	OutputSecretNamespace string `env:"OUTPUT_SECRET_NAMESPACE"`
	OutputSecretKey       string `env:"OUTPUT_SECRET_KEY"`
}

func parseFlags(cfg *Config) {
	flag.StringVar(&cfg.AuthType, "auth-type", cfg.AuthType, "Type of the authentication")
	flag.StringVar(&cfg.AuthRole, "auth-role", cfg.AuthRole, "Role for authentication")
	flag.StringVar(&cfg.AuthMount, "auth-mount", cfg.AuthMount, "Mount point for authentication")
	flag.StringVar(&cfg.AuthApproleSecretId, "auth-approle-secret-id", cfg.AuthApproleSecretId, "Approle Secret ID for authentication")
	flag.StringVar(&cfg.AuthApproleSecretIdFile, "auth-approle-secret-id-file", cfg.AuthApproleSecretIdFile, "Approle Secret ID file for authentication")

	flag.StringVar(&cfg.OutputType, "output-type", cfg.OutputType, "Type of output")
	flag.StringVar(&cfg.OutputSecretName, "output-secret-name", cfg.OutputSecretName, "Output secret name")
	flag.StringVar(&cfg.OutputSecretNamespace, "output-secret-namespace", cfg.OutputSecretNamespace, "Output secret namespace")
	flag.StringVar(&cfg.OutputSecretKey, "output-secret-key", cfg.OutputSecretKey, "Output secret key")

	flag.Parse()
}

func (c *Config) Validate() (err error) {
	if strings.TrimSpace(c.AuthType) == "" {
		err = multierr.Append(err, errors.New("no vault auth type provided"))
	}

	if strings.TrimSpace(c.AuthRole) == "" {
		err = multierr.Append(err, errors.New("no vault role provided"))
	}

	if strings.TrimSpace(c.OutputType) == "" {
		err = multierr.Append(err, errors.New("no output type provided"))
	}

	return
}

func buildApp(cfg Config) (*App, error) {
	authType, err := buildAuthMethod(cfg)
	if err != nil {
		return nil, err
	}

	tokenSource, err := NewVaultTokenSource(cfg, authType)
	if err != nil {
		return nil, err
	}

	tokenWriter, err := buildOutput(cfg)
	if err != nil {
		return nil, err
	}

	return &App{
		source: tokenSource,
		dest:   tokenWriter,
	}, nil
}

func buildAuthMethod(cfg Config) (api.AuthMethod, error) {
	switch cfg.AuthType {
	case "kubernetes":
		var opt []kubernetes.LoginOption
		if strings.TrimSpace(cfg.AuthMount) != "" {
			opt = append(opt, kubernetes.WithMountPath(cfg.AuthMount))
		}
		return kubernetes.NewKubernetesAuth(cfg.AuthRole, opt...)

	case "approle":
		var opt []approle.LoginOption
		if strings.TrimSpace(cfg.AuthMount) != "" {
			opt = append(opt, approle.WithMountPath(cfg.AuthMount))
		}
		secretId := &approle.SecretID{
			FromFile:   cfg.AuthApproleSecretIdFile,
			FromString: cfg.AuthApproleSecretId,
		}
		return approle.NewAppRoleAuth(cfg.AuthRole, secretId, opt...)
	default:
		return nil, errors.New("no valid auth type supplied")
	}
}

func buildOutput(cfg Config) (TokenWriter, error) {
	switch cfg.OutputType {
	case "stdout":
		return NewStdoutWriter(), nil
	case "kubernetes-secret":
		return NewKubernetesSecretWriter(cfg)
	case "file":
		return NewFileWriter(cfg)
	default:
		return nil, errors.New("no valid output type provided")
	}
}
