package core

import (
	"context"
	"encoding/json"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute"
)

func readConfig(filename string) (*Config, error) {
	var config Config
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	err = json.NewDecoder(f).Decode(&config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

/*

need the following env variables to be set for authentication:

AZURE_CLIENT_ID
AZURE_TENANT_ID
AZURE_CLIENT_SECRET

also need

AZURE_OBJECT_ID for the SP to be able to add the KV access policy

*/
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
	/*cred, err := authenticate()
	if err != nil {

		return nil, err
	}
	*/

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
