#!/usr/bin/env bash
az login --identity --username $MSI_ID > /dev/null
export ARM_SUBSCRIPTION_ID=$(az login --identity --username $MSI_ID | jq -r '.[0] | .id')
export ARM_TENANT_ID=$(az login --identity --username $MSI_ID | jq -r '.[0] | .tenantId')
export ARM_USE_MSI=true
terraform init -backend-config="storage_account_name=$BACKEND_STORAGE_ACCOUNT_NAME" -backend-config="resource_group_name=$BACKEND_RESOURCE_GROUP_NAME" -backend-config="container_name=$BACKEND_CONTAINER_NAME" -backend-config="key=$BACKEND_KEY"
terraform apply -auto-approve -input=false
