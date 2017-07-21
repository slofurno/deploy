#!/bin/bash

echo "deploying.."
echo $1

echo 'tag: $TRAVIS_TAG' | envsubst | cat

