package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	infrav1 "github.com/flux-iac/tofu-controller/api/v1alpha2"
	"github.com/flux-iac/tofu-controller/internal/web"
	sourcev1 "github.com/fluxcd/source-controller/api/v1"
	flag "github.com/spf13/pflag"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cluster"
)

var scheme = runtime.NewScheme()

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(sourcev1.AddToScheme(scheme))
	utilruntime.Must(infrav1.AddToScheme(scheme))
}

func main() {
	var port int
	var kubeconfig string

	flag.IntVar(&port, "port", 9080, "Port to listen on.")
	flag.StringVar(&kubeconfig, "kubeconfig", "", "Path to kubeconfig file. Defaults to in-cluster config.")
	flag.Parse()

	log := ctrl.Log.WithName("web-server")

	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	if kubeconfig != "" {
		loadingRules.ExplicitPath = kubeconfig
	}

	cfg, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		loadingRules,
		&clientcmd.ConfigOverrides{},
	).ClientConfig()
	if err != nil {
		log.Error(err, "unable to load kubeconfig")
		os.Exit(1)
	}

	c, err := cluster.New(cfg, func(o *cluster.Options) {
		o.Scheme = scheme
	})
	if err != nil {
		log.Error(err, "unable to create cluster client")
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Start the cluster (runs informers) in background
	go func() {
		if err := c.Start(ctx); err != nil {
			log.Error(err, "cluster stopped with error")
		}
	}()

	log.Info("Starting web server", "port", port)
	if err := web.RunServer(ctx, c, port); err != nil {
		log.Error(err, "web server error")
		os.Exit(1)
	}
}
