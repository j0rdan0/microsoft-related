#!/bin/bash

# get information from config.json file needed to set IAM access control
VM_NAME=$(cat ../config.json | jq .VMName |tr -d "\"")
RG=$(cat ../config.json | jq .ResourceGroup | tr -d "\"")
APP_NAME="decrypt-app"

RED='\033[0;31m'
NC='\033[0m'

#create app 
az ad app create --display-name $APP_NAME 1>/dev/null && printf "${RED}[*]${NC} created app ${RED}$APP_NAME${NC}\n"

# get app id
AZURE_CLIENT_ID=$(az ad app list --display-name $APP_NAME --query [].appId -o tsv)

#create sp 
az ad sp create --id $AZURE_CLIENT_ID 1>/dev/null && printf "${RED}[*]${NC} created service principal for app ${RED}$APP_NAME${NC}\n"

AZURE_OBJECT_ID=$(az ad sp list --display-name $APP_NAME --query [].id -o tsv)
AZURE_TENANT_ID=$(az ad sp list --display-name $APP_NAME --query [].appOwnerOrganizationId -o tsv)

printf "${RED}[*]${NC} creating client secret for app ${RED}$APP_NAME${NC}\n"

# create client secret
AZURE_CLIENT_SECRET=$(az ad app credential reset --id $AZURE_CLIENT_ID --append | jq .password | tr -d "\"")


printf "${RED}[*]${NC} creating RBAC rules\n"

# get os disk resource URI and KV resource URI
OS_DISK=$(az vm show -g $RG -n $VM_NAME| jq .storageProfile.osDisk.managedDisk.id| tr -d "\"")
KV=$(az vm show -g $RG -n $VM_NAME | jq .'resources[0]'.settings.KeyVaultResourceId | tr -d "\"")

# create IAM access roles
az role assignment create --assignee $AZURE_OBJECT_ID --role "Reader" --scope $OS_DISK 1>/dev/null
az role assignment create --assignee $AZURE_OBJECT_ID --role "Contributor" --scope $KV 1>/dev/null

printf "${RED}[*]${NC} save the below env variables for authentication\n\n"

# output env variables needed for authentication
echo "export AZURE_OBJECT_ID=$AZURE_OBJECT_ID"
echo "export AZURE_CLIENT_ID=$AZURE_CLIENT_ID"
echo "export AZURE_TENANT_ID=$AZURE_TENANT_ID"
echo "export AZURE_CLIENT_SECRET=$AZURE_CLIENT_SECRET"