#!/usr/bin/env bash
set -euo pipefail

# Why does this file exist?
#
# Kubernetes code-generator does not support writing to Go module directories.
# So we create a temporary directory, symlink this module to it, and then delete
# it afterward.

#
# Commands
#

LN_S="${LN_S:-ln -s}"
MKDIR_P="${MKDIR_P:-mkdir -p}"
MKTEMP="${MKTEMP:-mktemp}"
RM_F="${RM_F:-rm -f}"

#
#
#

BASEDIR="$( $MKTEMP -d -t nebula-tasks-k8s.XXXXXXX )"
trap '$RM_F -r "${BASEDIR}"' EXIT

MODULE_NAME="$( go list -m )"
MODULE_DIR="$( go list -m -f '{{ .Dir }}' )"

$MKDIR_P "${BASEDIR}/src/$( dirname "${MODULE_NAME}" )"
$LN_S "${MODULE_DIR}" "${BASEDIR}/src/${MODULE_NAME}"

GOPKGPATH="$( go env GOPATH )/pkg"
[ -d "${GOPKGPATH}" ] && $LN_S "${GOPKGPATH}" "${BASEDIR}/pkg"

GOPATH="${BASEDIR}" bash "${MODULE_DIR}/vendor/k8s.io/code-generator/generate-groups.sh" \
  all \
  "${MODULE_NAME}/pkg/generated" \
  "${MODULE_NAME}/pkg/apis" \
  nebula.puppet.com:v1 \
  --go-header-file "${MODULE_DIR}/hack/generate-boilerplate.go.txt"