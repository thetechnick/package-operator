package presets

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	manifestsv1alpha1 "package-operator.run/apis/manifests/v1alpha1"
)

var statefulSet = unstructured.Unstructured{
	Object: map[string]interface{}{
		"apiVersion": "apps/v1",
		"kind":       "StatefulSet",
		"metadata": map[string]interface{}{
			"name":      "banana",
			"namespace": "fruits",
		},
		"spec": map[string]interface{}{
			"replicas": int64(1),
			"template": map[string]interface{}{
				"spec": map[string]interface{}{
					"affinity": nil,
					"containers": []interface{}{
						map[string]interface{}{
							"name":  "banana",
							"image": "quay.io/package-operator/banana:latest",
							"env": []interface{}{
								map[string]interface{}{
									"name":  "HTTP_PROXY",
									"value": "xxx",
								},
							},
						},
					},
				},
			},
		},
	},
}

func TestStatefulSet(t *testing.T) {
	t.Parallel()
	scheme := &apiextensionsv1.JSONSchemaProps{}
	image := &ImageContainer{}
	out, ok, err := Parametrize(*statefulSet.DeepCopy(), scheme, image, ParametrizeOptions{
		Namespaces:    true,
		Replicas:      true,
		Images:        true,
		Resources:     true,
		NodeSelectors: true,

		// Need to figure out how to test with uuids present.
		// Tolerations: true,
		// Env:           true
	})
	require.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(t, `apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: banana
  namespace: {{ default (index .config.namespaces "fruits") .config.namespace }}
spec:
  replicas: {{ index .config "statefulsets" "fruits" "banana" "replicas" }}
  template:
    spec:
      affinity: null
      containers:
      - env:
        - name: HTTP_PROXY
          value: xxx
        image: {{ index .images "banana" }}
        name: banana
        resources: {{ index .config "statefulsets" "fruits" "banana" "containers" "banana" "resources" | toJson }}
      nodeSelector: {{ index .config "statefulsets" "fruits" "banana" "nodeSelector" | toJson }}
`, string(out))
	assert.Equal(t, []manifestsv1alpha1.PackageManifestImage{
		{
			Name:  "banana",
			Image: "quay.io/package-operator/banana:latest",
		},
	}, image.List())
}
