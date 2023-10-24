#!/usr/bin/env bash

# Copyright 2017 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -o errexit
set -o nounset
set -o pipefail

# corresponding to go mod init <module>
MODULE=igitlab.lenovo.com/lecp/mec/bailu
# api package
APIS_PKG=api
# generated output package
OUTPUT_PKG=client
# group-version such as foo:v1alpha1
GROUP_VERSION=cloud.edge:v1

SCRIPT_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
CODEGEN_PKG=${CODEGEN_PKG:-$(cd "${SCRIPT_ROOT}"; ls -d -1 ./vendor/k8s.io/code-generator 2>/dev/null || echo ../code-generator)}

# generate the code with:
# --output-base    because this script should also be able to run inside the vendor dir of
#                  k8s.io/kubernetes. The output-base is needed for the generators to output into the vendor dir
#                  instead of the $GOPATH directly. For normal projects this can be dropped.
bash "${CODEGEN_PKG}"/generate-groups.sh "client,lister,informer" \
  ${MODULE}/${OUTPUT_PKG} ${MODULE}/${APIS_PKG} \
  ${GROUP_VERSION} \
    --output-base /tmp \
  --go-header-file "${SCRIPT_ROOT}"/hack/boilerplate.go.txt

rm -rf ${SCRIPT_ROOT}/${OUTPUT_PKG}
mkdir -p ${SCRIPT_ROOT}/${OUTPUT_PKG}
mv /tmp/${MODULE}/${OUTPUT_PKG}/* ${SCRIPT_ROOT}/${OUTPUT_PKG}

# To use your own boilerplate text append:
#   --go-header-file "${SCRIPT_ROOT}"/hack/custom-boilerplate.go.txt
