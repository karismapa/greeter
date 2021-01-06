#!/bin/bash

protoc --go_out=plugins=grpc:. greetpb/greet.proto
