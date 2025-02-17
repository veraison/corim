#!/usr/bin/bash
# Copyright 2024 Contributors to the Veraison project.
# SPDX-License-Identifier: Apache-2.0
set -e

GEN_TESTCASE=$(go env GOPATH)/bin/gen-testcase

if [[ ! -f ${GEN_TESTCASE} ]]; then
	echo "installing gen-testcase"
	go install github.com/veraison/gen-testcase@v0.0.2
fi

testcases=(
	good-corim
	example-corim
	corim-with-extensions
)

for case in "${testcases[@]}"; do
	echo "generating unsigned-${case}.cbor"

	${GEN_TESTCASE} "src/${case}.yaml" -o "unsigned-${case}.cbor"

	echo "generating signed-${case}.cbor"

	${GEN_TESTCASE} -s src/ec-p256.jwk -m src/meta.yaml "src/${case}.yaml" \
		-o "signed-${case}.cbor"
done

echo "done."
