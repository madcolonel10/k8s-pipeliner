package builder_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/namely/k8s-pipeliner/pipeline/builder"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestContainersFromManifests(t *testing.T) {
	wd, _ := os.Getwd()

	t.Run("Deployment manifests are returned correctly", func(t *testing.T) {
		file := filepath.Join(wd, "testdata", "deployment.full.yml")
		group, err := builder.ContainersFromManifest(file)

		require.NoError(t, err, "error on retrieving the deployment manifests")

		assert.Len(t, group.Containers, 1)
		assert.Len(t, group.Annotations, 2)
		assert.Equal(t, "fake-namespace", group.Namespace)

		t.Run("Container VolumeMounts are copied in", func(t *testing.T) {
			c := group.Containers[0]

			require.Len(t, c.VolumeMounts, 1)
			assert.Equal(t, "configmap-volume", c.VolumeMounts[0].Name)
			assert.Equal(t, "/thisisthemount", c.VolumeMounts[0].MountPath)
			assert.Equal(t, true, c.VolumeMounts[0].ReadOnly)
		})
	})

	t.Run("Deployments schemes are converted to latest", func(t *testing.T) {
		file := filepath.Join(wd, "testdata", "deployment.v1beta1.yml")
		group, err := builder.ContainersFromManifest(file)

		require.NoError(t, err, "error on retrieving the deployment manifests")

		assert.Len(t, group.Containers, 1)
		assert.Len(t, group.Annotations, 2)
		assert.Equal(t, "fake-namespace", group.Namespace)
	})

	t.Run("Volume sources are copied", func(t *testing.T) {
		file := filepath.Join(wd, "testdata", "deployment.full.yml")
		group, err := builder.ContainersFromManifest(file)
		require.NoError(t, err)
		require.Len(t, group.VolumeSources, 3)

		t.Run("ConfigMaps are copied", func(t *testing.T) {
			cms := group.VolumeSources[0]
			require.NotNil(t, cms.ConfigMap)
			assert.Equal(t, cms.Type, "CONFIGMAP")
		})

		t.Run("Secrets are copied", func(t *testing.T) {
			sec := group.VolumeSources[1]
			require.NotNil(t, sec.Secret)
			assert.Equal(t, sec.Type, "SECRET")
		})

		t.Run("EmptyDirs are copied", func(t *testing.T) {
			ed := group.VolumeSources[2]
			require.NotNil(t, ed.EmptyDir)
			assert.Equal(t, ed.Type, "EMPTYDIR")
		})
	})
}