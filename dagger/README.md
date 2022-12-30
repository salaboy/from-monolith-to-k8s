# Using Dagger for Service and Environment Pipelines


This pipeline will clone the specified `repoUrl` and perform a `helmfile sync` using the
default kube context in the supplied `kubeCfgFilePath` file.

```
go run env_pipeline.go <repoUrl> <kubeCfgFilePath>
```
