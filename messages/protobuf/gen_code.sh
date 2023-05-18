#!/usr/bin/env bash
protoc --go_out=../../wtf-main/lib --python_betterproto_out=../../clients/py-client ./messages.proto
