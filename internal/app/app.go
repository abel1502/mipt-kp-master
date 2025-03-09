package app

import (
	"log"
	"net/url"
	"path"
	"path/filepath"

	"github.com/abel1502/mipt-kp-m-test/internal/backup"
	"github.com/spf13/cobra"
)

var argDirectory string

func MakeCmdRoot(appName string) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   appName,
		Short: "Microsoft Azure Blob Storage backup utility",
		Long:  "Microsoft Azure Blob Storage backup utility",
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
	}

	rootCmd.AddCommand(CmdInit)
	rootCmd.AddCommand(CmdBackup)

	rootCmd.PersistentFlags().StringVarP(&argDirectory, "directory", "C", ".", "Working directory")

	return rootCmd
}

var CmdInit = &cobra.Command{
	Use:   "init container_url [directory_name]",
	Short: "Initialize a new backup repository",
	Long:  "Initialize a new backup repository",
	Args:  cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		containerURL := args[0]

		parsedURL, err := url.Parse(containerURL)
		if err != nil {
			return err
		}

		backupName := path.Base(parsedURL.Path)
		if len(args) == 2 {
			backupName = args[1]
		}

		directory := filepath.Join(argDirectory, backupName)

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
		repo, err := backup.OpenRepository(argDirectory)
		if err != nil {
			return err
		}
		defer repo.Close()

		err = repo.TakeSnapshot(cmd.Context())
		if err != nil {
			return err
		}

		log.Printf("Successfully took a new snapshot")

		return nil
	},
}
