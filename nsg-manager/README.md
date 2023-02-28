---
# Update NSG Security Rule with CIDR source address prefix


```
Simple class for adding the CIDR your public IP is part of and add it into the NSG Security Rule for not manually doing this all the time 
manually.                                                                                                                                 
A Service Principal needs to be first created within the subscription and give proper permissions to read NSGs and update Security rules. 
Params:                                                                                                                                   
@tenant_id: The Tenant ID of the Service Principal                                                                                        
@app_id: The Service Principal App ID                                                                                                     
@sub_id: Subscription ID                                                                                                                  
@rg_name: Resource Group name where the NSG is created                                                                                    
@nsg_name: NSG Name                                                                                                                       
@rule_name: Security Rule name from the NSG                                                                                               
```

Requirements:

- python3.x 
- pip3, should be installed with: apt install python3-pip, yum install python3-pip etc
- Python modules used need to be installed also: pip3 install -r requirements.txt
- Azure Service Principal needs to be created and credentials/NSG details need to be provided to the NSGManager instance




