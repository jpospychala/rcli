#!/bin/bash

export PATH=$PATH:./

echo '{"a":{"b":true}}' | R path a.b | R eq true
echo '{"a":{"c":2}}' | R path a.c | R eq 2

echo '[1,2,3,4]' | R head | R eq 1
echo '[1,2,3,4]' | R tail | R eq '[2,3,4]'
