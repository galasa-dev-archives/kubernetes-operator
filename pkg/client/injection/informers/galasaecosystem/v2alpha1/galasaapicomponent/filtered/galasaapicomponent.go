/*
 * Copyright contributors to the Galasa Project
 */
// Code generated by injection-gen. DO NOT EDIT.

package filtered

import (
	context "context"

	apisgalasaecosystemv2alpha1 "github.com/galasa-dev/galasa-kubernetes-operator/pkg/apis/galasaecosystem/v2alpha1"
	versioned "github.com/galasa-dev/galasa-kubernetes-operator/pkg/client/clientset/versioned"
	v2alpha1 "github.com/galasa-dev/galasa-kubernetes-operator/pkg/client/informers/externalversions/galasaecosystem/v2alpha1"
	client "github.com/galasa-dev/galasa-kubernetes-operator/pkg/client/injection/client"
	filtered "github.com/galasa-dev/galasa-kubernetes-operator/pkg/client/injection/informers/factory/filtered"
	galasaecosystemv2alpha1 "github.com/galasa-dev/galasa-kubernetes-operator/pkg/client/listers/galasaecosystem/v2alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	cache "k8s.io/client-go/tools/cache"
	controller "knative.dev/pkg/controller"
	injection "knative.dev/pkg/injection"
	logging "knative.dev/pkg/logging"
)

func init() {
	injection.Default.RegisterFilteredInformers(withInformer)
	injection.Dynamic.RegisterDynamicInformer(withDynamicInformer)
}

// Key is used for associating the Informer inside the context.Context.
type Key struct {
	Selector string
}

func withInformer(ctx context.Context) (context.Context, []controller.Informer) {
	untyped := ctx.Value(filtered.LabelKey{})
	if untyped == nil {
		logging.FromContext(ctx).Panic(
			"Unable to fetch labelkey from context.")
	}
	labelSelectors := untyped.([]string)
	infs := []controller.Informer{}
	for _, selector := range labelSelectors {
		f := filtered.Get(ctx, selector)
		inf := f.Galasa().V2alpha1().GalasaApiComponents()
		ctx = context.WithValue(ctx, Key{Selector: selector}, inf)
		infs = append(infs, inf.Informer())
	}
	return ctx, infs
}

func withDynamicInformer(ctx context.Context) context.Context {
	untyped := ctx.Value(filtered.LabelKey{})
	if untyped == nil {
		logging.FromContext(ctx).Panic(
			"Unable to fetch labelkey from context.")
	}
	labelSelectors := untyped.([]string)
	for _, selector := range labelSelectors {
		inf := &wrapper{client: client.Get(ctx), selector: selector}
		ctx = context.WithValue(ctx, Key{Selector: selector}, inf)
	}
	return ctx
}

// Get extracts the typed informer from the context.
func Get(ctx context.Context, selector string) v2alpha1.GalasaApiComponentInformer {
	untyped := ctx.Value(Key{Selector: selector})
	if untyped == nil {
		logging.FromContext(ctx).Panicf(
			"Unable to fetch github.com/galasa-dev/galasa-kubernetes-operator/pkg/client/informers/externalversions/galasaecosystem/v2alpha1.GalasaApiComponentInformer with selector %s from context.", selector)
	}
	return untyped.(v2alpha1.GalasaApiComponentInformer)
}

type wrapper struct {
	client versioned.Interface

	namespace string

	selector string
}

var _ v2alpha1.GalasaApiComponentInformer = (*wrapper)(nil)
var _ galasaecosystemv2alpha1.GalasaApiComponentLister = (*wrapper)(nil)

func (w *wrapper) Informer() cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(nil, &apisgalasaecosystemv2alpha1.GalasaApiComponent{}, 0, nil)
}

func (w *wrapper) Lister() galasaecosystemv2alpha1.GalasaApiComponentLister {
	return w
}

func (w *wrapper) GalasaApiComponents(namespace string) galasaecosystemv2alpha1.GalasaApiComponentNamespaceLister {
	return &wrapper{client: w.client, namespace: namespace, selector: w.selector}
}

func (w *wrapper) List(selector labels.Selector) (ret []*apisgalasaecosystemv2alpha1.GalasaApiComponent, err error) {
	reqs, err := labels.ParseToRequirements(w.selector)
	if err != nil {
		return nil, err
	}
	selector = selector.Add(reqs...)
	lo, err := w.client.GalasaV2alpha1().GalasaApiComponents(w.namespace).List(context.TODO(), v1.ListOptions{
		LabelSelector: selector.String(),
		// TODO(mattmoor): Incorporate resourceVersion bounds based on staleness criteria.
	})
	if err != nil {
		return nil, err
	}
	for idx := range lo.Items {
		ret = append(ret, &lo.Items[idx])
	}
	return ret, nil
}

func (w *wrapper) Get(name string) (*apisgalasaecosystemv2alpha1.GalasaApiComponent, error) {
	// TODO(mattmoor): Check that the fetched object matches the selector.
	return w.client.GalasaV2alpha1().GalasaApiComponents(w.namespace).Get(context.TODO(), name, v1.GetOptions{
		// TODO(mattmoor): Incorporate resourceVersion bounds based on staleness criteria.
	})
}
