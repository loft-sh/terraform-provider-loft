//// Code generated by go-swagger; DO NOT EDIT.

package schemas

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	agentstoragev1 "github.com/loft-sh/agentapi/v2/pkg/apis/loft/storage/v1"
	managementv1 "github.com/loft-sh/api/v2/pkg/apis/management/v1"
	storagev1 "github.com/loft-sh/api/v2/pkg/apis/storage/v1"
	"github.com/loft-sh/terraform-provider-loft/pkg/utils"
)

func ManagementV1SpaceInstanceSpecSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"access": {
			Type: schema.TypeList,
			Elem: &schema.Resource{
				Schema: StorageV1AccessSchema(),
			},
			Description: "Access holds the access rights for users and teams",
			Optional:    true,
			Computed:    true,
		},
		"cluster_ref": {
			Type:     schema.TypeList,
			MinItems: 1,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: StorageV1ClusterRefSchema(),
			},
			Description: "ClusterRef is the reference to the connected cluster holding this space",
			Optional:    true,
			Computed:    true,
		},
		"description": {
			Type:        schema.TypeString,
			Description: "Description describes a space instance",
			Optional:    true,
		},
		"display_name": {
			Type:        schema.TypeString,
			Description: "DisplayName is the name that should be displayed in the UI",
			Optional:    true,
		},
		"extra_access_rules": {
			Type: schema.TypeList,
			Elem: &schema.Resource{
				Schema: StorageV1InstanceAccessRuleSchema(),
			},
			Description: "ExtraAccessRules defines extra rules which users and teams should have which access to the virtual cluster.",
			Optional:    true,
		},
		"owner": {
			Type:     schema.TypeList,
			MinItems: 1,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: StorageV1UserOrTeamSchema(),
			},
			Description: "Owner holds the owner of this object",
			Optional:    true,
		},
		"parameters": {
			Type:        schema.TypeString,
			Description: "Parameters are values to pass to the template",
			Optional:    true,
		},
		"template": {
			Type:     schema.TypeList,
			MinItems: 1,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: StorageV1SpaceTemplateDefinitionSchema(),
			},
			Description: "Template is the inline template to use for space creation. This is mutually exclusive with templateRef.",
			Optional:    true,
		},
		"template_ref": {
			Type:     schema.TypeList,
			MinItems: 1,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: StorageV1TemplateRefSchema(),
			},
			Description: "TemplateRef holds the space template reference",
			Optional:    true,
		},
	}
}

func CreateManagementV1SpaceInstanceSpec(in []interface{}) *managementv1.SpaceInstanceSpec {
	ret := storagev1.SpaceInstanceSpec{}

	if utils.HasValue(in) {

		data := in[0].(map[string]interface{})

		var accessItems []storagev1.Access
		for _, v := range data["access"].([]interface{}) {
			item := *CreateStorageV1Access(v.([]interface{}))
			accessItems = append(accessItems, item)
		}
		ret.Access = accessItems

		if v, ok := data["cluster_ref"]; ok && len(v.([]interface{})) > 0 {
			ret.ClusterRef = *CreateStorageV1ClusterRef(v.([]interface{}))
		}

		if v, ok := data["description"].(string); ok && len(v) > 0 {
			ret.Description = v
		}

		if v, ok := data["display_name"].(string); ok && len(v) > 0 {
			ret.DisplayName = v
		}

		var extraAccessRulesItems []agentstoragev1.InstanceAccessRule
		for _, v := range data["extra_access_rules"].([]interface{}) {
			item := *CreateStorageV1InstanceAccessRule(v.([]interface{}))
			extraAccessRulesItems = append(extraAccessRulesItems, item)
		}
		ret.ExtraAccessRules = extraAccessRulesItems

		ret.Owner = CreateStorageV1UserOrTeam(data["owner"].([]interface{}))

		if v, ok := data["parameters"].(string); ok && len(v) > 0 {
			ret.Parameters = v
		}

		ret.Template = CreateStorageV1SpaceTemplateDefinition(data["template"].([]interface{}))

		ret.TemplateRef = CreateStorageV1TemplateRef(data["template_ref"].([]interface{}))

	}

	return &managementv1.SpaceInstanceSpec{
		SpaceInstanceSpec: ret,
	}
}

