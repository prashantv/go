#!/bin/bash
set -euo pipefail
set -x

doc2go $(find . -iname go.mod -exec dirname {} \;)

# add meta tag for go import
find _site -iname *.html -exec go run scripts/appendmeta/appendmeta.go {} \;

