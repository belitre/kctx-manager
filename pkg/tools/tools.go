package tools

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/belitre/kctx-manager/pkg/kubeconfig"
)

func PrintContexts(contexts []*kubeconfig.ContextWithEndpoint) {
	if len(contexts) == 0 {
		fmt.Println("No contexts found.")
		return
	}

	w := new(tabwriter.Writer)

	w.Init(os.Stdout, 15, 8, 1, '\t', 0)
	defer w.Flush()

	fmt.Fprintf(w, "\n %s\t%s\t", "Context", "Endpoint")
	fmt.Fprintf(w, "\n %s\t%s\t", "-------", "--------")

	for _, v := range contexts {
		fmt.Fprintf(w, "\n %s\t%s\t", v.Name, v.Endpoint)
	}
}
