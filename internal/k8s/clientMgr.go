package k8s

import (
	"fmt"
	"io"

	"github.com/ForgeRock/forgeops-cli/internal/factory"
	"github.com/ForgeRock/forgeops-cli/internal/printer"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/resource"
)

// NewK8sClientMgr create a new instance of NewK8sClientMgr
func NewK8sClientMgr(f factory.Factory) ClientMgr {
	return &clientMgr{
		f,
	}
}

// ClientMgr a container for creating kubernetes rest clients
type ClientMgr interface {
	factory.Factory
	Namespace() (string, error)
	GetObjectsFromPath(path string) ([]*resource.Info, error)
	GetObjectsFromStream(reader io.Reader) ([]*resource.Info, error)
	GetObjectsFromServer(resourceType, name string) ([]*resource.Info, error)
	ApplyObject(info *resource.Info) error
	DeleteObject(info *resource.Info) error
	WaitForResource(timeoutSecs int, ns, name string, gvr schema.GroupVersionResource) (bool, error)
	WaitForResourceStatusCondition(timeoutSecs int, ns, name, conditionStr string, gvr schema.GroupVersionResource) (bool, error)
	WaitForResourceReplicas(timeoutSecs int, ns, name, replicas string, gvr schema.GroupVersionResource) (bool, error)
}

type clientMgr struct {
	factory.Factory
}

// NullSchema always validates bytes.
type NullSchema struct{}

// ValidateBytes never fails for NullSchema.
func (NullSchema) ValidateBytes(data []byte) error { return nil }

func (cmgr clientMgr) Namespace() (string, error) {
	cfg, err := cmgr.GetOverrideFlags()
	if err != nil {
		return "", err
	}
	// If no ns is provided in the flags, use the default kubeconfig value.
	ns, _, err := cfg.ToRawKubeConfigLoader().Namespace()
	return ns, nil
}

// GetObjectsFromPath Obtains objects from filepath or url
func (cmgr clientMgr) GetObjectsFromPath(path string) ([]*resource.Info, error) {
	usage := "contains the manifest to process"
	cfg, err := cmgr.GetOverrideFlags()
	if err != nil {
		return nil, err
	}
	filenames := []string{path}
	kustomize := ""
	recursive := false
	fileNameFlags := &genericclioptions.FileNameFlags{
		Usage:     usage,
		Filenames: &filenames,
		Kustomize: &kustomize,
		Recursive: &recursive,
	}
	fileNameOpts := fileNameFlags.ToOptions()
	builder := cmgr.Builder()
	r := builder.
		Unstructured().
		Schema(NullSchema{}).
		ContinueOnError().
		// Use cfg.Namespace in case there's an override. Otherwise default to the ns in the manifest
		NamespaceParam(*cfg.Namespace).DefaultNamespace().
		FilenameParam(false, &fileNameOpts).
		Flatten().
		Do()
	objects, err := r.Infos()
	return objects, err
}

// GetObjectsFromPath Obtains objects from a io.Reader stream
func (cmgr clientMgr) GetObjectsFromStream(reader io.Reader) ([]*resource.Info, error) {
	cfg, err := cmgr.GetOverrideFlags()
	if err != nil {
		return nil, err
	}

	builder := cmgr.Builder()
	r := builder.
		Unstructured().
		Schema(NullSchema{}).
		ContinueOnError().
		// Use cfg.Namespace in case there's an override. Otherwise default to the ns in the manifest
		NamespaceParam(*cfg.Namespace).DefaultNamespace().
		Stream(reader, "stream").
		Flatten().
		Do()
	objects, err := r.Infos()
	return objects, err
}

// if no name is provided, this function will return all objects of the given type
func (cmgr clientMgr) GetObjectsFromServer(resourceType, name string) ([]*resource.Info, error) {
	ns, err := cmgr.Namespace()
	if err != nil {
		return nil, err
	}
	selectAll := true
	nameSelector := ""
	if len(name) > 0 {
		selectAll = false
		nameSelector = fields.OneTermEqualSelector("metadata.name", name).String()
	}
	builder := cmgr.Builder()
	r := builder.
		Unstructured().
		ContinueOnError().
		NamespaceParam(ns).DefaultNamespace().
		SelectAllParam(selectAll).
		FieldSelectorParam(nameSelector).
		SingleResourceType().
		ResourceTypes(resourceType).
		Flatten().
		Do()
	objects, err := r.Infos()
	return objects, err
}

// ApplyObject Applies changes to the object
func (cmgr clientMgr) ApplyObject(info *resource.Info) error {
	helper := resource.NewHelper(info.Client, info.Mapping).
		WithFieldManager("forgeops-cli")

	// Clear "managedFields" before patching
	unstructured.RemoveNestedField(info.Object.(*unstructured.Unstructured).Object, "metadata", "managedFields")
	data, err := runtime.Encode(unstructured.UnstructuredJSONScheme, info.Object)
	if err != nil {
		if statusError, ok := err.(apierrors.APIStatus); ok {
			status := statusError.Status()
			status.Message = fmt.Sprintf("error when %s %q: %v", "serverside-apply", info.Source, status.Message)
			return &apierrors.StatusError{ErrStatus: status}
		}
		return fmt.Errorf("error when %s %q: %v", "serverside-apply", info.Source, err)
	}

	// if the object doesn't exist. Create it. Otherwise, patch it.
	if err := info.Get(); err != nil {
		if apierrors.IsNotFound(err) {
			obj, err := helper.Create(info.Namespace, true, info.Object)
			if err != nil {
				return err
			}
			info.Refresh(obj, true)
			printer.Noticef(fmt.Sprintf("%s %q created", info.ResourceMapping().GroupVersionKind.Kind, info.Name))
			return nil
		}
		return err
	}
	// The object exist. Patch/Apply instead
	// Send the full object to be applied on the server side.
	obj, err := helper.Patch(
		info.Namespace,
		info.Name,
		types.MergePatchType,
		data,
		nil,
	)
	if err != nil {
		return err
	}
	info.Refresh(obj, true)
	printer.Noticef(fmt.Sprintf("%s %q applied", info.ResourceMapping().GroupVersionKind.Kind, info.Name))
	return nil
}

func (cmgr clientMgr) DeleteObject(info *resource.Info) error {
	helper := resource.NewHelper(info.Client, info.Mapping).
		WithFieldManager("forgeops-cli")

	var gracePeriodSeconds int64 = 30
	propagationPolicy := metav1.DeletePropagationBackground
	options := metav1.DeleteOptions{
		GracePeriodSeconds: &gracePeriodSeconds,
		PropagationPolicy:  &propagationPolicy,
	}
	obj, err := helper.DeleteWithOptions(info.Namespace, info.Name, &options)
	if err != nil {
		return err
	}
	info.Refresh(obj, true)
	printer.Noticef(fmt.Sprintf("%s %q deleted", info.ResourceMapping().GroupVersionKind.Kind, info.Name))
	return nil
}
