// (c) Cartesi and individual authors (see AUTHORS)
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package deps

import (
	"context"
	"log/slog"
	"os/signal"
	"syscall"

	"github.com/cartesi/rollups-node/internal/deps"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:     "run-deps",
	Short:   "Run node dependencies with Docker",
	Example: examples,
	Run:     run,
}

const examples = `# Run all deps:
cartesi-rollups-cli run-deps`

var depsConfig = deps.NewDefaultDepsConfig()

func init() {
	Cmd.Flags().StringVar(&depsConfig.Postgres.DockerImage, "postgres-docker-image",
		deps.DefaultPostgresDockerImage,
		"Postgress docker image name")

	Cmd.Flags().StringVar(&depsConfig.Postgres.Port, "postgres-mapped-port",
		deps.DefaultPostgresPort,
		"Postgres local listening port number")

	Cmd.Flags().StringVar(&depsConfig.Postgres.Password, "postgres-password",
		deps.DefaultPostgresPassword,
		"Postgres password")

	Cmd.Flags().StringVar(&depsConfig.Devnet.DockerImage, "devnet-docker-image",
		deps.DefaultDevnetDockerImage,
		"Devnet docker image name")

	Cmd.Flags().StringVar(&depsConfig.Devnet.Port, "devnet-mapped-port",
		deps.DefaultDevnetPort,
		"devnet local listening port number")
}

func run(cmd *cobra.Command, args []string) {
	ctx, cancel := signal.NotifyContext(cmd.Context(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	depsContainers, err := deps.Run(ctx, *depsConfig)
	cobra.CheckErr(err)

	slog.Info("All dependencies are up")

	<-ctx.Done()

	err = deps.Terminate(context.Background(), depsContainers)
	cobra.CheckErr(err)
}
