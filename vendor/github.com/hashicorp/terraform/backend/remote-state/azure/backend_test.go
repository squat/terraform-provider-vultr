package azure

import (
	"context"
	"os"
	"testing"

	"github.com/hashicorp/terraform/backend"
	"github.com/hashicorp/terraform/helper/acctest"
)

func TestBackend_impl(t *testing.T) {
	var _ backend.Backend = new(Backend)
}

func TestBackendConfig(t *testing.T) {
	// This test just instantiates the client. Shouldn't make any actual
	// requests nor incur any costs.

	config := map[string]interface{}{
		"storage_account_name": "tfaccount",
		"container_name":       "tfcontainer",
		"key":                  "state",
		// Access Key must be Base64
		"access_key": "QUNDRVNTX0tFWQ0K",
	}

	b := backend.TestBackendConfig(t, New(), backend.TestWrapConfig(config)).(*Backend)

	if b.containerName != "tfcontainer" {
		t.Fatalf("Incorrect bucketName was populated")
	}
	if b.keyName != "state" {
		t.Fatalf("Incorrect keyName was populated")
	}
}

func TestBackendAccessKeyBasic(t *testing.T) {
	testAccAzureBackend(t)
	rs := acctest.RandString(4)
	res := testResourceNames(rs, "testState")
	armClient := buildTestClient(t, res)

	ctx := context.TODO()
	err := armClient.buildTestResources(ctx, &res)
	defer armClient.destroyTestResources(ctx, res)
	if err != nil {
		armClient.destroyTestResources(ctx, res)
		t.Fatalf("Error creating Test Resources: %q", err)
	}

	b := backend.TestBackendConfig(t, New(), backend.TestWrapConfig(map[string]interface{}{
		"storage_account_name": res.storageAccountName,
		"container_name":       res.storageContainerName,
		"key":                  res.storageKeyName,
		"access_key":           res.storageAccountAccessKey,
		"environment":          os.Getenv("ARM_ENVIRONMENT"),
		"endpoint":             os.Getenv("ARM_ENDPOINT"),
	})).(*Backend)

	backend.TestBackendStates(t, b)
}

func TestBackendManagedServiceIdentityBasic(t *testing.T) {
	testAccAzureBackendRunningInAzure(t)
	rs := acctest.RandString(4)
	res := testResourceNames(rs, "testState")
	armClient := buildTestClient(t, res)

	ctx := context.TODO()
	err := armClient.buildTestResources(ctx, &res)
	defer armClient.destroyTestResources(ctx, res)
	if err != nil {
		t.Fatalf("Error creating Test Resources: %q", err)
	}

	b := backend.TestBackendConfig(t, New(), backend.TestWrapConfig(map[string]interface{}{
		"storage_account_name": res.storageAccountName,
		"container_name":       res.storageContainerName,
		"key":                  res.storageKeyName,
		"resource_group_name":  res.resourceGroup,
		"use_msi":              true,
		"subscription_id":      os.Getenv("ARM_SUBSCRIPTION_ID"),
		"tenant_id":            os.Getenv("ARM_TENANT_ID"),
		"environment":          os.Getenv("ARM_ENVIRONMENT"),
		"endpoint":             os.Getenv("ARM_ENDPOINT"),
	})).(*Backend)

	backend.TestBackendStates(t, b)
}

func TestBackendSASTokenBasic(t *testing.T) {
	testAccAzureBackend(t)
	rs := acctest.RandString(4)
	res := testResourceNames(rs, "testState")
	armClient := buildTestClient(t, res)

	ctx := context.TODO()
	err := armClient.buildTestResources(ctx, &res)
	defer armClient.destroyTestResources(ctx, res)
	if err != nil {
		t.Fatalf("Error creating Test Resources: %q", err)
	}

	sasToken, err := buildSasToken(res.storageAccountName, res.storageAccountAccessKey)
	if err != nil {
		t.Fatalf("Error building SAS Token: %+v", err)
	}

	b := backend.TestBackendConfig(t, New(), backend.TestWrapConfig(map[string]interface{}{
		"storage_account_name": res.storageAccountName,
		"container_name":       res.storageContainerName,
		"key":                  res.storageKeyName,
		"sas_token":            *sasToken,
		"environment":          os.Getenv("ARM_ENVIRONMENT"),
		"endpoint":             os.Getenv("ARM_ENDPOINT"),
	})).(*Backend)

	backend.TestBackendStates(t, b)
}

func TestBackendServicePrincipalBasic(t *testing.T) {
	testAccAzureBackend(t)
	rs := acctest.RandString(4)
	res := testResourceNames(rs, "testState")
	armClient := buildTestClient(t, res)

	ctx := context.TODO()
	err := armClient.buildTestResources(ctx, &res)
	defer armClient.destroyTestResources(ctx, res)
	if err != nil {
		t.Fatalf("Error creating Test Resources: %q", err)
	}

	b := backend.TestBackendConfig(t, New(), backend.TestWrapConfig(map[string]interface{}{
		"storage_account_name": res.storageAccountName,
		"container_name":       res.storageContainerName,
		"key":                  res.storageKeyName,
		"resource_group_name":  res.resourceGroup,
		"subscription_id":      os.Getenv("ARM_SUBSCRIPTION_ID"),
		"tenant_id":            os.Getenv("ARM_TENANT_ID"),
		"client_id":            os.Getenv("ARM_CLIENT_ID"),
		"client_secret":        os.Getenv("ARM_CLIENT_SECRET"),
		"environment":          os.Getenv("ARM_ENVIRONMENT"),
		"endpoint":             os.Getenv("ARM_ENDPOINT"),
	})).(*Backend)

	backend.TestBackendStates(t, b)
}

