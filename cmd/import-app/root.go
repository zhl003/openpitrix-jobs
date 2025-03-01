package main

import (
	"github.com/spf13/cobra"
	"io"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
	"kubesphere.io/openpitrix-jobs/pkg/client/clientset/versioned"
	"kubesphere.io/openpitrix-jobs/pkg/s3"
	"kubesphere.io/openpitrix-jobs/pkg/types"
	"kubesphere.io/openpitrix-jobs/pkg/utils"
)

var kubeconfig string
var master string
var versionedClient *versioned.Clientset
var k8sClient *kubernetes.Clientset
var s3Options *s3.Options

func newRootCmd(out io.Writer, args []string) (*cobra.Command, error) {
	s3Options = s3.NewS3Options()
	cmd := &cobra.Command{
		Use:          "import-app",
		Short:        "import builtin app into kubesphere",
		SilenceUsage: true,
	}

	cobra.OnInitialize(func() {
		utils.DumpConfig()

		ksConfig, err := types.TryLoadFromDisk()
		if err != nil {
			klog.Fatalf("load config failed, error: %s", err)
		}

		if ksConfig.OpenPitrixOptions == nil {
			klog.Fatalf("openpitrix config is empty, please wait a minute")
		}

		if ksConfig.OpenPitrixOptions.S3Options == nil {
			klog.Fatalf("s3 config is empty, please wait a minute")
		}

		s3Options = ksConfig.OpenPitrixOptions.S3Options

		config, err := clientcmd.BuildConfigFromFlags(master, kubeconfig)
		if err != nil {
			klog.Fatalf("load kubeconfig failed, error: %s", err)
		}
		versionedClient, err = versioned.NewForConfig(config)
		if err != nil {
			klog.Fatalf("build config failed, error: %s", err)
		}
		k8sClient = kubernetes.NewForConfigOrDie(config)
	})

	flags := cmd.PersistentFlags()

	addKlogFlags(flags)
	flags.StringVar(&kubeconfig, "kubeconfig", "", "path to the kubeconfig file")
	flags.StringVar(&master, "master", "", "kubernetes master")

	flags.Parse(args)
	cmd.AddCommand(
		newImportCmd(),
		newConvertCmd(),
	)

	return cmd, nil
}
