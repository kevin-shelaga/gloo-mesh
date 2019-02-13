// Code generated by solo-kit. DO NOT EDIT.

// +build solokit

package v1

import (
	"context"
	"os"
	"time"

	gloo_solo_io "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	encryption_istio_io "github.com/solo-io/supergloo/pkg/api/external/istio/encryption/v1"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/go-utils/kubeutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	kuberc "github.com/solo-io/solo-kit/pkg/api/v1/clients/kube"
	"github.com/solo-io/solo-kit/pkg/utils/log"
	"github.com/solo-io/solo-kit/test/helpers"
	"github.com/solo-io/solo-kit/test/setup"
	"k8s.io/client-go/rest"

	// Needed to run tests in GKE
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"

	// From https://github.com/kubernetes/client-go/blob/53c7adfd0294caa142d961e1f780f74081d5b15f/examples/out-of-cluster-client-configuration/main.go#L31
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

var _ = Describe("V1Emitter", func() {
	if os.Getenv("RUN_KUBE_TESTS") != "1" {
		log.Printf("This test creates kubernetes resources and is disabled by default. To enable, set RUN_KUBE_TESTS=1 in your env.")
		return
	}
	var (
		namespace1               string
		namespace2               string
		name1, name2             = "angela" + helpers.RandString(3), "bob" + helpers.RandString(3)
		cfg                      *rest.Config
		emitter                  TranslatorEmitter
		secretClient             gloo_solo_io.SecretClient
		upstreamClient           gloo_solo_io.UpstreamClient
		meshClient               MeshClient
		routingRuleClient        RoutingRuleClient
		istioCacertsSecretClient encryption_istio_io.IstioCacertsSecretClient
	)

	BeforeEach(func() {
		namespace1 = helpers.RandString(8)
		namespace2 = helpers.RandString(8)
		var err error
		cfg, err = kubeutils.GetConfig("", "")
		Expect(err).NotTo(HaveOccurred())
		err = setup.SetupKubeForTest(namespace1)
		Expect(err).NotTo(HaveOccurred())
		err = setup.SetupKubeForTest(namespace2)
		Expect(err).NotTo(HaveOccurred())
		var kube kubernetes.Interface
		// Secret Constructor
		kube, err = kubernetes.NewForConfig(cfg)
		Expect(err).NotTo(HaveOccurred())

		secretClientFactory := &factory.KubeSecretClientFactory{
			Clientset: kube,
		}
		secretClient, err = gloo_solo_io.NewSecretClient(secretClientFactory)
		Expect(err).NotTo(HaveOccurred())
		// Upstream Constructor
		upstreamClientFactory := &factory.KubeResourceClientFactory{
			Crd:         gloo_solo_io.UpstreamCrd,
			Cfg:         cfg,
			SharedCache: kuberc.NewKubeCache(),
		}
		upstreamClient, err = gloo_solo_io.NewUpstreamClient(upstreamClientFactory)
		Expect(err).NotTo(HaveOccurred())
		// Mesh Constructor
		meshClientFactory := &factory.KubeResourceClientFactory{
			Crd:         MeshCrd,
			Cfg:         cfg,
			SharedCache: kuberc.NewKubeCache(),
		}
		meshClient, err = NewMeshClient(meshClientFactory)
		Expect(err).NotTo(HaveOccurred())
		// RoutingRule Constructor
		routingRuleClientFactory := &factory.KubeResourceClientFactory{
			Crd:         RoutingRuleCrd,
			Cfg:         cfg,
			SharedCache: kuberc.NewKubeCache(),
		}
		routingRuleClient, err = NewRoutingRuleClient(routingRuleClientFactory)
		Expect(err).NotTo(HaveOccurred())
		// IstioCacertsSecret Constructor
		kube, err = kubernetes.NewForConfig(cfg)
		Expect(err).NotTo(HaveOccurred())

		istioCacertsSecretClientFactory := &factory.KubeConfigMapClientFactory{
			Clientset: kube,
		}
		istioCacertsSecretClient, err = encryption_istio_io.NewIstioCacertsSecretClient(istioCacertsSecretClientFactory)
		Expect(err).NotTo(HaveOccurred())
		emitter = NewTranslatorEmitter(secretClient, upstreamClient, meshClient, routingRuleClient, istioCacertsSecretClient)
	})
	AfterEach(func() {
		setup.TeardownKube(namespace1)
		setup.TeardownKube(namespace2)
	})
	It("tracks snapshots on changes to any resource", func() {
		ctx := context.Background()
		err := emitter.Register()
		Expect(err).NotTo(HaveOccurred())

		snapshots, errs, err := emitter.Snapshots([]string{namespace1, namespace2}, clients.WatchOpts{
			Ctx:         ctx,
			RefreshRate: time.Second,
		})
		Expect(err).NotTo(HaveOccurred())

		var snap *TranslatorSnapshot

		/*
			Secret
		*/

		assertSnapshotSecrets := func(expectSecrets gloo_solo_io.SecretList, unexpectSecrets gloo_solo_io.SecretList) {
		drain:
			for {
				select {
				case snap = <-snapshots:
					for _, expected := range expectSecrets {
						if _, err := snap.Secrets.List().Find(expected.Metadata.Ref().Strings()); err != nil {
							continue drain
						}
					}
					for _, unexpected := range unexpectSecrets {
						if _, err := snap.Secrets.List().Find(unexpected.Metadata.Ref().Strings()); err == nil {
							continue drain
						}
					}
					break drain
				case err := <-errs:
					Expect(err).NotTo(HaveOccurred())
				case <-time.After(time.Second * 10):
					nsList1, _ := secretClient.List(namespace1, clients.ListOpts{})
					nsList2, _ := secretClient.List(namespace2, clients.ListOpts{})
					combined := nsList1.ByNamespace()
					combined.Add(nsList2...)
					Fail("expected final snapshot before 10 seconds. expected " + log.Sprintf("%v", combined))
				}
			}
		}
		secret1a, err := secretClient.Write(gloo_solo_io.NewSecret(namespace1, name1), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		secret1b, err := secretClient.Write(gloo_solo_io.NewSecret(namespace2, name1), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		assertSnapshotSecrets(gloo_solo_io.SecretList{secret1a, secret1b}, nil)
		secret2a, err := secretClient.Write(gloo_solo_io.NewSecret(namespace1, name2), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		secret2b, err := secretClient.Write(gloo_solo_io.NewSecret(namespace2, name2), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		assertSnapshotSecrets(gloo_solo_io.SecretList{secret1a, secret1b, secret2a, secret2b}, nil)

		err = secretClient.Delete(secret2a.Metadata.Namespace, secret2a.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		err = secretClient.Delete(secret2b.Metadata.Namespace, secret2b.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		assertSnapshotSecrets(gloo_solo_io.SecretList{secret1a, secret1b}, gloo_solo_io.SecretList{secret2a, secret2b})

		err = secretClient.Delete(secret1a.Metadata.Namespace, secret1a.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		err = secretClient.Delete(secret1b.Metadata.Namespace, secret1b.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		assertSnapshotSecrets(nil, gloo_solo_io.SecretList{secret1a, secret1b, secret2a, secret2b})

		/*
			Upstream
		*/

		assertSnapshotUpstreams := func(expectUpstreams gloo_solo_io.UpstreamList, unexpectUpstreams gloo_solo_io.UpstreamList) {
		drain:
			for {
				select {
				case snap = <-snapshots:
					for _, expected := range expectUpstreams {
						if _, err := snap.Upstreams.List().Find(expected.Metadata.Ref().Strings()); err != nil {
							continue drain
						}
					}
					for _, unexpected := range unexpectUpstreams {
						if _, err := snap.Upstreams.List().Find(unexpected.Metadata.Ref().Strings()); err == nil {
							continue drain
						}
					}
					break drain
				case err := <-errs:
					Expect(err).NotTo(HaveOccurred())
				case <-time.After(time.Second * 10):
					nsList1, _ := upstreamClient.List(namespace1, clients.ListOpts{})
					nsList2, _ := upstreamClient.List(namespace2, clients.ListOpts{})
					combined := nsList1.ByNamespace()
					combined.Add(nsList2...)
					Fail("expected final snapshot before 10 seconds. expected " + log.Sprintf("%v", combined))
				}
			}
		}
		upstream1a, err := upstreamClient.Write(gloo_solo_io.NewUpstream(namespace1, name1), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		upstream1b, err := upstreamClient.Write(gloo_solo_io.NewUpstream(namespace2, name1), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		assertSnapshotUpstreams(gloo_solo_io.UpstreamList{upstream1a, upstream1b}, nil)
		upstream2a, err := upstreamClient.Write(gloo_solo_io.NewUpstream(namespace1, name2), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		upstream2b, err := upstreamClient.Write(gloo_solo_io.NewUpstream(namespace2, name2), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		assertSnapshotUpstreams(gloo_solo_io.UpstreamList{upstream1a, upstream1b, upstream2a, upstream2b}, nil)

		err = upstreamClient.Delete(upstream2a.Metadata.Namespace, upstream2a.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		err = upstreamClient.Delete(upstream2b.Metadata.Namespace, upstream2b.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		assertSnapshotUpstreams(gloo_solo_io.UpstreamList{upstream1a, upstream1b}, gloo_solo_io.UpstreamList{upstream2a, upstream2b})

		err = upstreamClient.Delete(upstream1a.Metadata.Namespace, upstream1a.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		err = upstreamClient.Delete(upstream1b.Metadata.Namespace, upstream1b.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		assertSnapshotUpstreams(nil, gloo_solo_io.UpstreamList{upstream1a, upstream1b, upstream2a, upstream2b})

		/*
			Mesh
		*/

		assertSnapshotMeshes := func(expectMeshes MeshList, unexpectMeshes MeshList) {
		drain:
			for {
				select {
				case snap = <-snapshots:
					for _, expected := range expectMeshes {
						if _, err := snap.Meshes.List().Find(expected.Metadata.Ref().Strings()); err != nil {
							continue drain
						}
					}
					for _, unexpected := range unexpectMeshes {
						if _, err := snap.Meshes.List().Find(unexpected.Metadata.Ref().Strings()); err == nil {
							continue drain
						}
					}
					break drain
				case err := <-errs:
					Expect(err).NotTo(HaveOccurred())
				case <-time.After(time.Second * 10):
					nsList1, _ := meshClient.List(namespace1, clients.ListOpts{})
					nsList2, _ := meshClient.List(namespace2, clients.ListOpts{})
					combined := nsList1.ByNamespace()
					combined.Add(nsList2...)
					Fail("expected final snapshot before 10 seconds. expected " + log.Sprintf("%v", combined))
				}
			}
		}
		mesh1a, err := meshClient.Write(NewMesh(namespace1, name1), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		mesh1b, err := meshClient.Write(NewMesh(namespace2, name1), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		assertSnapshotMeshes(MeshList{mesh1a, mesh1b}, nil)
		mesh2a, err := meshClient.Write(NewMesh(namespace1, name2), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		mesh2b, err := meshClient.Write(NewMesh(namespace2, name2), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		assertSnapshotMeshes(MeshList{mesh1a, mesh1b, mesh2a, mesh2b}, nil)

		err = meshClient.Delete(mesh2a.Metadata.Namespace, mesh2a.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		err = meshClient.Delete(mesh2b.Metadata.Namespace, mesh2b.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		assertSnapshotMeshes(MeshList{mesh1a, mesh1b}, MeshList{mesh2a, mesh2b})

		err = meshClient.Delete(mesh1a.Metadata.Namespace, mesh1a.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		err = meshClient.Delete(mesh1b.Metadata.Namespace, mesh1b.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		assertSnapshotMeshes(nil, MeshList{mesh1a, mesh1b, mesh2a, mesh2b})

		/*
			RoutingRule
		*/

		assertSnapshotRoutingrules := func(expectRoutingrules RoutingRuleList, unexpectRoutingrules RoutingRuleList) {
		drain:
			for {
				select {
				case snap = <-snapshots:
					for _, expected := range expectRoutingrules {
						if _, err := snap.Routingrules.List().Find(expected.Metadata.Ref().Strings()); err != nil {
							continue drain
						}
					}
					for _, unexpected := range unexpectRoutingrules {
						if _, err := snap.Routingrules.List().Find(unexpected.Metadata.Ref().Strings()); err == nil {
							continue drain
						}
					}
					break drain
				case err := <-errs:
					Expect(err).NotTo(HaveOccurred())
				case <-time.After(time.Second * 10):
					nsList1, _ := routingRuleClient.List(namespace1, clients.ListOpts{})
					nsList2, _ := routingRuleClient.List(namespace2, clients.ListOpts{})
					combined := nsList1.ByNamespace()
					combined.Add(nsList2...)
					Fail("expected final snapshot before 10 seconds. expected " + log.Sprintf("%v", combined))
				}
			}
		}
		routingRule1a, err := routingRuleClient.Write(NewRoutingRule(namespace1, name1), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		routingRule1b, err := routingRuleClient.Write(NewRoutingRule(namespace2, name1), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		assertSnapshotRoutingrules(RoutingRuleList{routingRule1a, routingRule1b}, nil)
		routingRule2a, err := routingRuleClient.Write(NewRoutingRule(namespace1, name2), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		routingRule2b, err := routingRuleClient.Write(NewRoutingRule(namespace2, name2), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		assertSnapshotRoutingrules(RoutingRuleList{routingRule1a, routingRule1b, routingRule2a, routingRule2b}, nil)

		err = routingRuleClient.Delete(routingRule2a.Metadata.Namespace, routingRule2a.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		err = routingRuleClient.Delete(routingRule2b.Metadata.Namespace, routingRule2b.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		assertSnapshotRoutingrules(RoutingRuleList{routingRule1a, routingRule1b}, RoutingRuleList{routingRule2a, routingRule2b})

		err = routingRuleClient.Delete(routingRule1a.Metadata.Namespace, routingRule1a.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		err = routingRuleClient.Delete(routingRule1b.Metadata.Namespace, routingRule1b.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		assertSnapshotRoutingrules(nil, RoutingRuleList{routingRule1a, routingRule1b, routingRule2a, routingRule2b})

		/*
			IstioCacertsSecret
		*/

		assertSnapshotIstiocerts := func(expectIstiocerts encryption_istio_io.IstioCacertsSecretList, unexpectIstiocerts encryption_istio_io.IstioCacertsSecretList) {
		drain:
			for {
				select {
				case snap = <-snapshots:
					for _, expected := range expectIstiocerts {
						if _, err := snap.Istiocerts.List().Find(expected.Metadata.Ref().Strings()); err != nil {
							continue drain
						}
					}
					for _, unexpected := range unexpectIstiocerts {
						if _, err := snap.Istiocerts.List().Find(unexpected.Metadata.Ref().Strings()); err == nil {
							continue drain
						}
					}
					break drain
				case err := <-errs:
					Expect(err).NotTo(HaveOccurred())
				case <-time.After(time.Second * 10):
					nsList1, _ := istioCacertsSecretClient.List(namespace1, clients.ListOpts{})
					nsList2, _ := istioCacertsSecretClient.List(namespace2, clients.ListOpts{})
					combined := nsList1.ByNamespace()
					combined.Add(nsList2...)
					Fail("expected final snapshot before 10 seconds. expected " + log.Sprintf("%v", combined))
				}
			}
		}
		istioCacertsSecret1a, err := istioCacertsSecretClient.Write(encryption_istio_io.NewIstioCacertsSecret(namespace1, name1), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		istioCacertsSecret1b, err := istioCacertsSecretClient.Write(encryption_istio_io.NewIstioCacertsSecret(namespace2, name1), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		assertSnapshotIstiocerts(encryption_istio_io.IstioCacertsSecretList{istioCacertsSecret1a, istioCacertsSecret1b}, nil)
		istioCacertsSecret2a, err := istioCacertsSecretClient.Write(encryption_istio_io.NewIstioCacertsSecret(namespace1, name2), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		istioCacertsSecret2b, err := istioCacertsSecretClient.Write(encryption_istio_io.NewIstioCacertsSecret(namespace2, name2), clients.WriteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		assertSnapshotIstiocerts(encryption_istio_io.IstioCacertsSecretList{istioCacertsSecret1a, istioCacertsSecret1b, istioCacertsSecret2a, istioCacertsSecret2b}, nil)

		err = istioCacertsSecretClient.Delete(istioCacertsSecret2a.Metadata.Namespace, istioCacertsSecret2a.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		err = istioCacertsSecretClient.Delete(istioCacertsSecret2b.Metadata.Namespace, istioCacertsSecret2b.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		assertSnapshotIstiocerts(encryption_istio_io.IstioCacertsSecretList{istioCacertsSecret1a, istioCacertsSecret1b}, encryption_istio_io.IstioCacertsSecretList{istioCacertsSecret2a, istioCacertsSecret2b})

		err = istioCacertsSecretClient.Delete(istioCacertsSecret1a.Metadata.Namespace, istioCacertsSecret1a.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())
		err = istioCacertsSecretClient.Delete(istioCacertsSecret1b.Metadata.Namespace, istioCacertsSecret1b.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		assertSnapshotIstiocerts(nil, encryption_istio_io.IstioCacertsSecretList{istioCacertsSecret1a, istioCacertsSecret1b, istioCacertsSecret2a, istioCacertsSecret2b})
	})
})