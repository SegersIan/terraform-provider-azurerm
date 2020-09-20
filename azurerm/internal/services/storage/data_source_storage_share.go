package storage

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmStorageShare() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmStorageShareRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"storage_account_name": {
				Type:     schema.TypeString,
				Required: true,
			},

			/*
				"container_access_type": {
					Type:     schema.TypeString,
					Computed: true,
				},

				"url": {
					Type:     schema.TypeString,
					Computed: true,
				},

				"acl": {
					Type:     schema.TypeSet,
					Optional: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"id": {
								Type:         schema.TypeString,
								Required:     true,
								ValidateFunc: validation.StringLenBetween(1, 64),
							},
							"access_policy": {
								Type:     schema.TypeList,
								Optional: true,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"start": {
											Type:         schema.TypeString,
											Required:     true,
											ValidateFunc: validation.StringIsNotEmpty,
										},
										"expiry": {
											Type:         schema.TypeString,
											Required:     true,
											ValidateFunc: validation.StringIsNotEmpty,
										},
										"permissions": {
											Type:         schema.TypeString,
											Required:     true,
											ValidateFunc: validation.StringIsNotEmpty,
										},
									},
								},
							},
						},
					},
				},

				"quota": {
					Type:         schema.TypeInt,
					Optional:     true,
					Default:      5120,
					ValidateFunc: validation.IntBetween(1, 102400),
				},

				"metadata": MetaDataComputedSchema(),

				"resource_manager_id": {
					Type:     schema.TypeString,
					Computed: true,
				},
			*/
		},
	}
}

func dataSourceArmStorageShareRead(d *schema.ResourceData, meta interface{}) error {
	storageClient := meta.(*clients.Client).Storage
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	shareName := d.Get("name").(string)
	accountName := d.Get("storage_account_name").(string)

	account, err := storageClient.FindAccount(ctx, accountName)
	if err != nil {
		return fmt.Errorf("Error retrieving Account %q for Storage Share %q: %s", accountName, shareName, err)
	}
	if account == nil {
		return fmt.Errorf("Unable to locate Account %q for Storage Share %q", accountName, shareName)
	}

	client, err := storageClient.FileSharesClient(ctx, *account)
	if err != nil {
		return fmt.Errorf("Error building Share Client for Storage Account %q (Resource Group %q): %s", accountName, account.ResourceGroup, err)
	}

	d.SetId(client.GetResourceID(accountName, shareName))

	props, err := client.GetProperties(ctx, accountName, shareName)
	if err != nil {
		if utils.ResponseWasNotFound(props.Response) {
			return fmt.Errorf("Share %q was not found in Account %q / Resource Group %q", shareName, accountName, account.ResourceGroup)
		}

		return fmt.Errorf("Error retrieving Share %q (Account %q / Resource Group %q): %s", shareName, accountName, account.ResourceGroup, err)
	}

	d.Set("name", shareName)

	resourceManagerID := client.GetResourceManagerResourceID(storageClient.SubscriptionId, account.ResourceGroup, accountName, shareName)
	d.Set("resource_manager_id", resourceManagerID)

	/*

		d.Set("storage_account_name", accountName)

		d.Set("container_access_type", flattenStorageContainerAccessLevel(props.AccessLevel))

		if err := d.Set("metadata", FlattenMetaData(props.MetaData)); err != nil {
			return fmt.Errorf("Error setting `metadata`: %+v", err)
		}

		d.Set("has_immutability_policy", props.HasImmutabilityPolicy)
		d.Set("has_legal_hold", props.HasLegalHold)


	*/
	return nil
}
