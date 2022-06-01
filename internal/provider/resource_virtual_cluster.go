package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	v1 "github.com/loft-sh/agentapi/v2/pkg/apis/loft/storage/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	agentv1 "github.com/loft-sh/agentapi/v2/pkg/apis/loft/cluster/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func resourceVirtualCluster() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "A Loft VirtualCluster.",

		CreateContext: resourceVirtualClusterCreate,
		ReadContext:   resourceVirtualClusterRead,
		UpdateContext: resourceVirtualClusterUpdate,
		DeleteContext: resourceVirtualClusterDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: virtualClusterAttributes(),
	}
}

func resourceVirtualClusterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient, ok := meta.(*apiClient)
	if !ok {
		return diag.Errorf("Could not access apiClient")
	}

	clusterName := d.Get("cluster").(string)
	virtualClusterName := d.Get("name").(string)
	namespace := d.Get("namespace").(string)

	chartName := d.Get("chart_name").(string)
	chartVersion := d.Get("chart_version").(string)
	values := d.Get("values").(string)

	clusterClient, err := apiClient.LoftClient.Cluster(clusterName)

	if err != nil {
		return diag.FromErr(err)
	}

	virtualCluster := &agentv1.VirtualCluster{
		Spec: agentv1.VirtualClusterSpec{},
	}
	virtualClusterHelmRelease := v1.VirtualClusterHelmRelease{
		Chart: v1.VirtualClusterHelmChart{
			Name:    chartName,
			Version: chartVersion,
		},
		Values: values,
	}
	virtualCluster.Spec.VirtualClusterSpec.HelmRelease = &virtualClusterHelmRelease

	virtualCluster.SetName(virtualClusterName)

	rawAnnotations := d.Get("annotations").(map[string]interface{})
	annotations, err := attributesToMap(rawAnnotations)
	if err != nil {
		return diag.FromErr(err)
	}

	virtualCluster.SetAnnotations(annotations)

	rawLabels := d.Get("labels").(map[string]interface{})
	labels, err := attributesToMap(rawLabels)
	if err != nil {
		return diag.FromErr(err)
	}

	virtualCluster.SetLabels(labels)

	objects := d.Get("objects").(string)
	if objects != "" {
		virtualCluster.Spec.Objects = objects
	}

	virtualCluster, err = clusterClient.Agent().ClusterV1().VirtualClusters(namespace).Create(ctx, virtualCluster, metav1.CreateOptions{})
	if err != nil {
		return diag.FromErr(err)
	}

	if err := readVirtualCluster(clusterName, namespace, virtualCluster, d); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceVirtualClusterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient, ok := meta.(*apiClient)
	if !ok {
		return diag.Errorf("Could not access apiClient")
	}

	clusterName, namespace, virtualClusterName := parseVirtualClusterId(d.Id())
	clusterClient, err := apiClient.LoftClient.Cluster(clusterName)
	if err != nil {
		return diag.FromErr(err)
	}

	virtualCluster, err := clusterClient.Agent().ClusterV1().VirtualClusters(namespace).Get(ctx, virtualClusterName, metav1.GetOptions{})
	if err != nil {
		return diag.FromErr(err)
	}

	if err := readVirtualCluster(clusterName, namespace, virtualCluster, d); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceVirtualClusterUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// use the meta value to retrieve your client from the provider configure method
	apiClient, ok := meta.(*apiClient)
	if !ok {
		return diag.Errorf("Could not access apiClient")
	}

	clusterName, namespace, virtualClusterName := parseVirtualClusterId(d.Id())
	clusterClient, err := apiClient.LoftClient.Cluster(clusterName)
	if err != nil {
		return diag.FromErr(err)
	}

	oldVirtualCluster, err := clusterClient.Agent().ClusterV1().VirtualClusters(namespace).Get(ctx, virtualClusterName, metav1.GetOptions{})
	if err != nil {
		return diag.FromErr(err)
	}

	modifiedVirtualCluster := oldVirtualCluster.DeepCopy()

	if d.HasChange("objects") {
		_, newObjects := d.GetChange("objects")
		modifiedVirtualCluster.Spec.Objects = newObjects.(string)
	}

	if d.HasChange("annotations") {
		oldAnnotations, newAnnotations := d.GetChange("annotations")

		added, modified, deleted, err := getAddedModifiedAndDeleted(
			oldAnnotations.(map[string]interface{}),
			newAnnotations.(map[string]interface{}),
		)

		if err != nil {
			return diag.FromErr(err)
		}

		for k, v := range added {
			modifiedVirtualCluster.Annotations[k] = v.(string)
		}

		for k, v := range modified {
			modifiedVirtualCluster.Annotations[k] = v.(string)
		}

		for k := range deleted {
			delete(modifiedVirtualCluster.Annotations, k)
		}
	}

	if d.HasChange("labels") {
		oldLabels, newLabels := d.GetChange("labels")

		added, modified, deleted, err := getAddedModifiedAndDeleted(
			oldLabels.(map[string]interface{}),
			newLabels.(map[string]interface{}),
		)
		if err != nil {
			return diag.FromErr(err)
		}

		for k, v := range added {
			modifiedVirtualCluster.Labels[k] = v.(string)
		}

		for k, v := range modified {
			modifiedVirtualCluster.Labels[k] = v.(string)
		}

		for k := range deleted {
			delete(modifiedVirtualCluster.Labels, k)
		}
	}

	patch := client.MergeFrom(oldVirtualCluster)
	rawPatch, err := patch.Data(modifiedVirtualCluster)
	if err != nil {
		return diag.FromErr(err)
	}

	virtualCluster, err := clusterClient.Agent().ClusterV1().VirtualClusters(namespace).Patch(ctx, virtualClusterName, patch.Type(), rawPatch, metav1.PatchOptions{})
	if err != nil {
		return diag.FromErr(err)
	}

	if err := readVirtualCluster(clusterName, namespace, virtualCluster, d); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceVirtualClusterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient, ok := meta.(*apiClient)
	if !ok {
		return diag.Errorf("Could not access apiClient")
	}

	clusterName := d.Get("cluster").(string)
	namespace := d.Get("namespace").(string)
	virtualClusterName := d.Get("name").(string)

	clusterClient, err := apiClient.LoftClient.Cluster(clusterName)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := clusterClient.Agent().ClusterV1().VirtualClusters(namespace).Delete(context.TODO(), virtualClusterName, metav1.DeleteOptions{}); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
