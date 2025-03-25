package provider

import (
	"context"
	"errors"

	"github.com/Broadcom/terraform-provider-luminate/service"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type CreateLuminateContextFunc func(context.Context, *schema.ResourceData, *service.LuminateService) diag.Diagnostics

type ReadLuminateContextFunc func(context.Context, *schema.ResourceData, *service.LuminateService) diag.Diagnostics

type UpdateLuminateContextFunc func(context.Context, *schema.ResourceData, *service.LuminateService) diag.Diagnostics

type DeleteLuminateContextFunc func(context.Context, *schema.ResourceData, *service.LuminateService) diag.Diagnostics

func createLuminateContext(f CreateLuminateContextFunc) schema.CreateContextFunc {
	return func(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
		client, ok := m.(*service.LuminateService)
		if !ok {
			return diag.FromErr(errors.New("unable to cast Luminate service"))
		}

		return f(ctx, data, client)
	}
}

func readLuminateContext(f ReadLuminateContextFunc) schema.ReadContextFunc {
	return func(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
		client, ok := m.(*service.LuminateService)
		if !ok {
			return diag.FromErr(errors.New("unable to cast Luminate service"))
		}

		return f(ctx, data, client)
	}
}

func updateLuminateContext(f UpdateLuminateContextFunc) schema.UpdateContextFunc {
	return func(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
		client, ok := m.(*service.LuminateService)
		if !ok {
			return diag.FromErr(errors.New("unable to cast Luminate service"))
		}

		return f(ctx, data, client)
	}
}

func deleteLuminateContext(f DeleteLuminateContextFunc) schema.DeleteContextFunc {
	return func(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
		client, ok := m.(*service.LuminateService)
		if !ok {
			return diag.FromErr(errors.New("unable to cast Luminate service"))
		}

		return f(ctx, data, client)
	}
}
