#!/bin/bash
# Copyright 2026 Contributors to the Veraison project.
# SPDX-License-Identifier: Apache-2.0
#
# Generates unsigned CoRIM test cases for CCA Platform profile validation testing.
# Unlike the main testcases, we only generate unsigned CoRIMs since
# profile validation is independent of signing.
set -e

THIS_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

GEN_TESTCASE=$(go env GOPATH)/bin/gen-testcase

if [[ ! -f ${GEN_TESTCASE} ]]; then
	echo "installing gen-testcase"
	go install github.com/veraison/gen-testcase@v0.0.3
fi

testcases=(
	cca-platform-valid
	cca-platform-invalid-impl-id
	cca-platform-invalid-digests
	cca-realm-valid
	cca-realm-invalid-rim-size
	cca-realm-invalid-missing-rim
	no-profile
)

for case in "${testcases[@]}"; do
	echo "generating ${case}.cbor"

	${GEN_TESTCASE} "${THIS_DIR}/src/${case}.yaml" -o "${THIS_DIR}/${case}.cbor"
done

echo "done."
