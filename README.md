# Install Orc8r Operator

## Pre-requisites

- Golang
- make
- Operator-sdk
- A running Kubernetes cluster

#### Golang

To install Golang use the below command:

```
$ sudo apt update
$ sudo apt install golang -y
$ wget https://go.dev/dl/go1.22.0.linux-amd64.tar.gz
$ sudo tar -C /usr/local -xzf go1.22.0.linux-amd64.tar.gz
$ echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
$ source ~/.bashrc
$ go version
```

To verify the Golang installation:

```
go version
```

#### make

To install make use the below command:

```
$ sudo apt update
$ sudo apt install make -y
```

To verify the make installation:

```
make --version
```

#### Operator-sdk

To install operator-sdk use the below command:

```
$ export ARCH=$(uname -m)
$ export OS=$(uname | awk '{print tolower($0)}')
$ export OPERATOR_SDK_DL_URL=https://github.com/operator-framework/operator-sdk/releases/latest/download
$ curl -LO ${OPERATOR_SDK_DL_URL}/operator-sdk_${OS}_${ARCH}
$ chmod +x operator-sdk_${OS}_${ARCH}
$ sudo mv operator-sdk_${OS}_${ARCH} /usr/local/bin/operator-sdk

```

To verify the operator-sdk installation:

```
operator-sdk version
```

#### Create the 'PMN' namespace

```
$ kubectl create namespace pmn
```

#### Before installing the Orc8r-operator , ensure that secrets are created for Orc8r

```
$ cd pmn-opertor/pmn-operator/configs/secrets

```

1. magmalte-mysql-secrets - it contains MYSQL Username and Password so edit them before applying.

kubectl apply -f pmn-operator/configs/secrets/magmalte-mysql-secrets.yaml -n pmn

```

```

2. orc8r-controller - it contains encoded string of DB username and password so edit them before applying.

kubectl apply -f pmn-operator/configs/secrets/orc8r-controller.yaml -n pmn

```

```

3. pmn-configs - it contains four files which needs to be mounted analytics.yml, elastic.yml, metricsd.yml and orchestrator.yml. (edit them accordingly)

kubectl apply -f pmn-operator/configs/secrets/pmn-configs.yaml -n pmn

```
4. pmn-envdir - this secret defines a comma-separated list of enabled services for the Orc8r-Operator.

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
