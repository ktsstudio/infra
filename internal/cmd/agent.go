package cmd

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/cenkalti/backoff/v5"
	"github.com/shirou/gopsutil/v4/process"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"

	"github.com/infrahq/infra/internal"
	"github.com/infrahq/infra/internal/logging"
	"github.com/infrahq/infra/internal/repeat"
)

func newAgentCmd() *cobra.Command {
	return &cobra.Command{
		Use:    "agent",
		Short:  "Start the Infra agent",
		Long:   "Start the Infra agent that runs to sync Infra state",
		Args:   NoArgs,
		Hidden: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			infraDir, err := initInfraHomeDir()
			if err != nil {
				return fmt.Errorf("infra home directory: %w", err)
			}
			logging.UseFileLogger(filepath.Join(infraDir, "agent.log"))

			group, ctx := errgroup.WithContext(context.Background())
			group.Go(func() error {
				backOff := &backoff.ExponentialBackOff{
					InitialInterval:     time.Minute,
					MaxInterval:         2 * time.Minute,
					RandomizationFactor: 0.05,
					Multiplier:          1.5,
				}
				waiter := repeat.NewWaiterWithMaxElapsed(backOff, 5*time.Minute)
				for {
					if err := syncKubeConfig(ctx); err != nil {
						logging.L.Warn().Err(err).Msg("failed to sync kubeconfig")
					} else {
						waiter.Reset()
					}
					if err := waiter.Wait(ctx); err != nil {
						return err
					}
				}
			})
			// add the next agent task here

			logging.Infof("starting infra agent (%s)", internal.FullVersion())
			err = group.Wait()
			logging.L.Error().Err(err).Msg("infra agent exit")
			return err
		},
	}
}

// configAgentRunning checks if the agent process stored in config is still running
func configAgentRunning() (bool, error) {
	pid, err := readStoredAgentProcessID()
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			// this is the first time the agent is running, suppress the error and continue
			logging.Debugf("%s", err.Error())
			return false, nil
		}
		return false, err
	}

	if pid < 0 || pid > math.MaxInt32 {
		return false, fmt.Errorf("invalid pid: %d", pid)
	}
	return processRunning(int32(pid))
}

func processRunning(pid int32) (bool, error) {
	if pid == 0 { // on windows pid 0 is the system idle process, it will be running but its not the agent
		return false, nil
	}

	running, err := process.PidExists(pid)
	if err != nil {
		return false, err
	}

	return running, nil
}

// execAgent executes the agent in the background and stores its process ID in configuration
func execAgent() error {
	// find the location of the infra executable running
	// this ensures we run the agent for the correct version
	infraExe, err := os.Executable()
	if err != nil {
		return err
	}

	// use the current infra executable to start the agent
	// The agent outlives this CLI invocation, so a background context is intentional.
	// #nosec G204 -- Executable is resolved from os.Executable, never from request or shell input.
	cmd := exec.CommandContext(context.Background(), infraExe, "agent")
	if err := cmd.Start(); err != nil {
		return err
	}

	logging.Debugf("agent started, pid: %d", cmd.Process.Pid)

	return writeAgentConfig(cmd.Process.Pid)
}

func readStoredAgentProcessID() (int, error) {
	infraDir, err := initInfraHomeDir()
	if err != nil {
		return 0, err
	}

	// #nosec G304 -- Local file path explicitly selected by the CLI user or trusted administrator configuration.
	agentConfig, err := os.Open(filepath.Join(infraDir, "agent.pid"))
	if err != nil {
		return 0, err
	}
	defer func() {
		if err := agentConfig.Close(); err != nil {
			logging.Debugf("close agent PID file: %s", err)
		}
	}()

	var pid int

	_, err = fmt.Fscanf(agentConfig, "%d\n", &pid)
	if err != nil {
		return 0, err
	}

	return pid, nil
}

// writeAgentProcessConfig saves details about the agent to config
func writeAgentConfig(pid int) (retErr error) {
	infraDir, err := initInfraHomeDir()
	if err != nil {
		return err
	}

	// #nosec G304 -- Local file path explicitly selected by the CLI user or trusted administrator configuration.
	agentConfig, err := os.Create(filepath.Join(infraDir, "agent.pid"))
	if err != nil {
		return err
	}
	defer func() { retErr = errors.Join(retErr, agentConfig.Close()) }()

	_, err = agentConfig.WriteString(fmt.Sprintf("%d\n", pid))
	if err != nil {
		return err
	}

	return nil
}

// syncKubeConfig updates the local kubernetes configuration from Infra grants
func syncKubeConfig(_ context.Context) error {
	opts, err := defaultClientOpts()
	if err != nil {
		return err
	}
	client, err := NewAPIClient(opts)
	if err != nil {
		return err
	}
	client.Name = "agent"

	return updateKubeconfig(client)
}
