package credentials

import (
	"flag"
	"fmt"
	"path/filepath"

	"github.com/nlamot/sofibot/pkg/mmp/credentials"
	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all mmp credentials in a namespace",
	Long:  `List all mmp credentials in a namespace`,
	Run: func(cmd *cobra.Command, args []string) {
		var kubeconfig *string
		if home := homedir.HomeDir(); home != "" {
			kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
		} else {
			kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
		}
		flag.Parse()

		// use the current context in kubeconfig
		config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			panic(err.Error())
		}
		r, _ := credentials.List(config)
		for _, element := range r {
			fmt.Println(element)
		}
	},
}
