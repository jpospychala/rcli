#!/bin/bash
#set -e

export PATH=./:$PATH

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

# omit
echo '{}' | R omit a | R eq '{}'
echo '{"a":1,"b":2}' | R omit a | R eq '{"b":2}'

# path
echo '{"a":{"b":true}}' | R path a.b | R eq true
echo '{"a":{"b":true}}' | R path a.b eq true
echo '{"a":{"c":2}}' | R path a.c | R eq 2
echo '[{"a":1}]' | R path 0.a | R eq 1
echo '{}' | R path -1
echo '[0]' | R path 1
echo '{}' | R path a.b.c.d

# head
echo '[1,2,3,4]' | R head | R eq 1
echo '[1]' | R head | R eq 1
echo '[]' | R head

# tail
echo '[1,2,3,4]' | R tail | R eq '[2,3,4]'
echo '[1]' | R tail
echo '[]' | R tail

# each
echo '[1,2,3]' | R each | head -1 | R eq 1
echo '[1,2,3]' | R each | tail -1 | R eq 3

# map
echo '[{"a":1},{"a":2}]' | R map path a | R eq '[1,2]'

# append
echo '[1]' | R append 2 | R eq '[1,2]'

# concat
echo '[1, 2]' | R concat '[3,4]' eq '[1,2,3,4]'

# values
echo '{"a":1,"b":2}' | R values | R eq '[1,2]'
echo '[1,2,3]' | R values

# keys
echo '{"a":1,"b":2}' | R keys | R eq '["a","b"]'

# where
echo '{"a":1, "b":2}' | R where '{"a": 1}' | R eq '{"a":1, "b":2}'

# filter
echo '[{"a":1, "b":2}]' | R filter where '{"a":1}' | R eq '[{"a":1,"b":2}]'

# find
echo '[{"a":1, "b":2}]' | R find where '{"a":1}' | R eq '{"a":1,"b":2}'
echo '[{"a":1, "b":2}]' | R find where '{"a":2}'

# mixin
echo '{"a":1, "b":2}' | R mixin '{"b": 5,"c":3}' | R eq '{"a":1,"b":5,"c":3}'
