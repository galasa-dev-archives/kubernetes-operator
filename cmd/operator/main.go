package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/galasa-dev/galasa-kubernetes-operator/pkg/reconciler/api"
	"github.com/galasa-dev/galasa-kubernetes-operator/pkg/reconciler/cps"
	"github.com/galasa-dev/galasa-kubernetes-operator/pkg/reconciler/enginecontroller"
	"github.com/galasa-dev/galasa-kubernetes-operator/pkg/reconciler/galasaecosystem"
	"github.com/galasa-dev/galasa-kubernetes-operator/pkg/reconciler/metrics"
	"github.com/galasa-dev/galasa-kubernetes-operator/pkg/reconciler/ras"
	"github.com/galasa-dev/galasa-kubernetes-operator/pkg/reconciler/resmon"
	"github.com/galasa-dev/galasa-kubernetes-operator/pkg/reconciler/toolbox"

	"knative.dev/pkg/injection"
	"knative.dev/pkg/injection/sharedmain"
	"knative.dev/pkg/signals"
)

const (
	// ControllerLogKey is the name of the logger for the controller cmd
	ControllerLogKey = "galasa-ecosystem-operator"
)

var (
	namespace = flag.String("namespace", "default", "Namespace to restrict informer to. Optional, defaults to default.")
)

func main() {
	cfg := injection.ParseAndGetRESTConfigOrDie()
	ctx := injection.WithNamespaceScope(signals.NewContext(), *namespace)

	// Set up liveness and readiness probes.
	mux := http.NewServeMux()

	mux.HandleFunc("/", handler)
	mux.HandleFunc("/health", handler)
	mux.HandleFunc("/readiness", handler)

	port := os.Getenv("PROBES_PORT")
	if port == "" {
		port = "8080"
	}
	go func() {
		log.Printf("Readiness and health check server listening on port %s", port)
		log.Fatal(http.ListenAndServe(":"+port, mux))
	}()

	sharedmain.MainWithConfig(ctx, ControllerLogKey, cfg,
		galasaecosystem.NewController(*namespace),
		cps.NewController(*namespace),
		ras.NewController(*namespace),
		api.NewController(*namespace),
		enginecontroller.NewController(*namespace),
		metrics.NewController(*namespace),
		resmon.NewController(*namespace),
		toolbox.NewController(*namespace),
	)
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
