package tests

import (
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceArmStorageShare_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_storage_share", "test")
}
