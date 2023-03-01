---
# Update NSG Security Rule with CIDR source address prefix


Simple class for adding the CIDR your public IP is part of and add it into the NSG Security Rule for not manually doing this all the time manually.                                                                                                                                 

A Service Principal needs to be first created within the Subscription and given proper permissions to read NSGs and update Security rules. 

The Service Principal needs to be granted the following RBAC permissions at the subscription scope
		```
		- Microsoft.Network/networkSecurityGroups/securityRules/read
		- Microsoft.Network/networkSecurityGroups/securityRules/write
		-  Microsoft.Network/networkSecurityGroups/read
		```


An example of how to define the custom role needed for the SP can be found in the whitelist-app-role.json file. For role creation AZ cli can be used

	```
	az role definition create --role-definition whitelist-app-role.json
	```

In order to instantiate an NSGManager object the following parameters are required

```
Params for the NSGManager object:                                                                                                                                   
	@tenant_id: The Tenant ID of the Service Principal                                                                                        
	@app_id: The Service Principal App ID                                                                                                     
	@sub_id: Subscription ID                                                                                                                  
	@rg_name: Resource Group name where the NSG is created                                                                                    
	@nsg_name: NSG Name                                                                                                                       
	@rule_name: Security Rule name from the NSG                                                                                               
```

General requirements:

```
	- python3.x 
	- pip3, should be installed with: apt install python3-pip, yum install python3-pip etc
	- Python modules used need to be installed also: pip3 install -r requirements.txt
	- Azure Service Principal needs to be created and credentials/NSG details need to be provided to the NSGManager instance. The Service Principal can be created using AZ cli:
	az ad sp create-for-rbac --name whitelist-app --scopes /subscriptions/<subscription ID> --role whitelist-app-role
```