func TestBackendServicePrincipalCustomEndpoint(t *testing.T) {
	testAccAzureBackend(t)

	// this is only applicable for Azure Stack.
	endpoint := os.Getenv("ARM_ENDPOINT")
	if endpoint == "" {
		t.Skip("Skipping as ARM_ENDPOINT isn't configured")
	}

	rs := acctest.RandString(4)
	res := testResourceNames(rs, "testState")
	armClient := buildTestClient(t, res)

	ctx := context.TODO()
	err := armClient.buildTestResources(ctx, &res)
	defer armClient.destroyTestResources(ctx, res)
	if err != nil {
		t.Fatalf("Error creating Test Resources: %q", err)
	}

	b := backend.TestBackendConfig(t, New(), backend.TestWrapConfig(map[string]interface{}{
		"storage_account_name": res.storageAccountName,
		"container_name":       res.storageContainerName,
		"key":                  res.storageKeyName,
		"resource_group_name":  res.resourceGroup,
		"subscription_id":      os.Getenv("ARM_SUBSCRIPTION_ID"),
		"tenant_id":            os.Getenv("ARM_TENANT_ID"),
		"client_id":            os.Getenv("ARM_CLIENT_ID"),
		"client_secret":        os.Getenv("ARM_CLIENT_SECRET"),
		"environment":          os.Getenv("ARM_ENVIRONMENT"),
		"endpoint":             endpoint,
	})).(*Backend)

	backend.TestBackendStates(t, b)
}

func TestBackendAccessKeyLocked(t *testing.T) {
	testAccAzureBackend(t)
	rs := acctest.RandString(4)
	res := testResourceNames(rs, "testState")
	armClient := buildTestClient(t, res)

	ctx := context.TODO()
	err := armClient.buildTestResources(ctx, &res)
	defer armClient.destroyTestResources(ctx, res)
	if err != nil {
		t.Fatalf("Error creating Test Resources: %q", err)
	}

	b1 := backend.TestBackendConfig(t, New(), backend.TestWrapConfig(map[string]interface{}{
		"storage_account_name": res.storageAccountName,
		"container_name":       res.storageContainerName,
		"key":                  res.storageKeyName,
		"access_key":           res.storageAccountAccessKey,
		"environment":          os.Getenv("ARM_ENVIRONMENT"),
		"endpoint":             os.Getenv("ARM_ENDPOINT"),
	})).(*Backend)

	b2 := backend.TestBackendConfig(t, New(), backend.TestWrapConfig(map[string]interface{}{
		"storage_account_name": res.storageAccountName,
		"container_name":       res.storageContainerName,
		"key":                  res.storageKeyName,
		"access_key":           res.storageAccountAccessKey,
		"environment":          os.Getenv("ARM_ENVIRONMENT"),
		"endpoint":             os.Getenv("ARM_ENDPOINT"),
	})).(*Backend)

	backend.TestBackendStateLocks(t, b1, b2)
	backend.TestBackendStateForceUnlock(t, b1, b2)
}

func TestBackendServicePrincipalLocked(t *testing.T) {
	testAccAzureBackend(t)
	rs := acctest.RandString(4)
	res := testResourceNames(rs, "testState")
	armClient := buildTestClient(t, res)

	ctx := context.TODO()
	err := armClient.buildTestResources(ctx, &res)
	defer armClient.destroyTestResources(ctx, res)
	if err != nil {
		t.Fatalf("Error creating Test Resources: %q", err)
	}

	b1 := backend.TestBackendConfig(t, New(), backend.TestWrapConfig(map[string]interface{}{
		"storage_account_name": res.storageAccountName,
		"container_name":       res.storageContainerName,
		"key":                  res.storageKeyName,
		"access_key":           res.storageAccountAccessKey,
		"subscription_id":      os.Getenv("ARM_SUBSCRIPTION_ID"),
		"tenant_id":            os.Getenv("ARM_TENANT_ID"),
		"client_id":            os.Getenv("ARM_CLIENT_ID"),
		"client_secret":        os.Getenv("ARM_CLIENT_SECRET"),
		"environment":          os.Getenv("ARM_ENVIRONMENT"),
		"endpoint":             os.Getenv("ARM_ENDPOINT"),
	})).(*Backend)

	b2 := backend.TestBackendConfig(t, New(), backend.TestWrapConfig(map[string]interface{}{
		"storage_account_name": res.storageAccountName,
		"container_name":       res.storageContainerName,
		"key":                  res.storageKeyName,
		"access_key":           res.storageAccountAccessKey,
		"subscription_id":      os.Getenv("ARM_SUBSCRIPTION_ID"),
		"tenant_id":            os.Getenv("ARM_TENANT_ID"),
		"client_id":            os.Getenv("ARM_CLIENT_ID"),
		"client_secret":        os.Getenv("ARM_CLIENT_SECRET"),
		"environment":          os.Getenv("ARM_ENVIRONMENT"),
		"endpoint":             os.Getenv("ARM_ENDPOINT"),
	})).(*Backend)

	backend.TestBackendStateLocks(t, b1, b2)
	backend.TestBackendStateForceUnlock(t, b1, b2)
}
