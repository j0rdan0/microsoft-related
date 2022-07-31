package core

import (
	"context"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute"
	flags "github.com/jessevdk/go-flags"
)

func authenticate() (*azidentity.DefaultAzureCredential, error) {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {

		return nil, err
	}
	return cred, err
}

func isEncrypted(r armcompute.DisksClientGetResponse) bool {
	if *r.Properties.Encryption.Type == "EncryptionAtRestWithPlatformKey" {
		return true
	}
	return false
}
func getOSDisk(cred *azidentity.DefaultAzureCredential, subscriptionId string, resourceGroup string, vmName string) (*string, error) {
	vmClient, err := armcompute.NewVirtualMachinesClient(subscriptionId, cred, nil)
	if err != nil {

		return nil, err
	}
	resp, err := vmClient.Get(context.Background(), resourceGroup, vmName, nil)
	if err != nil {

		return nil, err
	}
	diskName := resp.Properties.StorageProfile.OSDisk.Name
	if diskName == nil {
		return nil, err
	}

	return diskName, err
}

func HandleError(err error) {
	if err != nil {
		log.Fatal(err)
		os.Exit(-1)
	}
}

func getArgs() {
	if _, err := parser.Parse(); err != nil {
		switch flagsErr := err.(type) {
		case flags.ErrorType:
			if flagsErr == flags.ErrHelp {
				os.Exit(0)
			}
			os.Exit(1)
		default:
			os.Exit(1)
		}
	}
}
