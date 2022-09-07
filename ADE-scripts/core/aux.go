package core

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resourcegraph/armresourcegraph"
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

func checkResourceGroup(keyvault string, subscription string, cred *azidentity.DefaultAzureCredential) (string, error) {

	ctx := context.Background()
	client, err := armresourcegraph.NewClient(cred, nil)
	if err != nil {
		return "", err
	}
	query := fmt.Sprintf("Resources | project id, name, resourceGroup | where name == '%s'", keyvault)
	res, err := client.Resources(ctx,
		armresourcegraph.QueryRequest{
			Query: to.Ptr(query),
			Subscriptions: []*string{
				to.Ptr(subscription)},
		},
		nil)
	if err != nil {
		return "", err
	}

	// uglies way to get data back as a struct, but too tired now to refactor
	data := make([]byte, 0)
	data, err = json.Marshal(res.Data)
	if err != nil {
		return "", err
	}
	var d []RG

	err = json.Unmarshal(data, &d)
	if err != nil {
		return "", nil
	}
	var rg string
	for _, elem := range d {
		rg = fmt.Sprint(elem.ResourceGroup)
	}
	fmt.Println(rg)
	return rg, nil
}
