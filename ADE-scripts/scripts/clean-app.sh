#!/bin/bash

VM_NAME=$(cat ../config.json | jq .VMName |tr -d "\"")
RG=$(cat ../config.json | jq .ResourceGroup | tr -d "\"")
APP_NAME="decrypt-app"

RED='\033[0;31m'
NC='\033[0m'

OS_DISK=$(az vm show -g $RG -n $VM_NAME| jq .storageProfile.osDisk.managedDisk.id| tr -d "\"")
KV=$(az vm show -g $RG -n $VM_NAME | jq .'resources[0]'.settings.KeyVaultResourceId | tr -d "\"")

AZURE_CLIENT_ID=$(az ad app list --display-name $APP_NAME --query [].appId -o tsv)
AZURE_OBJECT_ID=$(az ad sp list --display-name $APP_NAME --query [].id -o tsv)

az role assignment create --assignee $AZURE_OBJECT_ID --role "Reader" --scope $OS_DISK 1>/dev/null  && printf "${RED}[*]${NC} deleted IAM access rules for OS Disk ${RED}$APP_NAME${NC}\n"
az role assignment create --assignee $AZURE_OBJECT_ID --role "Contributor" --scope $KV 1>/dev/null  && printf "${RED}[*]${NC} deleted IAM access rules for KeyVault ${RED}$APP_NAME${NC}\n"

# give chance to have roles removed
sleep 10

# delete app 
az ad app delete --id $AZURE_CLIENT_ID && printf "${RED}[*]${NC} deleted app ${RED}$APP_NAME${NC}\n"



