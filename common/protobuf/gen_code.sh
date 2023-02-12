#!/usr/bin/env bash
protoc --go_out=../../server --python_betterproto_out=../../clients/py-client ./messages.proto
