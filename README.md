# Falcosidekick + Tekton

In this demo I won't explain how tekton works, there is great material on how to get started using tekton,
for example the [official docs](https://tekton.dev/docs/overview/).

I have taken lots and lots of inspiration from this awesome blog post
[Falcosidekick + OpenFaas = a Kubernetes Response Engine, Part 2](https://falco.org/blog/falcosidekick-openfaas/).
I will more or less copy Batuhan Apaydın but do it with tekton.

You can find all the pure yaml in this repo. TODO add link to the repo.

## Minikube

I'm sure you can use kind as well but falcosidekick complained a bit when i tried and I was to lazy to check out what extra flags i need to start in my kind cluster.

```shell
minikube start --cpus 3 --memory 8192 --vm-driver virtualbox
```

## Install Tekton

Install tekton quickly by using pipelines and triggers. From a operations point of view I would use the tekton operator.

```shell
kubectl apply --filename https://storage.googleapis.com/tekton-releases/pipeline/latest/release.yaml
kubectl apply --filename https://storage.googleapis.com/tekton-releases/triggers/latest/release.yaml

```

If this is your first time using tekton I would recomend you jump in to the tekton triggers [getting-started guide](https://github.com/tektoncd/triggers/tree/v0.10.1/docs/getting-started) to see how it works.

## Install Falco + Falcosidekick

```shell
kubectl create namespace falco
helm repo add falcosecurity https://falcosecurity.github.io/charts
helm repo update

cat <<'EOF' >> values.yaml
falcosidekick:
  config:
    webhook:
      address: http://el-falco-listener:8080
      customHeaders: |
        Falcon:true
  enabled: true


customRules:
  # Applications which are expected to communicate with the Kubernetes API
  rules_user_known_k8s_api_callers.yaml: |-
    - macro: user_known_contact_k8s_api_server_activities
      condition: >
        (container.image.repository = "gcr.io/tekton-releases/github.com/tektoncd/triggers/cmd/eventlistenersink")
EOF

helm upgrade --install falco falcosecurity/falco --namespace falco -f values.yaml
```

We need to setup a custom rule for event-listener since tekton and the event listener talks allot to the kubernetes API.

## Configure tekton

```shell
kubectl create ns falcoresponse
cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: ServiceAccount
metadata:
  name: falco-pod-delete
  namespace: falcoresponse
---
kind: ClusterRole
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: falco-pod-delete-cluster-role
rules:
  - apiGroups: [""]
    resources: ["pods"]
    verbs: ["get", "list", "delete"]
---
kind: ClusterRoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: falco-pod-delete-cluster-role-binding
roleRef:
  kind: ClusterRole
  name: falco-pod-delete-cluster-role
  apiGroup: rbac.authorization.k8s.io
subjects:
  - kind: ServiceAccount
    name: falco-pod-delete
    namespace: falcoresponse
EOF
```

### EventListener

Notice the Cel header match=Falcon
