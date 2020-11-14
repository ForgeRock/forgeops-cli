#!/usr/bin/env bash
set -x 
case "${OSTYPE}" in
    "darwin"*) os="darwin";;
    "linux"*) os="linux";;
esac

latest=$(curl -s -w "%{redirect_url}" https://github.com/golangci/golangci-lint/releases/latest -o /dev/null | awk -F "/v" '{ printf $NF }')

if ! curl -o /tmp/golinter.tar.gz -L "https://github.com/golangci/golangci-lint/releases/download/v${latest}/golangci-lint-${latest}-${os}-amd64.tar.gz";
then
    echo "failed to download go linter"
    exit 1;
fi

if ! mkdir -p /tmp/golinter; 
then
    echo "failed to setup temp directory"
    exit 1;
fi

if ! tar xf /tmp/golinter.tar.gz -C /tmp/golinter --strip-components 1
then
    echo "failed to extract go linter"
    exit 1;
fi

if ! install /tmp/golinter/golangci-lint $(go env GOPATH)/bin;
then
    echo "failed to install to bin/"
    exit 1;
fi

if ! rm -rf /tmp/{golinter,golinter.tar.gz};
then
    echo "failed to cleanup"
    exit 1;
fi
