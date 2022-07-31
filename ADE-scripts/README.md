# Script to unlock/decrypt ADE disk for Windows VM

Script to obtain the key file for decrypting ADE disks
---

Requires to have an AAD app registered for programmatic access to the ARM API and have the following env variables set for app authentication:

* AZURE_CLIENT_ID
* AZURE_TENANT_ID
* AZURE_CLIENT_SECRET
* AZURE_OBJECT_ID

For getting an app created in your AAD tenant use:
---

```
scripts/init-app.sh
```

The output can be used to fill in the environment variables required above.

e.g:

```
echo "export AZURE_CLIENT_ID=0000-000-000-0000\n\
export AZURE_TENANT_ID=0000-000-000-0000\n\
export AZURE_CLIENT_SECRET=000-000-000-0000\n\
export AZURE_OBJECT_ID= 0000-000-000-0000" >> ~/.bashrc

source ~/.bashrc
```

For getting a disk decrypted the following information needs to be filled in the config.json file:
---

* "SubscriptionID"
* "ResourceGroup"
* "VMName"

For cleaning up the app created in your AAD tenant and remove the RBAC roles assigned to it use:
---

```
scripts\clean-app.sh
```





