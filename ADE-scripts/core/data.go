package core

type Config struct {
	SubscriptionID string
	ResourceGroup  string
	DiskName       string
	VMName         string
}

// this needs to be used to get secret value and after that unwrap the value using the key

type KVData struct {
	KeyVaultName  string
	SecretName    string
	SecretVersion string
	KeyName       string
	KeyVersion    string
}
