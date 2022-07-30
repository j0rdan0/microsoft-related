# Script to unlock/decrypt ADE disk for VM ( troubleshooting purposes)

Script to obtain the key file for decrypting ADE disks
---

Requires to have an AD SP created and have the following env variables set:


* AZURE_CLIENT_ID
* AZURE_TENANT_ID
* AZURE_CLIENT_SECRET
* AZURE_OBJECT_ID

For getting an app created in your AAD run:
---

```
scripts/init-app.sh
```

The output can be used to fill in the environment variables required above


For getting a disk decrypted the following information needs to be filled in the config.json file:
---

* "SubscriptionID"
* "ResourceGroup"
* "DiskName"
* "VMName"

// need to find a way of finding the disk name directly from the VM name 

For cleaning up the app created in your AAD tenant run:
---

```
scripts\clean-app.sh
```





