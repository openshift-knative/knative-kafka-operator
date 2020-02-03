#!/usr/bin/env bash

DIR=${DIR:-$(cd $(dirname "$0")/.. && pwd)}
NAME=${NAME:-$(ls $DIR/deploy/olm-catalog)}

x=( $(echo $NAME | tr '-' ' ') )
DISPLAYNAME=${DISPLAYNAME:=${x[*]^}}

LATEST=$(find $DIR/deploy/olm-catalog -name '*version.yaml' | sort -n | sed "s/^.*\/\([^/]..*\).clusterserviceversion.yaml$/\1/" | tail -1)

indent() {
  INDENT="      "
  sed "s/^/$INDENT/" | sed "1 s/^${INDENT}\(.*\)/${INDENT:0:-2}- \1/"
}

rm -rf $DIR/.crds
mkdir $DIR/.crds
find $DIR/deploy/olm-catalog -name '*_crd.yaml' | sort -n | xargs -I{} cp {} $DIR/.crds/

CRD=$(for i in $(ls $DIR/.crds/*); do cat $i | grep -v -- "---" | indent; done)
CSV=$(for i in $(find $DIR/deploy/olm-catalog -name '*version.yaml' | sort -n); do cat $i | indent; done)
PKG=$(for i in $DIR/deploy/olm-catalog/$NAME/*package.yaml; do cat $i | indent; done)

cat <<EOF | sed 's/^  *$//'
kind: ConfigMap
apiVersion: v1
metadata:
  name: $NAME

data:
  customResourceDefinitions: |-
$CRD
  clusterServiceVersions: |-
$CSV
  packages: |-
$PKG
---
apiVersion: operators.coreos.com/v1alpha1
kind: CatalogSource
metadata:
  name: $NAME
spec:
  configMap: $NAME
  displayName: $DISPLAYNAME
  publisher: Red Hat
  sourceType: internal
EOF
