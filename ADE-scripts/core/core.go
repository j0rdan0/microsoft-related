package core

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/keyvault/armkeyvault"
	ct "github.com/daviddengcn/go-colortext"
)

func GetDiskEncryptionType(kvdata *KVData) {
	cred, err := authenticate()
	if err != nil {
		log.Fatal("failed authenticating: ", err)
		return
	}
	config, err := readConfig("./config.json")
	if err != nil {
		log.Fatal("error reading config: ", err)
		return
	}

	diskClient, err := armcompute.NewDisksClient(config.SubscriptionID, cred, nil)
	if err != nil {
		log.Fatal("failed creating disk client: ", err)
		return
	}
	resp, err := diskClient.Get(context.Background(), config.ResourceGroup, config.DiskName, nil)
	if err != nil {
		log.Fatal("failed getting disk: ", err)
		return
	}

	if !isEncrypted(resp) {
		ct.Background(ct.Red, false)
		log.Printf("*** disk %s is not encrypted, no decryption needed", config.DiskName)
		ct.ResetColor()
		os.Exit(-1)

	} else {
		encVersion, err := getADEVersion()
		if err != nil {
			log.Fatal(err)
			return
		}
		fmt.Printf("[***] disk %s is", config.DiskName)
		ct.Background(ct.Red, false)
		fmt.Printf(" encrypted ")
		ct.ResetColor()
		fmt.Printf("using ADE version %d\n", encVersion)
	}

	var secretURL string
	var keyURL string
	for _, elem := range resp.Properties.EncryptionSettingsCollection.EncryptionSettings {
		secretURL = *elem.DiskEncryptionKey.SecretURL
		keyURL = *elem.KeyEncryptionKey.KeyURL
	}

	secretURLData := strings.Split(secretURL, "/")
	keyURLData := strings.Split(keyURL, "/")

	// secretURLData[2] = keyvault name
	// secretURLData[4] = secretName
	// secretURLData[5] = secretVersion

	// keyURLData[2] = keyvault name
	// keyURLData[4] = keyname
	// keyURLData[5] = key version

	kvdata.SecretName = secretURLData[4]
	kvdata.SecretVersion = secretURLData[5]
	kvdata.KeyVaultName = secretURLData[2]
	kvdata.KeyName = keyURLData[4]
	kvdata.KeyVersion = keyURLData[5]

}

func getADEVersion() (int, error) {
	config, err := readConfig("./config.json")
	if err != nil {
		return -1, err
	}
	cred, err := authenticate() // autenticate the app to Azure
	if err != nil {

		return -1, err
	}
	extClient, err := armcompute.NewVirtualMachineExtensionsClient(config.SubscriptionID, cred, nil)
	if err != nil {
		return -1, err
	}

	resp, err := extClient.Get(context.Background(), config.ResourceGroup, config.VMName, "AzureDiskEncryption", nil)
	if err != nil {
		return -1, nil
	}

	version, err := strconv.Atoi(strings.Split(*resp.Properties.TypeHandlerVersion, ".")[0])
	if err != nil {
		return -1, err
	}
	return version, nil

}

func WriteBEKFile(value string) {
	data, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		fmt.Println(err)
	}

	nBytes := len(data)
	f, err := os.Create("bek.file")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer f.Close()
	ct.Background(ct.Blue, false)
	log.Println("[***] writing BEK file")
	ct.ResetColor()
	n, err := f.Write(data)
	if err != nil {
		log.Fatal(err)
		return
	}
	if n == nBytes {

		log.Printf("[**] all %d bytes written to file bek.file\n", n)

	} else {
		log.Fatal("only %d bytes written to file\n", n)

	}
}

func GetToken() (string, error) {

	URL := "https://login.microsoftonline.com/" + os.Getenv("AZURE_TENANT_ID") + "/oauth2/v2.0/token"

	body := url.Values{}
	body.Add("client_id", os.Getenv("AZURE_CLIENT_ID"))
	body.Add("scope", "https://vault.azure.net/.default")
	body.Add("grant_type", "client_credentials")
	body.Add("client_secret", os.Getenv("AZURE_CLIENT_SECRET"))
	resp, err := http.PostForm(URL, body)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", nil
	}

	temp := strings.Split(string(respData), ":")[4]
	temp = strings.Trim(temp, `"`)

	return (strings.Split(temp, `"}`)[0]), nil

}

