#!/bin/bash

VM_NAME=$(cat ../config.json | jq .VMName |tr -d "\"")
RG=$(cat ../config.json | jq .ResourceGroup | tr -d "\"")
APP_NAME="decrypt-app"

RED='\033[0;31m'
NC='\033[0m'

# get resource IDs in order to remove RBAC roles
OS_DISK=$(az vm show -g $RG -n $VM_NAME| jq .storageProfile.osDisk.managedDisk.id| tr -d "\"")
KV=$(az vm show -g $RG -n $VM_NAME | jq .'resources[0]'.settings.KeyVaultResourceId | tr -d "\"")
VM=$(az vm show -g $RG -n $VM_NAME | jq .id | tr -d "\"")

AZURE_CLIENT_ID=$(az ad app list --display-name $APP_NAME --query [].appId -o tsv)
AZURE_OBJECT_ID=$(az ad sp list --display-name $APP_NAME --query [].id -o tsv)


# delete RBAC roles
az role assignment delete --assignee $AZURE_OBJECT_ID --role "Reader" --scope $OS_DISK 1>/dev/null  && printf "${RED}[*]${NC} deleted Reader role for OS Disk resource ${RED}$OS_DISK${NC}\n"
az role assignment delete --assignee $AZURE_OBJECT_ID --role "Contributor" --scope $KV 1>/dev/null  && printf "${RED}[*]${NC} deleted Contributor role for KeyVault resource ${RED}$KV${NC}\n"
az role assignment delete --assignee $AZURE_OBJECT_ID --role "Reader" --scope $VM 1>/dev/null  && printf "${RED}[*]${NC} deleted Reader role for VM resource ${RED}$VM${NC}\n"


# delete app 
az ad app delete --id $AZURE_CLIENT_ID && printf "${RED}[*]${NC} deleted app ${RED}$APP_NAME${NC}\n"



