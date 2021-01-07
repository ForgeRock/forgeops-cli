package health

import (
	"github.com/antonmedv/expr"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"

	"github.com/ForgeRock/forgeops-cli/internal/k8s"
)

var ErrExpressionResult error = errors.New("only boolean expressions are permitted")

// &metav1.Duration{Duration: 100 * 365 * 24 * time.Hour}
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
	Group      string   `json:"group"`
	APIVersion string   `json:"apiversion"`
	Resource   string   `json:"resource"`
	Name       string   `json:"name"`
	Namespace  string   `json:"namespace"`
	Checks     []*Check `json:"checks"`
}

// conditionExpression wrap condition function using a closure allowing the event call back to have an expression
func conditionExpression(expression string) k8s.ConditionFunction {
	return func(event watch.Event, obj *unstructured.Unstructured) (bool, error) {
		// compile expression with object, Env help with typing during evaluation
		pgrm, err := expr.Compile(expression, expr.Env(obj.Object), expr.AsBool())
		if err != nil {
			return false, errors.WithMessage(ErrExpressionResult, err.Error())
		}
		// check against the object
		output, err := expr.Run(pgrm, obj.Object)
		if err != nil {
			return false, err
		}
		return output.(bool), nil
	}
}

// Check run wait on resource until expression passes
// note at the moment it doesn't track any success/fails
func (r *Resource) Check(clientMgr k8s.ClientMgr, namespace string) (bool, error) {
	gvr := schema.GroupVersionResource{
		Group:    r.Group,
		Version:  r.APIVersion,
		Resource: r.Resource,
	}
	passed := true
	for _, check := range r.Checks {
		// TODO WatchEventsForCondition should use a context
		met, err := clientMgr.WatchEventsForCondition(int(check.Timeout.Seconds()), namespace, r.Name, gvr, conditionExpression(check.Expression))
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
	Name               string            `json:"name"`
	Metadata           metav1.ObjectMeta `json:"metadata"`
	healthy, unhealthy []string
}

// CheckResources wait until all resources checks passed of have been exhausted
// return true if all resources passed checks
func (h *Health) CheckResources(client k8s.ClientMgr) (bool, error) {
	// track reuslts
	var err error
	for _, r := range h.Spec.Resources {
		ns, err := client.Namespace()
		if err != nil {
			return false, err
		}
		healthy, err := r.Check(client, ns)
		if err != nil {
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
