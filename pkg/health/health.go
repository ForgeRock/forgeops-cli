package health

import (
	"github.com/ForgeRock/forgeops-cli/internal/k8s"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// Check expression to be evaluated against a resource
type Check struct {
	// https://github.com/antonmedv/expr/blob/master/docs/Language-Definition.md
	Expression string `json:"expression"`
	// duration to test expression
	Timeout metav1.Duration `json:"timeout"`
}

// Resource kubernetes object to have checks run against
// An object has checks evaluated against the object
type Resource struct {
	Group      string   `json:"group,omitempty"`
	APIVersion string   `json:"apiversion"`
	Resource   string   `json:"resource"`
	Name       string   `json:"name"`
	Namespace  string   `json:"namespace"`
	Checks     []*Check `json:"checks"`
}

// Check run wait on resource until expression passes
// note at the moment it doesn't track any success/fails
// a resource is checked in a namespace with the following priority
// 1. Namespace on resource
// 2. if 1. is nil then namespace on hlth
func (r *Resource) Check(clientMgr k8s.ClientMgr, fallBackNamespace string) (bool, error) {
	gvr := schema.GroupVersionResource{
		Group:    r.Group,
		Version:  r.APIVersion,
		Resource: r.Resource,
	}
	namespace := fallBackNamespace
	if r.Namespace != "" {
		namespace = r.Namespace
	}
	passed := true
	for _, check := range r.Checks {
		// TODO WatchEventsForCondition should use a context
		met, err := clientMgr.WatchEventsForCondition(int(check.Timeout.Seconds()), namespace, r.Name, gvr, k8s.ConditionExpression(check.Expression))
		// A watch or condition not being met is failed, but not an _error_
		if errors.Is(err, k8s.ErrWatchTimeout) {
			passed = false

		} else if err != nil {
			return false, err
		}
		if !met {
			passed = false
		}
	}
	return passed, nil
}

// V1AlphaHealthSpec HealthSpec
type V1AlphaHealthSpec struct {
	Resources []*Resource     `json:"resources"`
	Timeout   metav1.Duration `json:"timeout"`
}

// Health is kuberenetes resources that should be checked together as a logical group
type Health struct {
	Spec               V1AlphaHealthSpec `json:"spec"`
	Metadata           metav1.ObjectMeta `json:"metadata"`
	healthy, unhealthy []string
}

// CheckResources wait until all resources checks passed of have been exhausted
// return true if all resources passed checks
func (h *Health) CheckResources(client k8s.ClientMgr, allNamespaces bool) (bool, error) {
	// track reuslts
	var err error = nil
	for _, r := range h.Spec.Resources {
		ns := ""
		if !allNamespaces {
			ns, err = client.Namespace()
			if err != nil {
				err = errors.Wrapf(err, "%s checks failed", r.Name)
			}
		}
		healthy, err := r.Check(client, ns)
		if err != nil {
			//nolint
			err = errors.Wrapf(err, "%s checks failed", r.Name)
		}
		if !healthy {
			h.unhealthy = append(h.unhealthy, r.Name)
			continue
		}
		h.healthy = append(h.healthy, r.Name)
	}
	return len(h.unhealthy) == 0, err
}