// need to replace the implementation with direct REST API, as this shit doesn`t work
// PUT https://management.azure.com/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/
// providers/Microsoft.KeyVault/vaults/{vaultName}/accessPolicies/{operationKind}?api-version=2021-10-01

func SetAccessPolicy(kvdata *KVData) (bool, error) {
	config, err := readConfig("./config.json")
	if err != nil {
		return false, err
	}
	cred, err := authenticate() // autenticate the app to Azure
	if err != nil {

		return false, err
	}
	client, err := armkeyvault.NewVaultsClient(config.SubscriptionID, cred, nil)
	if err != nil {
		return false, err
	}
	oid := os.Getenv("AZURE_OBJECT_ID")    // this calls requires an SP, not an APP
	tenant := os.Getenv("AZURE_TENANT_ID") // this seems to have a problem fetching env variable, need to check why

	params := armkeyvault.VaultAccessPolicyParameters{
		Properties: &armkeyvault.VaultAccessPolicyProperties{
			AccessPolicies: []*armkeyvault.AccessPolicyEntry{
				{

					ObjectID: to.Ptr(oid),
					Permissions: &armkeyvault.Permissions{
						Keys: []*armkeyvault.KeyPermissions{
							to.Ptr(armkeyvault.KeyPermissionsUnwrapKey),
						},
						Secrets: []*armkeyvault.SecretPermissions{
							to.Ptr(armkeyvault.SecretPermissionsGet),
						},
					},
					TenantID: to.Ptr(tenant),
				},
			},
		},
	}
	// resource group for keyvault can be different than resource group for disk!! need to adapt to this
	_, err = client.UpdateAccessPolicy(context.Background(), config.ResourceGroup, strings.Split(kvdata.KeyVaultName, ".")[0], armkeyvault.AccessPolicyUpdateKindAdd, params, nil)
	if err != nil {
		return false, err
	}

	log.Printf("[***] created KeyvaultAccessPolicy\n")

	return true, nil

}

func GetSecret(token string, kvdata *KVData) (string, error) {

	endpoint := "https://" + kvdata.KeyVaultName + "//secrets/" + kvdata.SecretName + "/" + kvdata.SecretVersion + "?api-version=7.3"

	client := &http.Client{}
	request, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return "", err
	}

	token = "Bearer " + token
	request.Header.Add("Authorization", token)
	resp, err := client.Do(request)

	if err != nil {
		return "", nil
	}

	defer resp.Body.Close()

	buff, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", nil
	}

	temp := strings.Split(string(buff), ":")[1]
	temp = strings.Split(temp, `"`)[1]

	log.Printf("[***] got secret`s %s value\n", kvdata.SecretName)

	return temp, nil

}

func UnwrapSecret(secret string, token string, kvdata *KVData) (string, error) {

	endpoint := "https://" + kvdata.KeyVaultName + "//keys/" + kvdata.KeyName + "/" + kvdata.KeyVersion + "/unwrapkey?api-version=7.3"

	body := map[string]string{"alg": "RSA-OAEP", "value": secret}
	json, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	token = "Bearer " + token
	client := &http.Client{Transport: nil}
	request, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(json))
	if err != nil {
		return "", err
	}
	request.Header.Add("Authorization", token)
	request.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// some value preprocessing before returning it for file write
	temp := string(data)
	temp = strings.Split(temp, ":")[3]
	temp = strings.Trim(temp, `"`)
	temp = strings.Trim(temp, `"}`)

	temp = strings.Replace(temp, "-", "+", -1)
	temp = strings.Replace(temp, "_", "/", -1)

	if len(temp)%4 == 2 {
		temp += "=="
	} else if len(temp)%4 == 3 {
		temp += "="
	}

	log.Printf("[***] unwrapped secret %s using key %s\n", kvdata.SecretName, kvdata.KeyName)

	return temp, nil

}

//only need the unwrapKey and get secret perms

func SetAccessPolicyREST(kvdata *KVData) (bool, error) {
	// to implement this function using REST API

	return true, nil
}
