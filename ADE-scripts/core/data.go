package core

// this needs to be used to get secret value and after that unwrap the value using the key

type KVData struct {
	KeyVaultName  string
	SecretName    string
	SecretVersion string
	KeyName       string
	KeyVersion    string
}

type opts struct {
	SubscriptionID string `short:"s" long:"subscription" description:"Subscription ID" required:"true"`
	ResourceGroup  string `short:"g" long:"resource-group" description:"Resource group" required:"true"`
	VMName         string `short:"n" long:"name" description:"VM Name" required:"true"`
}
