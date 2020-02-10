package consul

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/google/wire"
	"github.com/rotisserie/eris"
	mp_v1alpha1 "github.com/solo-io/mesh-projects/pkg/api/core.zephyr.solo.io/v1alpha1"
	v1alpha1_types "github.com/solo-io/mesh-projects/pkg/api/core.zephyr.solo.io/v1alpha1/types"
	"github.com/solo-io/mesh-projects/pkg/common/docker"
	"github.com/solo-io/mesh-projects/pkg/env"
	"github.com/solo-io/mesh-projects/services/mesh-discovery/pkg/discovery"
	k8s_apps_v1 "k8s.io/api/apps/v1"
	k8s_core_v1 "k8s.io/api/core/v1"
	k8s_meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	consulServerArg           = "-server"
	normalizedConsulImagePath = "library/consul"
	meshNamePrefix            = "consul"
)

var (
	WireProviderSet = wire.NewSet(
		NewConsulConnectInstallationFinder,
		NewConsulMeshFinder,
	)
	DiscoveryLabels = map[string]string{
		"discovered_by": "consul-mesh-discovery",
	}
	ErrorDetectingDeployment = func(err error) error {
		return eris.Wrap(err, "Error while detecting consul deployment")
	}

	datacenterRegex = regexp.MustCompile("-datacenter=([a-zA-Z0-9]*)")
)

// disambiguates this MeshFinder from the other MeshFinder implementations so that wire stays happy
type ConsulConnectMeshFinder discovery.MeshFinder

func NewConsulMeshFinder(imageNameParser docker.ImageNameParser, consulConnectInstallationFinder ConsulConnectInstallationFinder) ConsulConnectMeshFinder {
	return &consulMeshFinder{
		consulConnectInstallationFinder,
		imageNameParser,
	}
}

type consulMeshFinder struct {
	consulConnectInstallationFinder ConsulConnectInstallationFinder
	imageNameParser                 docker.ImageNameParser
}

func (c *consulMeshFinder) ScanDeployment(_ context.Context, deployment *k8s_apps_v1.Deployment) (*mp_v1alpha1.Mesh, error) {
	for _, container := range deployment.Spec.Template.Spec.Containers {
		isConsulInstallation, err := c.consulConnectInstallationFinder.IsConsulConnect(container)
		if err != nil {
			return nil, ErrorDetectingDeployment(err)
		}

		if !isConsulInstallation {
			continue
		}

		return c.buildConsulMeshObject(deployment, container, env.DefaultWriteNamespace)
	}

	return nil, nil
}

// returns an error if the image name is un-parsable
func (c *consulMeshFinder) buildConsulMeshObject(deployment *k8s_apps_v1.Deployment, container k8s_core_v1.Container, writeNamespace string) (*mp_v1alpha1.Mesh, error) {
	parsedImage, err := c.imageNameParser.Parse(container.Image)
	if err != nil {
		return nil, err
	}

	imageVersion := parsedImage.Tag
	if parsedImage.Digest != "" {
		imageVersion = parsedImage.Digest
	}

	return &mp_v1alpha1.Mesh{
		ObjectMeta: k8s_meta_v1.ObjectMeta{
			Name:      buildMeshName(deployment, container),
			Namespace: env.DefaultWriteNamespace,
			Labels:    DiscoveryLabels,
		},
		Spec: v1alpha1_types.MeshSpec{
			MeshType: &v1alpha1_types.MeshSpec_ConsulConnect{
				ConsulConnect: &v1alpha1_types.ConsulConnectMesh{
					Installation: &v1alpha1_types.MeshInstallation{
						InstallationNamespace: deployment.GetNamespace(),
						Version:               imageVersion,
					},
				},
			},
			Cluster: &v1alpha1_types.ResourceRef{
				Name:      deployment.GetClusterName(),
				Namespace: env.DefaultWriteNamespace,
			},
		},
	}, nil
}

// returns "consul(-$datacenterName)-$installNamespace(-$clusterName)"
func buildMeshName(deployment *k8s_apps_v1.Deployment, container k8s_core_v1.Container) string {
	meshName := meshNamePrefix

	cmd := strings.Join(container.Command, " ")
	datacenterNameMatch := datacenterRegex.FindStringSubmatch(cmd)

	if len(datacenterNameMatch) == 2 {
		meshName += fmt.Sprintf("-%s", datacenterNameMatch[1])
	}

	meshName += fmt.Sprintf("-%s", deployment.Namespace)

	if deployment.ClusterName != "" {
		meshName += fmt.Sprintf("-%s", deployment.ClusterName)
	}

	return meshName
}