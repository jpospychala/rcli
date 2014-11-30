#!/bin/bash

export PATH=$PATH:./

# path
echo '{"a":{"b":true}}' | R path a.b | R eq true
echo '{"a":{"c":2}}' | R path a.c | R eq 2

# head
echo '[1,2,3,4]' | R head | R eq 1

# tail
echo '[1,2,3,4]' | R tail | R eq '[2,3,4]'

# keys
echo '{"a":1,"b":2}' | R keys | R eq '["a", "b"]'