func ReadManagementV1SpaceInstanceSpec(obj *managementv1.SpaceInstanceSpec) (interface{}, error) {
	values := map[string]interface{}{}
	var accessItems []interface{}
	for _, v := range obj.Access {
		item, err := ReadStorageV1Access(&v)
		if err != nil {
			return nil, err
		}
		accessItems = append(accessItems, item)
	}
	values["access"] = accessItems
	// ComGithubLoftShAPIV3PkgApisStorageV1ClusterRef
	// {resolvedType:{IsAnonymous:false IsArray:false IsMap:false IsInterface:false IsPrimitive:false IsCustomFormatter:false IsAliased:false IsNullable:true IsStream:false IsEmptyOmitted:true IsJSONString:false IsEnumCI:false IsBase64:false IsExternal:false IsTuple:false HasAdditionalItems:false IsComplexObject:true IsBaseType:false HasDiscriminator:false GoType:ComGithubLoftShAPIV3PkgApisStorageV1ClusterRef Pkg:models PkgAlias: AliasedType: SwaggerType:object SwaggerFormat: Extensions:map[] ElemType:<nil> IsMapNullOverride:false IsSuperAlias:false IsEmbedded:false SkipExternalValidation:false} sharedValidations:{SchemaValidations:{CommonValidations:{Maximum:<nil> ExclusiveMaximum:false Minimum:<nil> ExclusiveMinimum:false MaxLength:<nil> MinLength:<nil> Pattern: MaxItems:<nil> MinItems:<nil> UniqueItems:false MultipleOf:<nil> Enum:[]} PatternProperties:map[] MaxProperties:<nil> MinProperties:<nil>} HasValidations:true HasContextValidations:true Required:false HasSliceValidations:false ItemsEnum:[]} Example: OriginalName:clusterRef Name:clusterRef Suffix: Path:"clusterRef" ValueExpression:m.ClusterRef IndexVar:i KeyVar: Title: Description:ClusterRef is the reference to the connected cluster holding this space Location:body ReceiverName:m Items:<nil> AllowsAdditionalItems:false HasAdditionalItems:false AdditionalItems:<nil> Object:<nil> XMLName: CustomTag: Properties:[] AllOf:[] HasAdditionalProperties:false IsAdditionalProperties:false AdditionalProperties:<nil> StrictAdditionalProperties:false ReadOnly:false IsVirtual:false IsBaseType:false HasBaseType:false IsSubType:false IsExported:true DiscriminatorField: DiscriminatorValue: Discriminates:map[] Parents:[] IncludeValidator:true IncludeModel:true Default:<nil> WantsMarshalBinary:true StructTags:[] ExtraImports:map[] ExternalDocs:<nil>}

	clusterRef, err := ReadStorageV1ClusterRef(&obj.ClusterRef)
	if err != nil {
		return nil, err
	}
	values["cluster_ref"] = []interface{}{clusterRef}
	values["description"] = obj.Description
	values["display_name"] = obj.DisplayName
	var extraAccessRulesItems []interface{}
	for _, v := range obj.ExtraAccessRules {
		item, err := ReadStorageV1InstanceAccessRule(&v)
		if err != nil {
			return nil, err
		}
		extraAccessRulesItems = append(extraAccessRulesItems, item)
	}
	values["extra_access_rules"] = extraAccessRulesItems
	// ComGithubLoftShAPIV3PkgApisStorageV1UserOrTeam
	// {resolvedType:{IsAnonymous:false IsArray:false IsMap:false IsInterface:false IsPrimitive:false IsCustomFormatter:false IsAliased:false IsNullable:true IsStream:false IsEmptyOmitted:true IsJSONString:false IsEnumCI:false IsBase64:false IsExternal:false IsTuple:false HasAdditionalItems:false IsComplexObject:true IsBaseType:false HasDiscriminator:false GoType:ComGithubLoftShAPIV3PkgApisStorageV1UserOrTeam Pkg:models PkgAlias: AliasedType: SwaggerType:object SwaggerFormat: Extensions:map[] ElemType:<nil> IsMapNullOverride:false IsSuperAlias:false IsEmbedded:false SkipExternalValidation:false} sharedValidations:{SchemaValidations:{CommonValidations:{Maximum:<nil> ExclusiveMaximum:false Minimum:<nil> ExclusiveMinimum:false MaxLength:<nil> MinLength:<nil> Pattern: MaxItems:<nil> MinItems:<nil> UniqueItems:false MultipleOf:<nil> Enum:[]} PatternProperties:map[] MaxProperties:<nil> MinProperties:<nil>} HasValidations:true HasContextValidations:true Required:false HasSliceValidations:false ItemsEnum:[]} Example: OriginalName:owner Name:owner Suffix: Path:"owner" ValueExpression:m.Owner IndexVar:i KeyVar: Title: Description:Owner holds the owner of this object Location:body ReceiverName:m Items:<nil> AllowsAdditionalItems:false HasAdditionalItems:false AdditionalItems:<nil> Object:<nil> XMLName: CustomTag: Properties:[] AllOf:[] HasAdditionalProperties:false IsAdditionalProperties:false AdditionalProperties:<nil> StrictAdditionalProperties:false ReadOnly:false IsVirtual:false IsBaseType:false HasBaseType:false IsSubType:false IsExported:true DiscriminatorField: DiscriminatorValue: Discriminates:map[] Parents:[] IncludeValidator:true IncludeModel:true Default:<nil> WantsMarshalBinary:true StructTags:[] ExtraImports:map[] ExternalDocs:<nil>}

	owner, err := ReadStorageV1UserOrTeam(obj.Owner)
	if err != nil {
		return nil, err
	}
	values["owner"] = []interface{}{owner}
	values["parameters"] = obj.Parameters
	// ComGithubLoftShAPIV3PkgApisStorageV1SpaceTemplateDefinition
	// {resolvedType:{IsAnonymous:false IsArray:false IsMap:false IsInterface:false IsPrimitive:false IsCustomFormatter:false IsAliased:false IsNullable:true IsStream:false IsEmptyOmitted:true IsJSONString:false IsEnumCI:false IsBase64:false IsExternal:false IsTuple:false HasAdditionalItems:false IsComplexObject:true IsBaseType:false HasDiscriminator:false GoType:ComGithubLoftShAPIV3PkgApisStorageV1SpaceTemplateDefinition Pkg:models PkgAlias: AliasedType: SwaggerType:object SwaggerFormat: Extensions:map[] ElemType:<nil> IsMapNullOverride:false IsSuperAlias:false IsEmbedded:false SkipExternalValidation:false} sharedValidations:{SchemaValidations:{CommonValidations:{Maximum:<nil> ExclusiveMaximum:false Minimum:<nil> ExclusiveMinimum:false MaxLength:<nil> MinLength:<nil> Pattern: MaxItems:<nil> MinItems:<nil> UniqueItems:false MultipleOf:<nil> Enum:[]} PatternProperties:map[] MaxProperties:<nil> MinProperties:<nil>} HasValidations:true HasContextValidations:true Required:false HasSliceValidations:false ItemsEnum:[]} Example: OriginalName:template Name:template Suffix: Path:"template" ValueExpression:m.Template IndexVar:i KeyVar: Title: Description:Template is the inline template to use for space creation. This is mutually exclusive with templateRef. Location:body ReceiverName:m Items:<nil> AllowsAdditionalItems:false HasAdditionalItems:false AdditionalItems:<nil> Object:<nil> XMLName: CustomTag: Properties:[] AllOf:[] HasAdditionalProperties:false IsAdditionalProperties:false AdditionalProperties:<nil> StrictAdditionalProperties:false ReadOnly:false IsVirtual:false IsBaseType:false HasBaseType:false IsSubType:false IsExported:true DiscriminatorField: DiscriminatorValue: Discriminates:map[] Parents:[] IncludeValidator:true IncludeModel:true Default:<nil> WantsMarshalBinary:true StructTags:[] ExtraImports:map[] ExternalDocs:<nil>}

	template, err := ReadStorageV1SpaceTemplateDefinition(obj.Template)
	if err != nil {
		return nil, err
	}
	values["template"] = []interface{}{template}
	// ComGithubLoftShAPIV3PkgApisStorageV1TemplateRef
	// {resolvedType:{IsAnonymous:false IsArray:false IsMap:false IsInterface:false IsPrimitive:false IsCustomFormatter:false IsAliased:false IsNullable:true IsStream:false IsEmptyOmitted:true IsJSONString:false IsEnumCI:false IsBase64:false IsExternal:false IsTuple:false HasAdditionalItems:false IsComplexObject:true IsBaseType:false HasDiscriminator:false GoType:ComGithubLoftShAPIV3PkgApisStorageV1TemplateRef Pkg:models PkgAlias: AliasedType: SwaggerType:object SwaggerFormat: Extensions:map[] ElemType:<nil> IsMapNullOverride:false IsSuperAlias:false IsEmbedded:false SkipExternalValidation:false} sharedValidations:{SchemaValidations:{CommonValidations:{Maximum:<nil> ExclusiveMaximum:false Minimum:<nil> ExclusiveMinimum:false MaxLength:<nil> MinLength:<nil> Pattern: MaxItems:<nil> MinItems:<nil> UniqueItems:false MultipleOf:<nil> Enum:[]} PatternProperties:map[] MaxProperties:<nil> MinProperties:<nil>} HasValidations:true HasContextValidations:true Required:false HasSliceValidations:false ItemsEnum:[]} Example: OriginalName:templateRef Name:templateRef Suffix: Path:"templateRef" ValueExpression:m.TemplateRef IndexVar:i KeyVar: Title: Description:TemplateRef holds the space template reference Location:body ReceiverName:m Items:<nil> AllowsAdditionalItems:false HasAdditionalItems:false AdditionalItems:<nil> Object:<nil> XMLName: CustomTag: Properties:[] AllOf:[] HasAdditionalProperties:false IsAdditionalProperties:false AdditionalProperties:<nil> StrictAdditionalProperties:false ReadOnly:false IsVirtual:false IsBaseType:false HasBaseType:false IsSubType:false IsExported:true DiscriminatorField: DiscriminatorValue: Discriminates:map[] Parents:[] IncludeValidator:true IncludeModel:true Default:<nil> WantsMarshalBinary:true StructTags:[] ExtraImports:map[] ExternalDocs:<nil>}

	templateRef, err := ReadStorageV1TemplateRef(obj.TemplateRef)
	if err != nil {
		return nil, err
	}
	values["template_ref"] = templateRef
	return values, nil
}
