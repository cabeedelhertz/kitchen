package service

import (
	"context"
	"errors"
	"fmt"
	"kitchen/pkg/common/config"
	"kitchen/pkg/common/logging"
	"math"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func init() {
	config.RegisterDefault("grpc_port", 50051)
	config.RegisterDefault("http_port", 8080)
	config.RegisterDefault("read_header_timeout", time.Second)
	config.RegisterDefault("read_timeout", 30*time.Second)
	config.RegisterDefault("write_timeout", 30*time.Second)
}

// Ensure Config conforms to ValidatableConfig
var _ ValidatableConfig = Config{}

// ValidatableConfig is an interface a Config can conform to in order to be
// validated
type ValidatableConfig interface {
	Root() config.Base
	Validate(context.Context) error
}

// LoadConfigForCommand loads the config for the supplied command
func LoadConfigForCommand(cmd *cobra.Command, c ValidatableConfig) error {
	return LoadAndValidateConfig(cmd.Context(), cmd.Flags(), c)
}

// LoadAndValidateConfig loads the config from the supplied context. If the
// config fails to load or is not valid, and error is returned
func LoadAndValidateConfig(ctx context.Context, flags *pflag.FlagSet, c ValidatableConfig) error {
	if err := config.LoadInto(flags, c); err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	if err := logging.Configure(c.Root()); err != nil {
		return fmt.Errorf("failed to configure logging: %w", err)
	}
	if err := c.Validate(ctx); err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}
	return nil
}

// Config is the gRPC config struct
type Config struct {
	config.Base       `config:",squash"`
	ValidIssuers      []string      `config:"valid_issuers"`
	GrpcPort          int           `config:"grpc_port"`
	HttpPort          int           `config:"http_port"`
	ReadHeaderTimeout time.Duration `config:"read_header_timeout"`
	ReadTimeout       time.Duration `config:"read_timeout"`
	WriteTimeout      time.Duration `config:"write_timeout"`
}

// Root returns the root config
func (c Config) Root() config.Base {
	return c.Base
}

// Validate validates this config
func (c Config) Validate(ctx context.Context) error {

	// Validate our base config
	violations := make(map[string]error)

	// Validate the grpc port is in a valid range
	if err := validatePort(c.GrpcPort); err != nil {
		violations["grpc_port"] = err
	}
	if err := validatePort(c.HttpPort); err != nil {
		violations["http_port"] = err
	}
	if len(violations) > 0 {
		return errors.New("validation failed")
	}
	return nil
}

// validatePort validates that the port number is in a valid range
func validatePort(port int) error {
	if port < 1 || port > math.MaxUint16 {
		return fmt.Errorf("invalid port %d, must be in range %d - %d", port, 1, math.MaxUint16)
	}
	return nil
}
