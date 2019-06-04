# Knative Kafka Operator

The following will install Knative Kafka and configure it
appropriately for your cluster in the `default` namespace:

    kubectl apply -f deploy/crds/eventing_v1alpha1_kafka_install_crd.yaml
    kubectl apply -f deploy/
    kubectl apply -f deploy/crds/eventing_v1alpha1_kafka_install_cr.yaml

## Prerequisites

### Operator SDK

This operator was created using the
[operator-sdk](https://github.com/operator-framework/operator-sdk/).
It's not strictly required but does provide some handy tooling.

## The KnativeEventingKafka Custom Resource

The installation of Knative Kafka is triggered by the creation of
[an `KnativeEventingKafka` custom
resource](deploy/crds/eventing_v1alpha1_kafka_install_crd.yaml).

The following are all equivalent, but the latter may suffer from name
conflicts.

    kubectl get knativeventingkafka.eventing.knative.dev -oyaml
    kubectl get kek -oyaml
    kubectl get knativeventingkafka -oyaml

To uninstall Knative Kafka, simply delete the `KnativeEventingKafka` resource.

    kubectl delete kek --all

## Development

It can be convenient to run the operator outside of the cluster to
test changes. The following command will build the operator and use
your current "kube config" to connect to the cluster:

    operator-sdk up local

Pass `--help` for further details on the various `operator-sdk`
subcommands, and pass `--help` to the operator itself to see its
available options:

    operator-sdk up local --operator-flags "--help"

### Building the Operator Image

To build the operator,

    operator-sdk build quay.io/$REPO/knative-kafka-operator:$VERSION

The image should match what's in
[deploy/operator.yaml](deploy/operator.yaml) and the `$VERSION` should
match [version.go](version/version.go) and correspond to the contents
of [deploy/resources](deploy/resources/).

There is a handy script that will build and push an image to
[quay.io](https://quay.io/repository/openshift-knative/knative-kafka-operator)
and tag the source:

    ./hack/release.sh
	
## Operator Framework

The remaining sections only apply if you wish to create the metadata
required by the [Operator Lifecycle
Manager](https://github.com/operator-framework/operator-lifecycle-manager)

### Create a CatalogSource

The OLM requires special manifests that the operator-sdk can help
generate.

Create a `ClusterServiceVersion` for the version that corresponds to
the manifest[s] beneath [deploy/resources](deploy/resources/). The
`$PREVIOUS_VERSION` is the CSV yours will replace.

    operator-sdk olm-catalog gen-csv \
        --csv-version $VERSION \
        --from-version $PREVIOUS_VERSION \
        --update-crds

Most values should carry over, but if you're starting from scratch,
some post-editing of the file it generates may be required:

* Add fields to address any warnings it reports
* Verify `description` and `displayName` fields for all owned CRD's
* Set the `fieldPath` for `WATCH_NAMESPACE` to `metadata.annotations['olm.targetNamespaces']`

The [catalog.sh](hack/catalog.sh) script should yield a valid
`CatalogSource` for you to publish.

### Using OLM on Minikube

You can test the operator using
[minikube](https://kubernetes.io/docs/setup/minikube/) after
installing OLM on it:

    minikube start
    kubectl apply -f https://github.com/operator-framework/operator-lifecycle-manager/releases/download/0.9.0/olm.yaml

Once all the pods in the `olm` namespace are running, install the
operator like so:
    
    ./hack/catalog.sh | kubectl apply -n olm -f -

Interacting with OLM is possible using `kubectl` but the OKD console
is "friendlier". If you have docker installed, use [this
script](https://github.com/operator-framework/operator-lifecycle-manager/blob/master/scripts/run_console_local.sh)
to fire it up on <http://localhost:9000>.

#### Using kubectl

To install Knative Kafka into the `knative-eventing` namespace, apply
the following resources:

```
cat <<-EOF | kubectl apply -f -
---
apiVersion: v1
kind: Namespace
metadata:
  name: knative-eventing
---
apiVersion: operators.coreos.com/v1
kind: OperatorGroup
metadata:
  name: knative-eventing
  namespace: knative-eventing
---
apiVersion: operators.coreos.com/v1alpha1
kind: Subscription
metadata:
  name: knative-kafka-operator-sub
  generateName: knative-kafka-operator-
  namespace: knative-eventing
spec:
  source: knative-kafka-operator
  sourceNamespace: olm
  name: knative-kafka-operator
  channel: alpha
---
apiVersion: eventing.knative.dev/v1alpha1
kind: KnativeEventingKafka
metadata:
  name: knative-eventing-kafka
  namespace: knative-eventing-kafka
EOF
```
