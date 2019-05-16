# kctx-manager

__A CLI tool to manage contexts in your $HOME/.kube/config or any other kubeconfig__

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

```
Usage:
  kctx-manager add kubeconfig_file [flags]

Flags:
  -h, --help   help for add

Global Flags:
  -k, --kubeconfig string   Path to kubeconfig file, if not provided it will default to the value of $KUBECONFIG, and if the environment variable is not set it will default to $HOME/.kube/config
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

### Rename

`rename` will change the name to a context. __This will only change the name of the context! It won't change the name of the cluster or the user in the kubeconfig file.__

```
Usage:
  kctx-manager rename current_context_name new_context_name [flags]

Flags:
  -h, --help   help for rename

Global Flags:
  -k, --kubeconfig string   Path to kubeconfig file, if not provided it will default to the value of $KUBECONFIG, and if the environment variable is not set it will default to $HOME/.kube/config
```
