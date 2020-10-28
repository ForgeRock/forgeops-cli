#!/usr/bin/env bash

if ! curl -o /tmp/goreleaser.tar.gz -L https://github.com/goreleaser/goreleaser/releases/latest/download/goreleaser_Linux_x86_64.tar.gz;
then
    echo "failed to download go releaser"
    exit 1;
fi
if ! mkdir /tmp/releaser && tar xf /tmp/goreleaser.tar.gz -C /tmp/releaser;
then
    echo "failed to extract go releaser"
    exit 1;
fi
if ! install /tmp/goreleaser bin/goreleaser; 
then
    echo "failed to install to bin/"
    exit 1;
fi
if ! rm -rf /tmp/{releaser,goreleaser.tar.gz};
then
    echo "failed to cleanup"
    exit 1;
fi

