#!/usr/bin/bash
# Copyright 2024-2026 Contributors to the Veraison project.
# SPDX-License-Identifier: Apache-2.0
set -e

THIS_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

GEN_TESTCASE=$(go env GOPATH)/bin/gen-testcase

if [[ "$(type -p diag2cbor.rb)" == "" ]]; then
	echo "ERROR: please install ruby-cbor-diag package"
	exit 1
fi

if [[ ! -f ${GEN_TESTCASE} ]]; then
	echo "installing gen-testcase"
	go install github.com/veraison/gen-testcase@v0.0.3
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

for case in "$THIS_DIR"/src/*.diag; do
	outfile=$(basename "${case%%.diag}").cbor

	echo "generating $outfile"

	diag2cbor.rb "$case" > "$THIS_DIR/$outfile"
done

echo "done."
