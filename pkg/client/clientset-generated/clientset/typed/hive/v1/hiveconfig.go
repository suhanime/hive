// Code generated by client-gen. DO NOT EDIT.

package v1

import (
	"time"

	v1 "github.com/openshift/hive/pkg/apis/hive/v1"
	scheme "github.com/openshift/hive/pkg/client/clientset-generated/clientset/scheme"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// HiveConfigsGetter has a method to return a HiveConfigInterface.
// A group's client should implement this interface.
type HiveConfigsGetter interface {
	HiveConfigs() HiveConfigInterface
}

// HiveConfigInterface has methods to work with HiveConfig resources.
type HiveConfigInterface interface {
	Create(*v1.HiveConfig) (*v1.HiveConfig, error)
	Update(*v1.HiveConfig) (*v1.HiveConfig, error)
	UpdateStatus(*v1.HiveConfig) (*v1.HiveConfig, error)
	Delete(name string, options *metav1.DeleteOptions) error
	DeleteCollection(options *metav1.DeleteOptions, listOptions metav1.ListOptions) error
	Get(name string, options metav1.GetOptions) (*v1.HiveConfig, error)
	List(opts metav1.ListOptions) (*v1.HiveConfigList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.HiveConfig, err error)
	HiveConfigExpansion
}

// hiveConfigs implements HiveConfigInterface
type hiveConfigs struct {
	client rest.Interface
}

// newHiveConfigs returns a HiveConfigs
func newHiveConfigs(c *HiveV1Client) *hiveConfigs {
	return &hiveConfigs{
		client: c.RESTClient(),
	}
}

// Get takes name of the hiveConfig, and returns the corresponding hiveConfig object, and an error if there is any.
func (c *hiveConfigs) Get(name string, options metav1.GetOptions) (result *v1.HiveConfig, err error) {
	result = &v1.HiveConfig{}
	err = c.client.Get().
		Resource("hiveconfigs").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of HiveConfigs that match those selectors.
func (c *hiveConfigs) List(opts metav1.ListOptions) (result *v1.HiveConfigList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1.HiveConfigList{}
	err = c.client.Get().
		Resource("hiveconfigs").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested hiveConfigs.
func (c *hiveConfigs) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Resource("hiveconfigs").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch()
}

// Create takes the representation of a hiveConfig and creates it.  Returns the server's representation of the hiveConfig, and an error, if there is any.
func (c *hiveConfigs) Create(hiveConfig *v1.HiveConfig) (result *v1.HiveConfig, err error) {
	result = &v1.HiveConfig{}
	err = c.client.Post().
		Resource("hiveconfigs").
		Body(hiveConfig).
		Do().
		Into(result)
	return
}

// Update takes the representation of a hiveConfig and updates it. Returns the server's representation of the hiveConfig, and an error, if there is any.
func (c *hiveConfigs) Update(hiveConfig *v1.HiveConfig) (result *v1.HiveConfig, err error) {
	result = &v1.HiveConfig{}
	err = c.client.Put().
		Resource("hiveconfigs").
		Name(hiveConfig.Name).
		Body(hiveConfig).
		Do().
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().

func (c *hiveConfigs) UpdateStatus(hiveConfig *v1.HiveConfig) (result *v1.HiveConfig, err error) {
	result = &v1.HiveConfig{}
	err = c.client.Put().
		Resource("hiveconfigs").
		Name(hiveConfig.Name).
		SubResource("status").
		Body(hiveConfig).
		Do().
		Into(result)
	return
}

// Delete takes name of the hiveConfig and deletes it. Returns an error if one occurs.
func (c *hiveConfigs) Delete(name string, options *metav1.DeleteOptions) error {
	return c.client.Delete().
		Resource("hiveconfigs").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *hiveConfigs) DeleteCollection(options *metav1.DeleteOptions, listOptions metav1.ListOptions) error {
	var timeout time.Duration
	if listOptions.TimeoutSeconds != nil {
		timeout = time.Duration(*listOptions.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Resource("hiveconfigs").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Timeout(timeout).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched hiveConfig.
func (c *hiveConfigs) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.HiveConfig, err error) {
	result = &v1.HiveConfig{}
	err = c.client.Patch(pt).
		Resource("hiveconfigs").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
