Noddy 
# Galasa Kubernetes operator
This is a alpha level release of the kubernetes operator, moved and re-architected from galasa-dev/extensions repository.

If any problems are found please open issues at https://github.com/galasa-dev/projectmanagement/issues

## Operator Install

For a basic install that creates a namespace called `galasa` and installs the operator and relevant ecosystem CRD's use:

```
kubectl apply -f https://raw.githubusercontent.com/galasa-dev/galasa-kubernetes-operator/main/releases/0.18.1/release.yaml
```
This limits all work to the galasa namespace

If you wish to install the operator and ecosystem in another namespace, please amend the yaml as neccessary.

## Ecosystem Install

For installing using the default values for your cluster you can simply apply a yaml like the one found in the `examples/basic.yaml`. This will use default values for attributes e.g storageClassName

Certain values can be overridden, please see `examples/ecosystem_with_overrides.yaml` for details

Once the ecosystem have been applied, you can see the state of the ecosystem with a `kubectl get galasaecosystem` which should show if the ecosystem is ready for work load and the bootstrapURL:
```
NAME               CREATED                READY   BOOTSTRAPURL
Galasa-Ecosystem   2021-11-22T10:12:21Z   true    http://example.com:30000/bootstrap
```
## Development

#### Code Generation
```
# Deepcopy and Client gen
hack/generate-groups.sh all github.com/galasa-dev/galasa-kubernetes-operator/pkg/client github.com/galasa-dev/galasa-kubernetes-operator/pkg/apis "galasaecosystem:v2alpha1" -h hack/boilerplate/boilerplate.go.txt

# Knative Injection clients
hack/generate-knative.sh injection github.com/galasa-dev/galasa-kubernetes-operator/pkg/client github.com/galasa-dev/galasa-kubernetes-operator/pkg/apis "galasaecosystem:v2alpha1" -h hack/boilerplate/boilerplate.go.txt
```

#### Deploying with Ko:
```
# Build and deploy
ko apply -f config
```
