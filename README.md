# Install Orc8r Operator

## Pre-requisites

- Helm
- A running Kubernetes cluster

#### Helm

To install Helm use the below command:

```
$ curl -fsSL https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash
```

To verify the Helm installation:

```
$ helm version
```

#### Create the 'PMN' namespace

```
$ kubectl create namespace pmn
```

#### Before installing the Orc8r-operator , ensure that secrets are created for Orc8r

```
$ cd pmn-opertor/pmn-operator/configs/secrets

kubectl apply -f pmn-operator/configs/secrets/magmalte-mysql-secrets.yaml -n pmn

kubectl apply -f pmn-operator/configs/secrets/orc8r-controller.yaml -n pmn

kubectl apply -f pmn-operator/configs/secrets/pmn-configs.yaml -n pmn

kubectl apply -f pmn-operator/configs/secrets/pmn-envdir.yaml -n pmn
```

#### Modify the values of pmnsystems_v1alpha1_pmnsystem.yaml (CR) before installing the Operator:

```
Edit the values in

$ cd pmn-operator/config/samples

$ vim pmn-operator/config/samples/pmnsystems_v1alpha1_pmnsystem.yaml
```

## After editing the values of the CR the most important step is to apply the CRD:

```
$ cd pmn-operator/config/crd/bases

$ kubectl apply -f pmnsystems.pmnsystem.com_pmnsystems.yaml
```

### To Uninstall the Orc8r-operator

$ cd pmn-operator

i. run command "make generate"
ii. run command "make manifests"
iii. run command "make run"

#### To uninstall the Orc8r-operator:

```
$ run command "make uninstall"
```
