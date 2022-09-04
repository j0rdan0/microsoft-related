package core

import (
	"bytes"
	"context"
	"io"
	"log"
	"os"

	"net/http"

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

func CheckResourceGroup() (string, error) {

	//// resource group for keyvault can be different than resource group for disk!! need to handle this in the future

	// will use Azure Resource Graph API to find RG for KV automatically:
	// see https://docs.microsoft.com/en-us/rest/api/azureresourcegraph/resourcegraph(2020-04-01-preview)/resources/resources?tabs=HTTP

	URL := "https://management.azure.com/providers/Microsoft.ResourceGraph/resources?api-version=2020-04-01-preview"
	client := &http.Client{}
	body := []byte(`{
		"query":  "Resources | where type =~ 'Microsoft.KeyVault/vaults'"
	}`)
	req, err := http.NewRequest("POST", URL, bytes.NewBuffer(body))
	token, err := GetToken(true)
	if err != nil {
		return "", nil
	}
	t := "Bearer " + token

	req.Header.Set("Authorization", t)
	req.Header.Set("Content-type", "application/json")

	if err != nil {
		return "", nil
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", nil
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", nil
	}
	return string(data), nil

	// need to process the returned data to confirm KV Resource group

}
