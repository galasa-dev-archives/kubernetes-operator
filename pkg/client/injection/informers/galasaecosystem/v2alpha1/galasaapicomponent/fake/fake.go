/*
 * Copyright contributors to the Galasa Project
 */
// Code generated by injection-gen. DO NOT EDIT.

package fake

import (
	context "context"

	fake "github.com/galasa-dev/galasa-kubernetes-operator/pkg/client/injection/informers/factory/fake"
	galasaapicomponent "github.com/galasa-dev/galasa-kubernetes-operator/pkg/client/injection/informers/galasaecosystem/v2alpha1/galasaapicomponent"
	controller "knative.dev/pkg/controller"
	injection "knative.dev/pkg/injection"
)

var Get = galasaapicomponent.Get

func init() {
	injection.Fake.RegisterInformer(withInformer)
}

func withInformer(ctx context.Context) (context.Context, controller.Informer) {
	f := fake.Get(ctx)
	inf := f.Galasa().V2alpha1().GalasaApiComponents()
	return context.WithValue(ctx, galasaapicomponent.Key{}, inf), inf.Informer()
}
