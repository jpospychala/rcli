#!/bin/bash

export PATH=$PATH:./

# eq
echo '"abc"' | R eq '"abc"'
echo '1' | R eq '1'
echo '1.0' | R eq '1.0'
echo 'true' | R eq 'true'
echo '[1,2,3]' | R eq '[1,2,3]'
echo '{"a":{"b":1}}' | R eq '{"a":{"b":1}}'

# not eq
echo '"abc"' | R not eq '"jkl"'
echo '0' | R not eq '1'
echo '1' | R not eq '{}'
echo '{"a":1}' | R not eq '{"b":1,"c":2,"a":1}'

# pick
echo '{}' | R pick a | R eq '{}'
echo '{"a":1,"b":2}' | R pick a | R eq '{"a":1}'
# path
echo '{"a":{"b":true}}' | R path a.b | R eq true
echo '{"a":{"c":2}}' | R path a.c | R eq 2

# head
echo '[1,2,3,4]' | R head | R eq 1

# tail
echo '[1,2,3,4]' | R tail | R eq '[2,3,4]'

# keys
echo '{"a":1,"b":2}' | R keys | R eq '["a", "b"]'
