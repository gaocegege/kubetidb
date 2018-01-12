package main

import (
	"flag"
	"os"
	"runtime"
	"time"

	"github.com/golang/glog"
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	clientset "github.com/gaocegege/kubetidb/pkg/clientset/versioned"
	"github.com/gaocegege/kubetidb/pkg/controller"
	tidbInformers "github.com/gaocegege/kubetidb/pkg/informers/externalversions"
	"github.com/gaocegege/kubetidb/pkg/util/signals"
	"github.com/gaocegege/kubetidb/pkg/version"
)

var (
	masterURL    string
	printVersion bool
	kubeconfig   string
)

func run() {
	// set up signals so we handle the first shutdown signal gracefully
	stopCh := signals.SetupSignalHandler()

	cfg, err := clientcmd.BuildConfigFromFlags(masterURL, kubeconfig)
	if err != nil {
		glog.Fatalf("Error building kubeconfig: %s", err.Error())
	}

	kubeClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		glog.Fatalf("Error building kubernetes clientset: %s", err.Error())
	}

	tidbClient, err := clientset.NewForConfig(cfg)
	if err != nil {
		glog.Fatalf("Error building example clientset: %s", err.Error())
	}

	kubeInformerFactory := kubeinformers.NewSharedInformerFactory(kubeClient, time.Second*30)
	tidbInformerFactory := tidbInformers.NewSharedInformerFactory(tidbClient, time.Second*30)

	controller := controller.NewController(kubeClient, tidbClient, kubeInformerFactory, tidbInformerFactory)

	go kubeInformerFactory.Start(stopCh)
	go tidbInformerFactory.Start(stopCh)

	if err = controller.Run(2, stopCh); err != nil {
		glog.Fatalf("Error running controller: %s", err.Error())
	}

}

func init() {
	flag.StringVar(&kubeconfig, "kubeconfig", "", "Path to a kubeconfig. Only required if out-of-cluster.")
	flag.StringVar(&masterURL, "master", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")
	flag.BoolVar(&printVersion, "version", false, "Show version and quit")
}

func main() {
	// This is to solve https://github.com/golang/glog/commit/65d674618f712aa808a7d0104131b9206fc3d5ad, which is definitely NOT cool.
	flag.Parse()

	glog.Infof("kubetidb Version: %v", version.Version)
	glog.Infof("Git SHA: %s", version.GitSHA)
	glog.Infof("Go Version: %s", runtime.Version())
	glog.Infof("Go OS/Arch: %s/%s", runtime.GOOS, runtime.GOARCH)
	if printVersion {
		os.Exit(0)
	}
	run()
}
