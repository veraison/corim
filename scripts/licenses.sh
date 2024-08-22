#!/bin/bash
# Copyright 2024 Contributors to the Veraison project.
# SPDX-License-Identifier: Apache-2.0

set -e

type go-licenses &> /dev/null || go get github.com/google/go-licenses

MODULES+=("github.com/veraison/corim/corim")
MODULES+=("github.com/veraison/corim/comid")

for module in ${MODULES[@]}
do
  echo ">> retrieving licenses [ ${module} ]"
  go-licenses csv ${module}
done
