package adoption

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/controller-runtime/pkg/client"

	coordinationv1alpha1 "github.com/thetechnick/package-operator/apis/coordination/v1alpha1"
	"github.com/thetechnick/package-operator/internal/testutil"
)

func newRRA[T operandPtr[O], O operand](
	c client.Client, o O) *RoundRobinAdoptionReconciler[T, O] {
	return &RoundRobinAdoptionReconciler[T, O]{
		client: c,
	}
}

// Noop case
func TestRoundRobinAdoptionReconciler(t *testing.T) {
	c := testutil.NewClient()
	ctx := context.Background()

	var (
		operandListOpts client.ListOptions
	)

	c.
		On(
			"List",
			mock.Anything, mock.IsType(&unstructured.UnstructuredList{}),
			mock.Anything,
		).
		Run(func(args mock.Arguments) {
			list := args.Get(1).(*unstructured.UnstructuredList)
			list.Items = make([]unstructured.Unstructured, 4)

			optFns := args.Get(2).([]client.ListOption)
			for _, fn := range optFns {
				fn.ApplyToList(&operandListOpts)
			}
		}).
		Once().
		Return(nil)

	var updatedOperands []unstructured.Unstructured
	c.
		On(
			"Update",
			mock.Anything, mock.IsType(&unstructured.Unstructured{}),
			mock.Anything,
		).
		Run(func(args mock.Arguments) {
			operand := args.Get(1).(*unstructured.Unstructured)
			updatedOperands = append(updatedOperands, *operand)
		}).
		Return(nil)

	adoption := &coordinationv1alpha1.Adoption{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "test",
		},
		Spec: coordinationv1alpha1.AdoptionSpec{
			Strategy: coordinationv1alpha1.AdoptionStrategy{
				Type: coordinationv1alpha1.AdoptionStrategyRoundRobin,
				RoundRobin: &coordinationv1alpha1.AdoptionStrategyRoundRobinSpec{
					Always: map[string]string{"always": "always"},
					Options: []map[string]string{
						{"option1": "option1"},
						{"option2": "option2"},
						{"option3": "option3"},
					},
				},
			},
			TargetAPI: coordinationv1alpha1.TargetAPI{
				Group:   "test.thetechnick.ninja",
				Version: "v1alpha1",
				Kind:    "Test",
			},
		},
	}
	rra := newRRA(c, coordinationv1alpha1.Adoption{})
	res, err := rra.Reconcile(ctx, adoption)

	require.NoError(t, err)
	assert.True(t, res.IsZero())

	if !c.AssertExpectations(t) {
		return
	}

	if assert.Len(t, updatedOperands, 4) {
		assert.Equal(t, updatedOperands[0].GetLabels(), map[string]string{
			"always":  "always",
			"option1": "option1",
		})
		assert.Equal(t, updatedOperands[1].GetLabels(), map[string]string{
			"always":  "always",
			"option2": "option2",
		})
		assert.Equal(t, updatedOperands[2].GetLabels(), map[string]string{
			"always":  "always",
			"option3": "option3",
		})
		assert.Equal(t, updatedOperands[3].GetLabels(), map[string]string{
			"always":  "always",
			"option1": "option1",
		})
	}

	if assert.NotNil(t, adoption.Status.RoundRobin) {
		assert.Equal(t, 0, adoption.Status.RoundRobin.LastIndex)
	}

	assert.Equal(t, adoption.GetNamespace(), operandListOpts.Namespace)
	assert.Equal(t,
		"!option1,!option2,!option3",
		operandListOpts.LabelSelector.String())
}

// Noop case
func TestRoundRobinAdoptionReconciler_WrongStrategy(t *testing.T) {
	c := testutil.NewClient()
	ctx := context.Background()

	rra := newRRA(c, coordinationv1alpha1.Adoption{})
	res, err := rra.Reconcile(ctx, &coordinationv1alpha1.Adoption{})

	require.NoError(t, err)
	assert.True(t, res.IsZero())
}

func Test_roundRobinIndex(t *testing.T) {
	var (
		index   = -1 // no start index
		results []int
	)
	for i := 0; i < 4; i++ {
		index = roundRobinIndex(index, 2)
		results = append(results, index)
	}

	assert.Equal(t, []int{0, 1, 2, 0}, results)
}
