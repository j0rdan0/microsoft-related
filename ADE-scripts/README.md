# Script to unlock/decrypt ADE disk for VM ( troubleshooting purposes)

Script to obtain the key file for decrypted ADE disks
---

Requires to have an AD SP created and have the following env variabler set:

* AZURE_CLIENT_ID
* AZURE_TENANT_ID
* AZURE_CLIENT_SECRET


SP can be created prior to this running:


az ad sp create-for-rbac -n "disk-decrypt"


The output can be used to fill in the environment variables required above



For getting a disk decrypted the following information needs to be filled in the config.json file:
---

* "SubscriptionID"
* "ResourceGroup"
* "DiskName"
* "VMName"


