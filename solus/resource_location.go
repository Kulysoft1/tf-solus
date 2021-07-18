package solus

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/solusio/solus-go-sdk"
)

func resourceLocation() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLocationCreate,
		ReadContext:   resourceLocationRead,
		UpdateContext: resourceLocationUpdate,
		DeleteContext: resourceLocationDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"description": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"icon_id": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"is_default": {
				Type:         schema.TypeBool,
				Optional:     true,
				Default:      false,
				ValidateFunc: validation.NoZeroValues,
			},
			"is_visible": {
				Type:         schema.TypeBool,
				Optional:     true,
				Default:      true,
				ValidateFunc: validation.NoZeroValues,
			},
		},
	}
}

func resourceLocationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	metadata, ok := m.(metadata)
	if !ok {
		return diag.Errorf("invalid metadata type %T", m)
	}
	client := metadata.Client
	timeout := metadata.RequestTimeout

	name := d.Get("name").(string)
	description := d.Get("description").(string)
	iconID := d.Get("icon_id").(int)
	isDefault := d.Get("is_default").(bool)
	isVisible := d.Get("is_visible").(bool)

	reqCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	l, err := client.Locations.Create(reqCtx, solus.LocationCreateRequest{
		Name:        name,
		Description: description,
		IconID:      newNullableIntForID(iconID),
		IsDefault:   isDefault,
		IsVisible:   isVisible,
	})
	if err != nil {
		return diag.Errorf("failed to create new location: %s", err)
	}

	d.SetId(strconv.Itoa(l.ID))
	return resourceLocationRead(ctx, d, m)
}

func resourceLocationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	metadata, ok := m.(metadata)
	if !ok {
		return diag.Errorf("invalid metadata type %T", m)
	}
	client := metadata.Client
	timeout := metadata.RequestTimeout

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	reqCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	l, err := client.Locations.Get(reqCtx, id)
	if err != nil {
		return diag.Errorf("failed to get location by id %d: %s", id, err)
	}

	err = (&schemaChainSetter{d: d}).
		SetID(l.ID).
		Set("name", l.Name).
		Set("icon_id", l.Icon.ID).
		Set("description", l.Description).
		Set("is_default", l.IsDefault).
		Set("is_visible", l.IsVisible).
		Error()
	if err != nil {
		return diag.Errorf("failed to map location response to resource: %s", err)
	}

	return nil
}

func resourceLocationUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	metadata, ok := m.(metadata)
	if !ok {
		return diag.Errorf("invalid metadata type %T", m)
	}
	client := metadata.Client
	timeout := metadata.RequestTimeout

	name := d.Get("name").(string)
	description := d.Get("description").(string)
	iconID := d.Get("icon_id").(int)
	isDefault := d.Get("is_default").(bool)
	isVisible := d.Get("is_visible").(bool)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	reqCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	l, err := client.Locations.Update(reqCtx, id, solus.LocationCreateRequest{
		Name:        name,
		Description: description,
		IconID:      newNullableIntForID(iconID),
		IsDefault:   isDefault,
		IsVisible:   isVisible,
	})
	if err != nil {
		return diag.Errorf("failed to update location with id %d: %s", id, err)
	}

	d.SetId(strconv.Itoa(l.ID))
	return resourceLocationRead(ctx, d, m)
}

func resourceLocationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	metadata, ok := m.(metadata)
	if !ok {
		return diag.Errorf("invalid metadata type %T", m)
	}
	client := metadata.Client
	timeout := metadata.RequestTimeout

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	reqCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	err = client.Locations.Delete(reqCtx, id)
	if err != nil {
		return diag.Errorf("failed to delete location by id %d: %s", id, err)
	}
	return nil
}
