package controller_test

import (
	"context"
	"fmt"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/rotisserie/eris"
	. "github.com/solo-io/go-utils/testutils"
	core_types "github.com/solo-io/mesh-projects/pkg/api/core.zephyr.solo.io/v1alpha1/types"
	discoveryv1alpha1 "github.com/solo-io/mesh-projects/pkg/api/discovery.zephyr.solo.io/v1alpha1"
	discovery_types "github.com/solo-io/mesh-projects/pkg/api/discovery.zephyr.solo.io/v1alpha1/types"
	networkingv1alpha1 "github.com/solo-io/mesh-projects/pkg/api/networking.zephyr.solo.io/v1alpha1"
	v1alpha1_types "github.com/solo-io/mesh-projects/pkg/api/networking.zephyr.solo.io/v1alpha1/types"
	mock_core "github.com/solo-io/mesh-projects/pkg/clients/zephyr/discovery/mocks"
	"github.com/solo-io/mesh-projects/services/mesh-group/pkg/controller"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("validator", func() {

	var (
		ctrl       *gomock.Controller
		validator  controller.MeshGroupValidator
		meshClient *mock_core.MockMeshClient
		ctx        context.Context

		testErr = eris.New("hello")
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		meshClient = mock_core.NewMockMeshClient(ctrl)
		validator = controller.MeshGroupValidatorProvider(meshClient)
		ctx = context.TODO()
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	It("will return an error if mesh client list fails", func() {
		meshClient.EXPECT().List(ctx).Return(nil, testErr)
		status, err := validator.Validate(ctx, nil)
		Expect(err).To(HaveOccurred())
		Expect(err).To(Equal(testErr))
		Expect(status).To(Equal(v1alpha1_types.MeshGroupStatus{
			Config: v1alpha1_types.MeshGroupStatus_PROCESSING_ERROR,
		}))
	})

	It("will return an error if a mesh ref doesn't exist", func() {
		ref := &core_types.ResourceRef{
			Name:      "incorrect",
			Namespace: "ref",
		}
		meshClient.EXPECT().List(ctx).Return(&discoveryv1alpha1.MeshList{}, nil)
		status, err := validator.Validate(ctx, &networkingv1alpha1.MeshGroup{
			Spec: v1alpha1_types.MeshGroupSpec{
				Meshes: []*core_types.ResourceRef{ref},
			},
		})
		Expect(err).To(HaveOccurred())
		Expect(err).To(HaveInErrorChain(controller.InvalidMeshRefsError([]string{
			fmt.Sprintf("%s.%s", ref.GetName(), ref.GetNamespace()),
		})))
		Expect(status).To(Equal(v1alpha1_types.MeshGroupStatus{
			Config: v1alpha1_types.MeshGroupStatus_INVALID,
		}))
	})

	It("will return an error if a non-istio mesh is referenced", func() {
		ref := &core_types.ResourceRef{
			Name:      "valid",
			Namespace: "ref",
		}
		mesh := discoveryv1alpha1.Mesh{
			ObjectMeta: v1.ObjectMeta{
				Name:      ref.GetName(),
				Namespace: ref.GetNamespace(),
			},
			Spec: discovery_types.MeshSpec{
				MeshType: &discovery_types.MeshSpec_ConsulConnect{},
			},
		}
		meshClient.EXPECT().List(ctx).Return(&discoveryv1alpha1.MeshList{
			Items: []discoveryv1alpha1.Mesh{mesh},
		}, nil)
		status, err := validator.Validate(ctx, &networkingv1alpha1.MeshGroup{
			Spec: v1alpha1_types.MeshGroupSpec{
				Meshes: []*core_types.ResourceRef{ref},
			},
		})
		Expect(err).To(HaveOccurred())
		Expect(err).To(HaveInErrorChain(controller.OnlyIstioSupportedError(mesh.Name)))
		Expect(status).To(Equal(v1alpha1_types.MeshGroupStatus{
			Config: v1alpha1_types.MeshGroupStatus_INVALID,
		}))
	})

	It("will return valid and no error if all went fine", func() {
		ref := &core_types.ResourceRef{
			Name:      "valid",
			Namespace: "ref",
		}
		mesh := discoveryv1alpha1.Mesh{
			ObjectMeta: v1.ObjectMeta{
				Name:      ref.GetName(),
				Namespace: ref.GetNamespace(),
			},
			Spec: discovery_types.MeshSpec{
				MeshType: &discovery_types.MeshSpec_Istio{
					Istio: &discovery_types.IstioMesh{},
				},
			},
		}
		meshClient.EXPECT().List(ctx).Return(&discoveryv1alpha1.MeshList{
			Items: []discoveryv1alpha1.Mesh{mesh},
		}, nil)
		status, err := validator.Validate(ctx, &networkingv1alpha1.MeshGroup{
			Spec: v1alpha1_types.MeshGroupSpec{
				Meshes: []*core_types.ResourceRef{ref},
			},
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(status).To(Equal(v1alpha1_types.MeshGroupStatus{
			Config: v1alpha1_types.MeshGroupStatus_VALID,
		}))
	})
})
