## GitLab Token Injector
This tool injects a JWT token into a secret variable of all GitLab projects that the owner of the GitLab API token is a member of.
The GitLab user for API token needs to have the `Master` role for a project to be able to use the secret variables API.
The secret variable has the default name `K8S_TOKEN` and the generated JWT token has a default lifetime of `48h`.

### Build
```
glide install
go build
```

### Usage
```
gitlab-token-injector -host git.example.com -token MyAPIToken -key /path/to/private.key
```

The command line arguments are:
```
Required:

  -host=""          GitLab instance hostname
  -key=""           Path to a private key file
  -token=""         GitLab API token

Optional:

  -debug=false      Enable debug output
  -config="":       Path to a config file
  -ttl=48h0m0s      Token lifetime in hours
  -var="KUBE_TOKEN" GitLab secret variable

```

Instead of using arguments, you can also use environment variables (using upppercase argument names) or a ini-style config file.
