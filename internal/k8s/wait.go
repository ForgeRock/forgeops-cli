package k8s

import (
	"context"
	"strings"
	"time"

	"github.com/pkg/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	watchtools "k8s.io/client-go/tools/watch"
)

var ErrWatchTimeout error = errors.New("wait timeout")

type ConditionFunction func(event watch.Event, obj *unstructured.Unstructured) (bool, error)

// WatchEventsForCondition sets a watch for events and evaluatest the provided ConditionFunction
func (cmgr clientMgr) WatchEventsForCondition(timeoutSecs int, ns, name string, gvr schema.GroupVersionResource, condition ConditionFunction) (bool, error) {
	for {
		endTime := time.Now().Add(time.Duration(timeoutSecs) * time.Second)
		dynamicClient, err := cmgr.DynamicClient()
		if err != nil {
			return false, err
		}
		nameSelector := fields.OneTermEqualSelector("metadata.name", name).String()
		gottenObjList, err := dynamicClient.Resource(gvr).Namespace(ns).List(context.TODO(), metav1.ListOptions{FieldSelector: nameSelector})
		if apierrors.IsNotFound(err) {
			gottenObjList = &unstructured.UnstructuredList{}
		} else if err != nil && !apierrors.IsNotFound(err) {
			return false, err
		}
		// If the object is present, let's evaluate if the condition has already been met.
		if len(gottenObjList.Items) != 0 {
			if ok, err := condition(watch.Event{}, &gottenObjList.Items[0]); ok && err == nil {
				return true, nil
			}
		}
		// The condition has not been met. Set a watch on the object
		watchOptions := metav1.ListOptions{}
		watchOptions.FieldSelector = nameSelector
		watchOptions.ResourceVersion = gottenObjList.GetResourceVersion()
		objWatch, err := dynamicClient.Resource(gvr).Namespace(ns).Watch(context.TODO(), watchOptions)
		if err != nil {
			return false, err
		}
		timeout := endTime.Sub(time.Now())
		if timeout < 0 {
			// we're out of time
			return false, errors.WithMessagef(ErrWatchTimeout, "timedout on %s/%s", gvr, name)
		}
		isConditionMet := func(event watch.Event) (bool, error) {
			if event.Type == watch.Error {
				err := apierrors.FromObject(event.Object)
				return false, err
			}
			obj := event.Object.(*unstructured.Unstructured)
			return condition(event, obj)
		}
		ctx, cancel := watchtools.ContextWithOptionalTimeout(context.Background(), timeout)
		lastEvent, err := watchtools.UntilWithoutRetry(ctx, objWatch, isConditionMet)
		cancel()
		if err == watchtools.ErrWatchClosed {
			continue
		}
		if err != nil || lastEvent == nil {
			return false, err
		}
		return true, nil
	}
}

// WaitForResource waits until a resource is present in the k8s API
func (cmgr clientMgr) WaitForResource(timeoutSecs int, ns, name string, gvr schema.GroupVersionResource) (bool, error) {
	var condFunc ConditionFunction = func(event watch.Event, obj *unstructured.Unstructured) (bool, error) {
		if event.Type == watch.Deleted {
			return false, nil
		}
		return true, nil
	}
	return cmgr.WatchEventsForCondition(timeoutSecs, ns, name, gvr, condFunc)

}

// WaitForResourceStatusCondition waits until a resource is present in the k8s API with a given status.conditions
func (cmgr clientMgr) WaitForResourceStatusCondition(timeoutSecs int, ns, name, conditionStr string, gvr schema.GroupVersionResource) (bool, error) {
	var condFunc ConditionFunction = func(event watch.Event, obj *unstructured.Unstructured) (bool, error) {
		conditions, found, err := unstructured.NestedSlice(obj.Object, "status", "conditions")
		if err != nil {
			return false, err
		}
		if !found {
			return false, nil
		}
		for _, conditionUncast := range conditions {
			condition := conditionUncast.(map[string]interface{})
			name, found, err := unstructured.NestedString(condition, "type")
			if !found || err != nil || !strings.EqualFold(name, conditionStr) {
				continue
			}
			status, found, err := unstructured.NestedString(condition, "status")
			if !found || err != nil {
				continue
			}
			return strings.EqualFold(status, "True"), nil
		}

		return false, nil
	}
	return cmgr.WatchEventsForCondition(timeoutSecs, ns, name, gvr, condFunc)

}

// WaitForResourceReplicas waits until a resource has the given number of replicas in "ready" state
func (cmgr clientMgr) WaitForResourceReplicas(timeoutSecs int, ns, name, replicas string, gvr schema.GroupVersionResource) (bool, error) {
	var condFunc ConditionFunction = func(event watch.Event, obj *unstructured.Unstructured) (bool, error) {
		readyReplicas, found, err := unstructured.NestedString(obj.Object, "status", "readyReplicas")
		if err != nil {
			return false, err
		}
		if !found {
			return false, nil
		}
		return strings.EqualFold(readyReplicas, replicas), nil
	}
	return cmgr.WatchEventsForCondition(timeoutSecs, ns, name, gvr, condFunc)
}
