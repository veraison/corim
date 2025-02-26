#!/usr/bin/bash
# Copyright 2024 Contributors to the Veraison project.
# SPDX-License-Identifier: Apache-2.0
set -e

THIS_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

if [[ "$(type -p diag2cbor.rb)" == "" ]]; then
	echo "ERROR: please install ruby-cbor-diag package"
	exit 1
fi

for case in "$THIS_DIR"/src/*.diag; do
	outfile=$(basename "${case%%.diag}").cbor

	echo "generating $outfile"

	diag2cbor.rb "$case" > "$THIS_DIR/$outfile"
done

echo "done."
