# Copyright 2021 Richard Kosegi
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.


IMAGE_NAME := "rkosegi/cert-manager-webhook-active24"
IMAGE_TAG := "1.0.0"

OUT := $(shell pwd)/deploy

$(shell mkdir -p "$(OUT)")

dist-clean:
	rm -fr .cache kubebuilder cert-manager-webhook-active24

fetch-test-binaries:
	mkdir .cache || true
	test -f .cache/kubebuilder-tools-1.15.5-linux-amd64.tar.gz || \
		curl https://storage.googleapis.com/kubebuilder-tools/kubebuilder-tools-1.15.5-linux-amd64.tar.gz -o .cache/kubebuilder-tools-1.15.5-linux-amd64.tar.gz
	tar -zvxf .cache/kubebuilder-tools-1.15.5-linux-amd64.tar.gz


verify:
	TEST_ZONE_NAME=mydomain.tld. go test -v .

build:
	docker build -t "$(IMAGE_NAME):$(IMAGE_TAG)" .

.PHONY: rendered-manifest.yaml
rendered-manifest.yaml:
	helm template \
	    cert-manager-webhook-active24 \
        --set image.repository=$(IMAGE_NAME) \
        --set image.tag=$(IMAGE_TAG) \
		--namespace cert-manager \
        deploy/chart > "$(OUT)/rendered-manifest.yaml"
