#!/usr/bin/bash
# Copyright 2024 Contributors to the Veraison project.
# SPDX-License-Identifier: Apache-2.0
set -e

GEN_TESTCASE=$(go env GOPATH)/bin/gen-testcase

if [[ ! -f ${GEN_TESTCASE} ]]; then
	echo "installing gen-testcase"

	go install github.com/veraison/gen-testcase@v0.0.1
fi

testcases=(
	psa-refval
	test-comid
	test-coswid
	test-cots
)

signed_testcases=(
	signed-corim-valid
	signed-corim-invalid
	signed-corim-valid-with-cots
)

for case in "${testcases[@]}"; do
	echo "generating ${case}.cbor"

	${GEN_TESTCASE} "src/${case}.yaml" -o "${case}.cbor"
done

for case in "${signed_testcases[@]}"; do
	echo "generating ${case}.cbor"

	${GEN_TESTCASE} -s ec-p256.jwk -m "src/${case}-meta.yaml" "src/${case}.yaml" \
		-o "${case}.cbor"
done

echo "done."
