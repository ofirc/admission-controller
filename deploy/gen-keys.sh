#!/usr/bin/env bash

# Copyright (c) 2024 Ofir Cohen.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# gen-keys.sh
#
# Sets up the environment for the admission controller webhook demo in the active cluster.

set -euo pipefail

# Read the PEM-encoded CA certificate, base64 encode it, and replace the `${CA_PEM_B64}` placeholder in the YAML
# template with it. Then, create the Kubernetes resources.
ca_pem_b64="$(openssl base64 -A <certs/ca.crt)"
webhook_key_pem_b64="$(openssl base64 -A <certs/webhook.key)"
webhook_crt_pem_b64="$(openssl base64 -A <certs/webhook.crt)"
sed \
  -e 's@${CA_PEM_B64}@'"$ca_pem_b64"'@g' \
  -e 's@${TLS_CRT_B64}@'"$webhook_crt_pem_b64"'@g' \
  -e 's@${TLS_KEY_B64}@'"$webhook_key_pem_b64"'@g' \
  < deploy/kubernetes/deployment.yaml.template \
  > deploy/kubernetes/deployment.yaml
