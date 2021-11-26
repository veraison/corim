#!/bin/bash

set -e

type go-licenses &> /dev/null || go get github.com/google/go-licenses

MODULES+=("github.com/veraison/corim/corim")
MODULES+=("github.com/veraison/corim/comid")
MODULES+=("github.com/veraison/corim/cocli")

for module in ${MODULES[@]}
do
  echo ">> retrieving licenses [ ${module} ]"
  go-licenses csv ${module}
done
