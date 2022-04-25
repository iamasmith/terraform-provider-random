package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceUuid() *schema.Resource {
	return &schema.Resource{
		Description: "The resource `random_uuid` generates random uuid string that is intended to be " +
			"used as unique identifiers for other resources.\n" +
			"\n" +
			"This resource uses [hashicorp/go-uuid](https://github.com/hashicorp/go-uuid) to generate a " +
			"UUID-formatted string for use with services needed a unique string identifier.",
		CreateContext: CreateUuid,
		ReadContext:   schema.NoopContext,
		DeleteContext: RemoveResourceFromState,
		Importer: &schema.ResourceImporter{
			StateContext: ImportUuid,
		},

		Schema: map[string]*schema.Schema{
			"keepers": {
				Description: "Arbitrary map of values that, when changed, will trigger recreation of " +
					"resource. See [the main provider documentation](../index.html) for more information.",
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},

			"result": {
				Description: "The generated uuid presented in string format.",
				Type:        schema.TypeString,
				Computed:    true,
			},

			"id": {
				Description: "The generated uuid presented in string format.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func CreateUuid(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	result, err := uuid.GenerateUUID()
	if err != nil {
		return append(diags, diag.Errorf("error generating uuid: %s", err)...)
	}

	if err := d.Set("result", result); err != nil {
		return append(diags, diag.Errorf("error setting result: %s", err)...)
	}

	d.SetId(result)

	return nil
}

func ImportUuid(_ context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	id := d.Id()

	bytes, err := uuid.ParseUUID(id)
	if err != nil {
		return nil, fmt.Errorf("error parsing uuid bytes: %w", err)
	}

	result, err := uuid.FormatUUID(bytes)
	if err != nil {
		return nil, fmt.Errorf("error formatting uuid bytes: %w", err)
	}

	if err := d.Set("result", result); err != nil {
		return nil, fmt.Errorf("error setting result: %w", err)
	}

	d.SetId(result)

	return []*schema.ResourceData{d}, nil
}
