package cmd

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/infrahq/infra/internal"
	"github.com/infrahq/infra/internal/logging"
)

func newVersionCmd(cli *CLI) *cobra.Command {
	return &cobra.Command{
		Use:     "version",
		Short:   "Display the Infra version",
		GroupID: groupOther,
		Args:    NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return version(cli)
		},
	}
}

func version(cli *CLI) (retErr error) {
	ctx := context.Background()

	w := tabwriter.NewWriter(cli.Stdout, 0, 0, 1, ' ', tabwriter.AlignRight)
	defer func() { retErr = errors.Join(retErr, w.Flush()) }()

	if _, err := fmt.Fprintln(w); err != nil {
		logging.Debugf("write command output: %s", err)
	}
	if _, err := fmt.Fprintln(w, "Client:\t", strings.TrimPrefix(internal.FullVersion(), "v")); err != nil {
		logging.Debugf("write command output: %s", err)
	}

	// Note that we use the client to get this version, but it is in fact the server version
	client, err := cli.apiClient()
	if err != nil {
		if _, err := fmt.Fprintln(w, "Server:\t", "disconnected"); err != nil {
			logging.Debugf("write command output: %s", err)
		}
		logging.Debugf("%s", err.Error())
		if _, err := fmt.Fprintln(w); err != nil {
			logging.Debugf("write command output: %s", err)
		}

		return nil
	}

	config, err := currentHostConfig()
	if err != nil {
		return err
	}

	// Don't bother printing the server version for SaaS
	if strings.HasSuffix(config.Host, ".infrahq.com") {
		if _, err := fmt.Fprintln(w); err != nil {
			logging.Debugf("write command output: %s", err)
		}
		return nil
	}

	version, err := client.GetServerVersion(ctx)
	if err != nil {
		if _, err := fmt.Fprintln(w, "Server:\t", "disconnected"); err != nil {
			logging.Debugf("write command output: %s", err)
		}
		logging.Debugf("%s", err.Error())
		if _, err := fmt.Fprintln(w); err != nil {
			logging.Debugf("write command output: %s", err)
		}

		return nil
	}

	if _, err := fmt.Fprintln(w, "Server:\t", strings.TrimPrefix(version.Version, "v")); err != nil {
		logging.Debugf("write command output: %s", err)
	}
	if _, err := fmt.Fprintln(w); err != nil {
		logging.Debugf("write command output: %s", err)
	}

	return nil
}
