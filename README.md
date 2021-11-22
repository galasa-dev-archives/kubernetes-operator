# galasa-kubernetes-operator
This is a alpha level release of the kubernetes operator, moved and re-architected from galasa-dev/extensions repository.

## Install



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
