package health

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/watch"

	"github.com/ForgeRock/forgeops-cli/internal/k8s"
	imock "github.com/ForgeRock/forgeops-cli/internal/mock"
)

// TESTING OF HEALTH/RESOURCE LOGIC

// attributes used when testing a resource
// used to mock responses from events condition
type tResource struct {
	// name of resource
	rname string
	// return value from WatchEventsForCondition
	conditionMet bool
	// return err value from WatchEventsForCondition
	err error
}

// newHealthFromResources builds health object for given resources
func newHealthFromResources(resources []tResource) *Health {
	builtResources := make([]*Resource, 0)
	for _, r := range resources {
		resource := &Resource{
			Group:      "mytestgroup",
			APIVersion: "v1alpha1",
			Resource:   "mytestresource",
			Name:       r.rname,
			Namespace:  "test_namespace",
			Checks: []*Check{
				{
					Expression: "status.state == \"Completed\"",
					Timeout:    metav1.Duration{Duration: 1 * time.Second},
				},
			},
		}
		builtResources = append(builtResources, resource)
	}
	return &Health{
		Spec: V1AlphaHealthSpec{
			Resources: builtResources,
		},
		Metadata: metav1.ObjectMeta{
			Name: "testhealth",
		},
	}
}

// TestHealthyLogic tests Heatlh/Resource logic
func TestHealthyLogic(t *testing.T) {

	// test data
	tdResources := []struct {
		resources   []tResource
		expect      bool
		expectedErr error
	}{
		// multiple passing
		{
			resources: []tResource{
				{"r1", true, nil},
				{"r2", true, nil},
			},
			expect:      true,
			expectedErr: nil,
		},
		// a resource times out, should only be unhealthy
		{
			resources: []tResource{
				{"r1", true, nil},
				{"r2", true, k8s.ErrWatchTimeout},
			},
			expect:      false,
			expectedErr: nil,
		},
		// errors return err
		{
			resources: []tResource{
				{"r2", true, errors.New("test error")},
				{"r1", false, nil},
			},
			expect:      false,
			expectedErr: nil,
		},
		// all unhealthy
		{
			resources: []tResource{
				{"r1", false, nil},
				{"r2", false, nil},
			},
			expect:      false,
			expectedErr: nil,
		},
		// some unhealthy
		{
			resources: []tResource{
				{"r1", false, nil},
				{"r2", true, nil},
			},
			expect:      false,
			expectedErr: nil,
		},
	}

	for _, tc := range tdResources {
		testHealth := newHealthFromResources(tc.resources)
		testClientMgr := &imock.ClientMgr{}
		testClientMgr.On("Namespace").Return("test_namespace", nil)
		//
		for _, resource := range tc.resources {
			testClientMgr.On("WatchEventsForCondition",
				1,
				"test_namespace",
				resource.rname,
				mock.AnythingOfType("schema.GroupVersionResource"),
				mock.AnythingOfType("k8s.ConditionFunction"),
			).Return(resource.conditionMet, resource.err)
		}

		res, resultErr := testHealth.CheckResources(testClientMgr, false)
		if resultErr != tc.expectedErr {
			t.Errorf("expected no error but found %s", resultErr.Error())
		}
		if res != tc.expect {
			t.Error("expected check to pass")
		}
	}
}

// END TESTING OF HEALTH/RESOURCE LOGIC

// TESTING OF CONDITON EXPRESSION

// TestConditionExpression tests expressions
func TestConditionExpression(t *testing.T) {
	testEvent := watch.Event{}

	// test table data
	td := []struct {
		// comment about test case
		testComment string
		// obj during event
		obj map[string]interface{}
		// expected return result
		expectedResult bool
		// expected error
		expectedError error
		// test expression
		testExpression string
	}{
		// test cases
		{
			testComment: "status object of Ready should return true",
			// spec of object in "conditionExpression" arg
			obj: map[string]interface{}{
				"status": map[string]interface{}{
					"state": "Ready",
				},
			},
			expectedError:  nil,
			expectedResult: true,
			testExpression: "status.state == \"Ready\"",
		},
		{
			testComment: "status of object NotReady should return false",
			obj: map[string]interface{}{
				"status": map[string]interface{}{
					"state": "NotReady",
				},
			},
			expectedError:  nil,
			expectedResult: false,
			testExpression: "status.state == \"Ready\"",
		},
		{
			testComment: "non boolean expressions are errors",
			obj: map[string]interface{}{
				"status": map[string]interface{}{
					"state": "NotReady",
				},
			},
			expectedError:  ErrExpressionResult,
			expectedResult: false,
			testExpression: "sprintf(state.status)",
		},
	}

	for _, tc := range td {
		fn := conditionExpression(tc.testExpression)
		testObject := &unstructured.Unstructured{
			Object: tc.obj,
		}
		testResult, testErr := fn(testEvent, testObject)

		// don't test nil errors
		if !errors.Is(testErr, tc.expectedError) {
			t.Errorf("%s expected error: %+v, found: %+v", tc.testComment, tc.expectedError, testErr)
		}
		if testResult != tc.expectedResult {
			t.Errorf("%s expected: %t, found: %t", tc.testComment, tc.expectedResult, testResult)
		}
	}
}

// END TESTING OF CONDITON EXPRESSION
