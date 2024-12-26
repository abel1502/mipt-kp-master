package main

import (
	"context"
	"flag"

	"github.com/abel1502/mipt-kp-m-test/internal"
)

func main() {
	var containerName = flag.String("container", "mycontainer", "Container name")
	var blobName = flag.String("blob", "myblob", "Blob name")
	var keepSnapshot = flag.Bool("keep", false, "Keep snapshot in cloud")
	flag.Parse()

	err := doMain(*containerName, *blobName, *keepSnapshot)
	if err != nil {
		panic(err)
	}
}

func doMain(
	containerName string,
	blobName string,
	keepSnapshot bool,
) error {
	backup, err := internal.NewDefaultAzureContainerBackup(containerName)
	if err != nil {
		return err
	}

	err = backup.BackupBlob(context.Background(), blobName, "", keepSnapshot)
	if err != nil {
		return err
	}

	return nil
}
