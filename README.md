# kctx-manager

__A CLI tool to manage contexts in your $HOME/.kube/config or any other kubeconfig__

![Build and test](https://github.com/belitre/kctx-manager/workflows/Build%20and%20test/badge.svg)
![Publish release](https://github.com/belitre/kctx-manager/workflows/Publish%20release/badge.svg)

## Install

You can download kctx-manager binaries for windows, linux and mac from here: https://github.com/belitre/kctx-manager/releases

## Usage

__IMPORTANT:__ `kctx-manager` will use as default kubeconfig:

1. The kubeconfig file provided with the flag `-k` or `--kubeconfig`
2. If no kubeconfig is provided using the flags, it will use the kubeconfig from the environment variable `$KUBECONFIG`
3. If no kubeconfig is provided using the flags or the environment variable, it will use as default kubeconfig `$HOME/.kube/config`

```
Usage:
  kctx-manager [command]

Available Commands:
  add         Add the contexts defined in kubeconfig_file
  delete      Delete the context specified in context_name
  help        Help about any command
  list        Show current contexts
  rename      Rename current_context_name to new_context_name
  version     Print the version number of kctx-manager

Flags:
  -h, --help                help for kctx-manager
  -k, --kubeconfig string   Path to kubeconfig file, if not provided it will default to the value of $KUBECONFIG, and if the environment variable is not set it will default to $HOME/.kube/config

Use "kctx-manager [command] --help" for more information about a command.
```

### Add

`add` will add all the contexts found in the `kubeconfig_file` used as parameter to the default kubeconfig. __If a context with the same name is found in the default kubeconfig file, `kctx-manager` will override it.__

`add` will also change the name of the cluster and the name of the user before adding them to the kubeconfig file. Notice that it will only update the name, the rest of the information of each item will stay the same (server url, certificate, etc...). This is to avoid problems with auto generated kubeconfigs where you can have for example: for eks kubeconfigs with context and user with name `aws` and `kubernetes` for the cluster.

```
Usage:
  kctx-manager add kubeconfig_file [flags]

Flags:
  -h, --help          help for add
  -n, --name string   name of the cluster. Use this argument if you want to rename the cluster while adding it.

Global Flags:
  -k, --kubeconfig string   Path to kubeconfig file, if not provided it will default to the value of $KUBECONFIG, and if the environment variable is not set it will default to $HOME/.kube/config
```

```
$ kctx-manager add $HOME/Downloads/my_kubeconfig.yaml
Context coolcluster added/updated
Context patata added/updated
Context bobedilla added/updated
```

If the kubeconfig we want to add has only one context we can add it overriding the name using the flag `-n`:

```
$ kctx-manager add $HOME/Downloads/my_kubeconfig.yaml -n bleh
Context bleh added/updated
```

### Delete

`delete` will delete the context used as parameter from the default kubeconfig.

```
Usage:
  kctx-manager delete context_name [flags]

Flags:
  -h, --help   help for delete

Global Flags:
  -k, --kubeconfig string   Path to kubeconfig file, if not provided it will default to the value of $KUBECONFIG, and if the environment variable is not set it will default to $HOME/.kube/config
```

```
$ kctx-manager delete bobedilla
Context bobedilla deleted successfully!
```

### List

`list` will list all the contexts, with their endpoints, found in the default kubeconfig.

```
Usage:
  kctx-manager list [flags]

Flags:
  -h, --help   help for list

Global Flags:
  -k, --kubeconfig string   Path to kubeconfig file, if not provided it will default to the value of $KUBECONFIG, and if the environment variable is not set it will default to $HOME/.kube/config
```

```
$ kctx-manager list

 Context            Endpoint
 -------            --------
 minikube           https://192.168.64.5:8443
 patata             https://asd.yl4.eu-west-1.eks.amazonaws.com
 rancher            https://rancher.myhost.com/k8s/clusters/local
 bobedilla          https://bgd.sk1.eu-west-1.eks.amazonaws.com
 coolcluster        https://35.12.76.234
 docker-for-desktop https://localhost:6443
```

### Rename

`rename` will change the name to a context. __This will also change the name of the cluster and the name of the user in the kubeconfig file.__ Notice that it will only update the name, the rest of the information of each item will stay the same (endpoint url, certificates, etc...). This is to avoid problems with auto generated kubeconfigs where you can have for example: for eks kubeconfigs with context and user with name `aws` and `kubernetes` for the cluster.

If a context with the same name already exists `rename` will fail. To override the context with the same name use the flag `-f` or `--force`

```
Usage:
  kctx-manager rename current_context_name new_context_name [flags]

Flags:
  -f, --force   Forces rename. If new_context_name already exists it will be deleted.
  -h, --help    help for rename

Global Flags:
  -k, --kubeconfig string   Path to kubeconfig file, if not provided it will default to the value of $KUBECONFIG, and if the environment variable is not set it will default to $HOME/.kube/config
```

```
$ kctx-manager rename patata zorrocotropo -f
Context patata renamed to zorrocotropo successfully!
```
