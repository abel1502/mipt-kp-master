package app

import (
	"log"

	"github.com/abel1502/mipt-kp-m-test/internal/backup"
	"github.com/spf13/cobra"
)

func MakeCmdRoot(appName string) *cobra.Command {
	result := &cobra.Command{
		Use:   appName,
		Short: "Microsoft Azure Blob Storage backup utility",
		Long:  "Microsoft Azure Blob Storage backup utility",
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
	}

	result.AddCommand(CmdInit)
	result.AddCommand(CmdBackup)

	return result
}

var CmdInit = &cobra.Command{
	Use:   "init container_url [directory]",
	Short: "Initialize a new backup repository",
	Long:  "Initialize a new backup repository",
	Args:  cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		containerURL := args[0]
		directory := "."
		if len(args) > 1 {
			directory = args[1]
		}

		repo, err := backup.NewRepository(containerURL, directory)
		if err != nil {
			return err
		}

		err = repo.Close()
		if err != nil {
			return err
		}

		log.Printf("Successfully initialized backup repository for %v at %v", containerURL, directory)

		return nil
	},
}

var CmdBackup = &cobra.Command{
	Use:   "backup",
	Short: "Make a new incremental backup in the current repository",
	Long:  "Make a new incremental backup in the current repository",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		repo, err := backup.OpenRepository(".")
		if err != nil {
			return err
		}

		err = repo.TakeSnapshot(cmd.Context())
		if err != nil {
			return err
		}

		log.Printf("Successfully took a new snapshot")

		return nil
	},
}
