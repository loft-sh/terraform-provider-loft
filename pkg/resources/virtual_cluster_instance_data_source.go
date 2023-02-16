//// Code generated by go-swagger; DO NOT EDIT.

package resources

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/loft-sh/terraform-provider-loft/pkg/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func VirtualClusterInstanceDataSource() *schema.Resource {
	return &schema.Resource{
		Description: "VirtualClusterInstance holds the VirtualClusterInstance information",
		Schema:      virtualClusterInstanceDataSourceSchema(),
		ReadContext: dataSourceVirtualClusterInstanceRead,
	}
}

func virtualClusterInstanceDataSourceSchema() map[string]*schema.Schema {
	attributes := virtualClusterInstanceAttributes()

	metadataSchema := attributes["metadata"].Elem.(*schema.Resource)

	metadataSchema.Schema["name"].Computed = false
	metadataSchema.Schema["name"].Optional = false
	metadataSchema.Schema["name"].Required = true
	metadataSchema.Schema["name"].ConflictsWith = nil

	metadataSchema.Schema["generate_name"].ConflictsWith = nil
	metadataSchema.Schema["generate_name"].AtLeastOneOf = nil

	metadataSchema.Schema["namespace"].Computed = false
	metadataSchema.Schema["namespace"].Optional = false
	metadataSchema.Schema["namespace"].Required = true

	attributes["spec"].Required = false
	attributes["spec"].Computed = true

	return attributes
}

func dataSourceVirtualClusterInstanceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	metadata := metav1.ObjectMeta{
		Namespace: d.Get("metadata.0.namespace").(string),
		Name:      d.Get("metadata.0.name").(string),
	}
	d.SetId(utils.ReadId(metadata))

	return virtualClusterInstanceRead(ctx, d, meta)
}
