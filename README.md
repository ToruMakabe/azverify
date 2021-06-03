# [Under developing] Azure Resource Verifier with Resource Graph

## Table of Contents

- [About](#about)
- [Getting Started](#getting_started)
- [Usage](#usage)
- [Limitations](#limitations)
## About <a name = "about"></a>

You can verify if there is a difference between your desired properties and actual with this CLI.
This CLI read your desired properties as JSON files, and query to Azure Resource Graph API, then check the difference.

[Azure Resource Graph table and resource type reference](https://docs.microsoft.com/en-us/azure/governance/resource-graph/reference/supported-tables-resources)

## Getting Started <a name = "getting_started"></a>

### Prerequisites

* Go 1.16 or later (if you will build)
* Azure Resource Graph [permissions](https://docs.microsoft.com/en-us/azure/governance/resource-graph/overview#permissions-in-azure-resource-graph)

### Installing

```
go install github.com/ToruMakabe/azverify@latest
```

Or download [binary](https://github.com/ToruMakabe/azverify/releases)

## Usage <a name = "usage"></a>

### Global
```
  azverify [command]

Available Commands:
  help        Help about any command
  match       Matching check the difference between desired and actual
  version     Display the version

Flags:
      --cert-password string     cert file password
      --cert-path string         PKCS12 (.pfx) cert file path
      --client-id string         Azure AD service principal App ID
      --client-secret string     Azure AD service principal App secret
      --config string            config file path (default "$HOME/.azverify/config.toml")
      --env-prefix string        env prefix (default "AZV")
      --environment string       Azure environment ([public]/usgovernment/german/china)
  -h, --help                     help for azverify
      --log-level string         log level (DEBUG/[INFO]/ERROR)
      --subscription-id string   Azure subscription ID
      --tenant-id string         Azure AD tenant ID
```

#### Config options and the evaluation order

Each item takes precedence over the item below it.

* flag
* env. var
  * default prefix is "AZV_". e.g. AZV_TENANT_ID
* config file
  * [sample](https://github.com/ToruMakabe/azverify/blob/main/config.toml)

#### Auth methods and the evaluation order

It is determined by elements set with flags, env variables and config file.

* Service Principal Client Certificate
  * tenant-id, subscription-id, cert-path, cert-password
* Service Principal Client Secret
  * tenant-id, subscription-id, client-id, client-secret
* Managed Identiy
  * --use-msi flag in subcommand
* Azure CLI token
  * just run without auth options on the machine has valid Azure CLI token

### subcommand [match]
```
  azverify match [flags]

Flags:
  -f, --file string   path(glob) of the file(s) where the desired resources are written
  -h, --help          help for match
      --use-msi       flag for using Managed Identity to auth (defalut: false)
```

You have to prepare your "desired" resources as JSON(array) file. "id" is mandatry to identify target resource, but do not have to describe all keys/values. You can write only keys/values you want to verify. [Sample](https://github.com/ToruMakabe/azverify/blob/main/testdata/desired/desired_1.json) is in [testdata](https://github.com/ToruMakabe/azverify/tree/main/testdata/desired).

Then, run match command.

```
azverify match -f ./desired.json
```

If the result of matching desired resources against to the actual returns from Resource Graph API is [SupersetMatch](https://github.com/ToruMakabe/azverify/tree/main/testdata/actual) or FullMatch, result of the match command will be successful. On the other hand, if there was any other matching result, exits with status code 1. If all matching(s) was successful, exit code is 0.

## Limitations <a name = "limitations"></a>

* Currently do not have a subcommand to generate template of desired resource by resource types. Investigating the feasibility.
