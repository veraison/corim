#!/bin/bash
# SPDX-License-Identifier: Apache-2.0
set -eu
set -o pipefail

ROOT_CERT_NAME=rootCA
INTERMEDIATE_CERT_NAME=intermediateCA
END_ENTITY_CERT_NAME=endEntity

THIS_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
MISC_DIR="$THIS_DIR/../testdata"

mkdir -p "$MISC_DIR"

trap '[[ $_should_clean_certs_artifacts == true ]] && clean_certs_artifacts' EXIT

function write_intermediate_extfile() {
    cat > "${MISC_DIR}/${INTERMEDIATE_CERT_NAME}.ext" <<'EOF'
[ v3_intermediate ]
basicConstraints = critical, CA:TRUE, pathlen:0
keyUsage = critical, keyCertSign, cRLSign
subjectKeyIdentifier = hash
authorityKeyIdentifier = keyid,issuer
EOF
}

function write_end_entity_extfile() {
    cat > "${MISC_DIR}/${END_ENTITY_CERT_NAME}.ext" <<'EOF'
[ v3_end_entity ]
basicConstraints = critical, CA:FALSE
keyUsage = critical, digitalSignature
subjectKeyIdentifier = hash
authorityKeyIdentifier = keyid,issuer
EOF
}

function create_root_cert() {
    _check_openssl

    if [[ -f "${MISC_DIR}/${ROOT_CERT_NAME}.der" ]] && [[ "${FORCE:-}" != "1" ]]; then
        echo "Root certificate already exists. Skipping creation."
        return
    fi

    if [[ ! -f "${MISC_DIR}/${ROOT_CERT_NAME}.key" ]]; then
        openssl ecparam -name prime256v1 -genkey -noout -out ${MISC_DIR}/${ROOT_CERT_NAME}.key
    fi

    openssl req -x509 -new -nodes -key ${MISC_DIR}/${ROOT_CERT_NAME}.key \
        -sha256 -days 3650 \
        -subj "/CN=Acme Inc." \
        -addext "basicConstraints=critical,CA:TRUE" \
        -addext "keyUsage=critical,keyCertSign,cRLSign" \
        -addext "subjectKeyIdentifier=hash" \
        -out ${MISC_DIR}/${ROOT_CERT_NAME}.crt
    openssl x509 -in ${MISC_DIR}/${ROOT_CERT_NAME}.crt -outform der \
        -out ${MISC_DIR}/${ROOT_CERT_NAME}.der
    rm -f ${MISC_DIR}/${ROOT_CERT_NAME}.crt

    echo "Created ${MISC_DIR}/${ROOT_CERT_NAME}.der and ${MISC_DIR}/${ROOT_CERT_NAME}.key"
}

function create_intermediate_cert() {
    _check_openssl
    _check_root_cert

    if [[ -f "${MISC_DIR}/${INTERMEDIATE_CERT_NAME}.der" ]] && [[ "${FORCE:-}" != "1" ]]; then
        echo "Intermediate certificate already exists. Skipping creation."
        return
    fi

    if [[ ! -f "${MISC_DIR}/${INTERMEDIATE_CERT_NAME}.key" ]]; then
        openssl ecparam -name prime256v1 -genkey -noout -out ${MISC_DIR}/${INTERMEDIATE_CERT_NAME}.key
    fi

    write_intermediate_extfile
    openssl req -new -key ${MISC_DIR}/${INTERMEDIATE_CERT_NAME}.key \
        -out ${MISC_DIR}/${INTERMEDIATE_CERT_NAME}.csr \
        -subj "/CN=Acme Gizmos"
    openssl x509 -req -in ${MISC_DIR}/${INTERMEDIATE_CERT_NAME}.csr \
        -CA ${MISC_DIR}/${ROOT_CERT_NAME}.der \
        -CAkey ${MISC_DIR}/${ROOT_CERT_NAME}.key \
        -CAcreateserial \
        -out ${MISC_DIR}/${INTERMEDIATE_CERT_NAME}.crt \
        -days 3650 -sha256 \
        -CAform der \
        -extensions v3_intermediate \
        -extfile ${MISC_DIR}/${INTERMEDIATE_CERT_NAME}.ext
    openssl x509 -in ${MISC_DIR}/${INTERMEDIATE_CERT_NAME}.crt \
        -outform der -out ${MISC_DIR}/${INTERMEDIATE_CERT_NAME}.der
    rm -f ${MISC_DIR}/${INTERMEDIATE_CERT_NAME}.crt ${MISC_DIR}/${INTERMEDIATE_CERT_NAME}.ext

    echo "Created ${MISC_DIR}/${INTERMEDIATE_CERT_NAME}.der and ${MISC_DIR}/${INTERMEDIATE_CERT_NAME}.key"
}

