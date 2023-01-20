//// Code generated by go-swagger; DO NOT EDIT.

package schemas

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	agentstoragev1 "github.com/loft-sh/agentapi/v2/pkg/apis/loft/storage/v1"
	storagev1 "github.com/loft-sh/api/v2/pkg/apis/storage/v1"
	"github.com/loft-sh/terraform-provider-loft/pkg/utils"
)

func StorageV1VirtualClusterSpaceTemplateDefinitionSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"apps": {
			Type: schema.TypeList,
			Elem: &schema.Resource{
				Schema: StorageV1AppReferenceSchema(),
			},
			Description: "Apps specifies the apps that should get deployed by this template",
			Optional:    true,
		},
		"charts": {
			Type: schema.TypeList,
			Elem: &schema.Resource{
				Schema: StorageV1TemplateHelmChartSchema(),
			},
			Description: "Charts are helm charts that should get deployed",
			Optional:    true,
		},
		"metadata": {
			Type:     schema.TypeList,
			MinItems: 1,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: StorageV1TemplateMetadataSchema(),
			},
			Description: "The space metadata",
			Optional:    true,
		},
		"objects": {
			Type:        schema.TypeString,
			Description: "Objects are Kubernetes style yamls that should get deployed into the virtual cluster namespace",
			Optional:    true,
		},
	}
}

func CreateStorageV1VirtualClusterSpaceTemplateDefinition(in []interface{}) *storagev1.VirtualClusterSpaceTemplateDefinition {
	if !utils.HasValue(in) {
		return nil
	}

	ret := &storagev1.VirtualClusterSpaceTemplateDefinition{}

	data := in[0].(map[string]interface{})

	var appsItems []agentstoragev1.AppReference
	for _, v := range data["apps"].([]interface{}) {
		item := *CreateStorageV1AppReference(v.([]interface{}))
		appsItems = append(appsItems, item)
	}
	ret.Apps = appsItems

	var chartsItems []agentstoragev1.TemplateHelmChart
	for _, v := range data["charts"].([]interface{}) {
		item := *CreateStorageV1TemplateHelmChart(v.([]interface{}))
		chartsItems = append(chartsItems, item)
	}
	ret.Charts = chartsItems

	ret.TemplateMetadata = *CreateStorageV1TemplateMetadata(data["metadata"].([]interface{}))

	if v, ok := data["objects"].(string); ok && len(v) > 0 {
		ret.Objects = v
	}

	return ret
}

func ReadStorageV1VirtualClusterSpaceTemplateDefinition(obj *storagev1.VirtualClusterSpaceTemplateDefinition) (interface{}, error) {
	values := map[string]interface{}{}
	var appsItems []interface{}
	for _, v := range obj.Apps {
		item, err := ReadStorageV1AppReference(&v)
		if err != nil {
			return nil, err
		}
		appsItems = append(appsItems, item)
	}
	values["apps"] = appsItems
	var chartsItems []interface{}
	for _, v := range obj.Charts {
		item, err := ReadStorageV1TemplateHelmChart(&v)
		if err != nil {
			return nil, err
		}
		chartsItems = append(chartsItems, item)
	}
	values["charts"] = chartsItems
	// ComGithubLoftShAPIV3PkgApisStorageV1TemplateMetadata
	// {resolvedType:{IsAnonymous:false IsArray:false IsMap:false IsInterface:false IsPrimitive:false IsCustomFormatter:false IsAliased:false IsNullable:true IsStream:false IsEmptyOmitted:true IsJSONString:false IsEnumCI:false IsBase64:false IsExternal:false IsTuple:false HasAdditionalItems:false IsComplexObject:true IsBaseType:false HasDiscriminator:false GoType:ComGithubLoftShAPIV3PkgApisStorageV1TemplateMetadata Pkg:models PkgAlias: AliasedType: SwaggerType:object SwaggerFormat: Extensions:map[] ElemType:<nil> IsMapNullOverride:false IsSuperAlias:false IsEmbedded:false SkipExternalValidation:false} sharedValidations:{SchemaValidations:{CommonValidations:{Maximum:<nil> ExclusiveMaximum:false Minimum:<nil> ExclusiveMinimum:false MaxLength:<nil> MinLength:<nil> Pattern: MaxItems:<nil> MinItems:<nil> UniqueItems:false MultipleOf:<nil> Enum:[]} PatternProperties:map[] MaxProperties:<nil> MinProperties:<nil>} HasValidations:true HasContextValidations:true Required:false HasSliceValidations:false ItemsEnum:[]} Example: OriginalName:metadata Name:metadata Suffix: Path:"metadata" ValueExpression:m.Metadata IndexVar:i KeyVar: Title: Description:The space metadata Location:body ReceiverName:m Items:<nil> AllowsAdditionalItems:false HasAdditionalItems:false AdditionalItems:<nil> Object:<nil> XMLName: CustomTag: Properties:[] AllOf:[] HasAdditionalProperties:false IsAdditionalProperties:false AdditionalProperties:<nil> StrictAdditionalProperties:false ReadOnly:false IsVirtual:false IsBaseType:false HasBaseType:false IsSubType:false IsExported:true DiscriminatorField: DiscriminatorValue: Discriminates:map[] Parents:[] IncludeValidator:true IncludeModel:true Default:<nil> WantsMarshalBinary:true StructTags:[] ExtraImports:map[] ExternalDocs:<nil>}

	metadata, err := ReadStorageV1TemplateMetadata(&obj.TemplateMetadata)
	if err != nil {
		return nil, err
	}
	values["metadata"] = metadata
	values["objects"] = obj.Objects
	return values, nil
}