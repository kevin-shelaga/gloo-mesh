// Code generated by solo-kit. DO NOT EDIT.

package v1

import (
	"sync"
	"time"

	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/solo-kit/pkg/utils/errutils"
)

var (
	mInstallSnapshotIn  = stats.Int64("install.supergloo.solo.io/snap_emitter/snap_in", "The number of snapshots in", "1")
	mInstallSnapshotOut = stats.Int64("install.supergloo.solo.io/snap_emitter/snap_out", "The number of snapshots out", "1")

	installsnapshotInView = &view.View{
		Name:        "install.supergloo.solo.io_snap_emitter/snap_in",
		Measure:     mInstallSnapshotIn,
		Description: "The number of snapshots updates coming in",
		Aggregation: view.Count(),
		TagKeys:     []tag.Key{},
	}
	installsnapshotOutView = &view.View{
		Name:        "install.supergloo.solo.io/snap_emitter/snap_out",
		Measure:     mInstallSnapshotOut,
		Description: "The number of snapshots updates going out",
		Aggregation: view.Count(),
		TagKeys:     []tag.Key{},
	}
)

func init() {
	view.Register(installsnapshotInView, installsnapshotOutView)
}

type InstallEmitter interface {
	Register() error
	Install() InstallClient
	Snapshots(watchNamespaces []string, opts clients.WatchOpts) (<-chan *InstallSnapshot, <-chan error, error)
}

func NewInstallEmitter(installClient InstallClient) InstallEmitter {
	return NewInstallEmitterWithEmit(installClient, make(chan struct{}))
}

func NewInstallEmitterWithEmit(installClient InstallClient, emit <-chan struct{}) InstallEmitter {
	return &installEmitter{
		install:   installClient,
		forceEmit: emit,
	}
}

type installEmitter struct {
	forceEmit <-chan struct{}
	install   InstallClient
}

func (c *installEmitter) Register() error {
	if err := c.install.Register(); err != nil {
		return err
	}
	return nil
}

func (c *installEmitter) Install() InstallClient {
	return c.install
}

func (c *installEmitter) Snapshots(watchNamespaces []string, opts clients.WatchOpts) (<-chan *InstallSnapshot, <-chan error, error) {

	if len(watchNamespaces) == 0 {
		watchNamespaces = []string{""}
	}

	for _, ns := range watchNamespaces {
		if ns == "" && len(watchNamespaces) > 1 {
			return nil, nil, errors.Errorf("the \"\" namespace is used to watch all namespaces. Snapshots can either be tracked for " +
				"specific namespaces or \"\" AllNamespaces, but not both.")
		}
	}

	errs := make(chan error)
	var done sync.WaitGroup
	ctx := opts.Ctx
	/* Create channel for Install */
	type installListWithNamespace struct {
		list      InstallList
		namespace string
	}
	installChan := make(chan installListWithNamespace)

	for _, namespace := range watchNamespaces {
		/* Setup namespaced watch for Install */
		installNamespacesChan, installErrs, err := c.install.Watch(namespace, opts)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "starting Install watch")
		}

		done.Add(1)
		go func(namespace string) {
			defer done.Done()
			errutils.AggregateErrs(ctx, errs, installErrs, namespace+"-installs")
		}(namespace)

		/* Watch for changes and update snapshot */
		go func(namespace string) {
			for {
				select {
				case <-ctx.Done():
					return
				case installList := <-installNamespacesChan:
					select {
					case <-ctx.Done():
						return
					case installChan <- installListWithNamespace{list: installList, namespace: namespace}:
					}
				}
			}
		}(namespace)
	}

	snapshots := make(chan *InstallSnapshot)
	go func() {
		originalSnapshot := InstallSnapshot{}
		currentSnapshot := originalSnapshot.Clone()
		timer := time.NewTicker(time.Second * 1)
		sync := func() {
			if originalSnapshot.Hash() == currentSnapshot.Hash() {
				return
			}

			stats.Record(ctx, mInstallSnapshotOut.M(1))
			originalSnapshot = currentSnapshot.Clone()
			sentSnapshot := currentSnapshot.Clone()
			snapshots <- &sentSnapshot
		}

		for {
			record := func() { stats.Record(ctx, mInstallSnapshotIn.M(1)) }

			select {
			case <-timer.C:
				sync()
			case <-ctx.Done():
				close(snapshots)
				done.Wait()
				close(errs)
				return
			case <-c.forceEmit:
				sentSnapshot := currentSnapshot.Clone()
				snapshots <- &sentSnapshot
			case installNamespacedList := <-installChan:
				record()

				namespace := installNamespacedList.namespace
				installList := installNamespacedList.list

				currentSnapshot.Installs[namespace] = installList
			}
		}
	}()
	return snapshots, errs, nil
}