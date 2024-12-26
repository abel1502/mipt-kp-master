package internal

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
)

type AzureContainerBackup struct {
	containerClient *container.Client
}

func NewAzureContainerBackup(containerClient *container.Client) *AzureContainerBackup {
	return &AzureContainerBackup{
		containerClient,
	}
}

func NewDefaultAzureContainerBackup(containerName string) (*AzureContainerBackup, error) {
	defaultCred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, err
	}

	containerClient, err := container.NewClient(
		fmt.Sprintf("https://%s.blob.core.windows.net", containerName),
		defaultCred,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return NewAzureContainerBackup(containerClient), nil
}

func sanitizeName(blobName string) string {
	return strings.NewReplacer("/", "_", ".", "_").Replace(blobName)
}

func (a *AzureContainerBackup) BackupBlob(ctx context.Context, blobName string, fileName string, keepSnapshot bool) error {
	blobClient := a.containerClient.NewBlobClient(blobName)

	snapshotResp, err := blobClient.CreateSnapshot(ctx, nil)
	if err != nil {
		return err
	}

	snapshotClient, err := blobClient.WithSnapshot(*snapshotResp.Snapshot)
	if err != nil {
		return err
	}
	if !keepSnapshot {
		// To make sure we don't take up unnecessary storage space after we're done with the snapshot
		defer func() {
			_, _ = snapshotClient.Delete(ctx, nil)
		}()
	}

	if fileName == "" {
		fileName = fmt.Sprintf("./snapshot-%s-%s.bin", sanitizeName(blobName), *snapshotResp.Snapshot)
	}

	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = snapshotClient.DownloadFile(ctx, file, nil)
	if err != nil {
		return err
	}

	return nil
}

func (a *AzureContainerBackup) BackupAll(ctx context.Context, dirName string, keepSnapshots bool) error {
	blobPager := a.containerClient.NewListBlobsFlatPager(nil)

	for blobPager.More() {
		blobPage, err := blobPager.NextPage(ctx)
		if err != nil {
			return err
		}
		for _, blob := range blobPage.Segment.BlobItems {
			if *blob.Snapshot != "" || *blob.Deleted {
				continue
			}

			err := a.BackupBlob(ctx, *blob.Name, fmt.Sprintf("%s/%s.bin", dirName, sanitizeName(*blob.Name)), keepSnapshots)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