function create_end_entity_cert() {
    _check_openssl
    _check_root_cert

    if [[ ! -f "${MISC_DIR}/${INTERMEDIATE_CERT_NAME}.der" ]]; then
        create_intermediate_cert
    fi

    if ([[ -f "${MISC_DIR}/${END_ENTITY_CERT_NAME}.der" ]] && [[ -f "${MISC_DIR}/${END_ENTITY_CERT_NAME}.key" ]] && [[ "${FORCE:-}" != "1" ]]); then
        echo "End-entity certificate and key already exist. Skipping creation."
        return
    fi

    if [[ ! -f "${MISC_DIR}/${END_ENTITY_CERT_NAME}.key" ]]; then
        openssl ecparam -name prime256v1 -genkey -noout -out ${MISC_DIR}/${END_ENTITY_CERT_NAME}.key
    fi

    write_end_entity_extfile
    openssl req -new -key ${MISC_DIR}/${END_ENTITY_CERT_NAME}.key \
        -out ${MISC_DIR}/${END_ENTITY_CERT_NAME}.csr \
        -subj "/CN=Acme Gizmo CoRIM signer"
    openssl x509 -req -in ${MISC_DIR}/${END_ENTITY_CERT_NAME}.csr \
        -CA ${MISC_DIR}/${INTERMEDIATE_CERT_NAME}.der \
        -CAkey ${MISC_DIR}/${INTERMEDIATE_CERT_NAME}.key \
        -CAcreateserial \
        -out ${MISC_DIR}/${END_ENTITY_CERT_NAME}.crt \
        -days 1825 -sha256 \
        -CAform der \
        -extensions v3_end_entity \
        -extfile ${MISC_DIR}/${END_ENTITY_CERT_NAME}.ext
    openssl x509 -in ${MISC_DIR}/${END_ENTITY_CERT_NAME}.crt \
        -outform der -out ${MISC_DIR}/${END_ENTITY_CERT_NAME}.der
    rm -f ${MISC_DIR}/${END_ENTITY_CERT_NAME}.crt ${MISC_DIR}/${END_ENTITY_CERT_NAME}.ext

    echo "Created ${MISC_DIR}/${END_ENTITY_CERT_NAME}.der and ${MISC_DIR}/${END_ENTITY_CERT_NAME}.key"
}

function clean_certs_artifacts() {
    pushd "$MISC_DIR" > /dev/null || exit 1
    echo "rm -f -- *.csr *.srl *.ext"
    rm -f -- *.csr *.srl *.ext
    popd > /dev/null || exit 1
}

function clean_cert() {
    pushd "$MISC_DIR" > /dev/null || exit 1
    local cert="$1"
    echo "rm -f \"${cert}.der\" \"${cert}.key\""
    rm -f "${cert}.der" "${cert}.key"
    popd > /dev/null || exit 1
}

function clean_all() {
    clean_certs_artifacts
    clean_cert "$ROOT_CERT_NAME"
    clean_cert "$INTERMEDIATE_CERT_NAME"
    clean_cert "$END_ENTITY_CERT_NAME"
}

function regenerate() {
    FORCE=1
    rm -f "${MISC_DIR}/${ROOT_CERT_NAME}.der" \
          "${MISC_DIR}/${INTERMEDIATE_CERT_NAME}.der" \
          "${MISC_DIR}/${END_ENTITY_CERT_NAME}.der"
    create_root_cert
    create_intermediate_cert
    create_end_entity_cert
    if [[ $_should_clean_certs_artifacts == true ]]; then
        clean_certs_artifacts
    fi
}

function help() {
    set +e
    read -r -d '' usage <<-EOF
    Usage: gen-certs [-h] [-C] [COMMAND]

    This script is used to (re-)generate certificates used for a veraison
    deployment. The certificates are signed by a CA certificate called
    ${ROOT_CERT_NAME}.crt. If this does not exist, a self-signed one will
    be (re-)generated.

    Commands:

    create
            Create the root, intermediate, and end-entity certificates.

    regenerate
            Re-issue DER certificates with PKIX extensions, keeping existing keys.

    clean_certs_artifacts
            Clean output artifacts for the certificates.

    clean_all
            Clean both intermediate and output artifacts for everything (including
            the root CA cert).

    help
            Print this message and exit (same as -h option).

    Options:

    -h         Print this message and exit.
    -C         Do not clean up intermediate artifacts (e.g., CSRs).

    Note: Regenerating the certificate chain is an exceptional action and should
    only be done when necessary (e.g., when certificates expire).    

EOF

    echo "$usage"
    set -e
}

function _check_openssl() {
    if [[ "$(which openssl 2>/dev/null)" == "" ]]; then
        echo -e "ERROR: openssl executable must be installed to use this command."
        exit 1
    fi
}

function _check_root_cert() {
    if [[ ! -f "${MISC_DIR}/${ROOT_CERT_NAME}.der" ]]; then
        create_root_cert
    fi
}

_should_clean_certs_artifacts=true

OPTIND=1

while getopts "hC" opt; do
    case "$opt" in
        h) help; exit 0;;
        C) _should_clean_certs_artifacts=false;;
        *) break;;
    esac
done

shift $((OPTIND-1))
[ "${1:-}" = "--" ] && shift

command=$1
case $command in
    help)
        help
        exit 0
        ;;
    clean)
        clean_certs_artifacts
        ;;
    clean_all)
        clean_all
        ;;
    regenerate)
        regenerate
        ;;
    create)
        create_root_cert
        create_intermediate_cert
        create_end_entity_cert
        if [[ $_should_clean_certs_artifacts == true ]]; then
            clean_certs_artifacts
        fi
        ;;
    *)
        echo -e "ERROR: unexpected command: \"$command\" (use -h for help)"
        ;;
esac
