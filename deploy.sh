#!/bin/bash

echo "deploying.."

echo 'tag: $TRAVIS_TAG' | envsubst | cat

