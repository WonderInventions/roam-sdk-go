#!/usr/bin/env sh

set -e

SRCDIR=../developer-ro-am/static
SCRIPTDIR=$(dirname "$0")


PKG=roamv0
DSTDIR=$SCRIPTDIR/gen/$PKG
mkdir -p $DSTDIR
go run github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@v2.0.0 \
    -generate types -package $PKG $SRCDIR/openapi-alpha.json > $DSTDIR/$PKG.gen.go

PKG=roamv1
DSTDIR=$SCRIPTDIR/gen/$PKG
mkdir -p $DSTDIR
go run github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@v2.0.0 \
    -generate types -package $PKG $SRCDIR/openapi.json > $DSTDIR/$PKG.gen.go

PKG=roamwebhooks
DSTDIR=$SCRIPTDIR/gen/$PKG
mkdir -p $DSTDIR
go run github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@v2.0.0 \
    -generate types -package $PKG $SRCDIR/openapi-webhooks.json > $DSTDIR/$PKG.gen.go
