RCLI: Ramda for Command Line
============================

[![Build Status](https://travis-ci.org/jpospychala/rcli.svg?branch=master)](https://travis-ci.org/jpospychala/rcli)

JSON manipulation tools for command line inspired by Ramda.js, a practical functional JavaScript library inspired by Clojure.

```bash
$ echo '{"status":"RUNNING"}' | R path status
RUNNING

$ echo '{"age":60,"color":"blue", "score": 3}' | R pick age score
{"age": 60, "score": 3}

$ echo '[1, 2, 3, 4]' | R head
1

$ echo '[1, 2, 3, 4]' | R tail
[2, 3, 4]

$ echo '{"age":60}' | R eq '{"age":60}'
true

$ echo '{"age":60}' | R not eq '{"name":"joe"}'
true
```

Building
========
RCLI is written in Go and requires ```go``` command to build.
```bash
git clone https://github.com/jpospychala/rcli.git
cd rcli
make
sudo make install
```

Usage
=====
Usage: R <func> [arguments...]

Functions can be stacked, so that result of one function is passed as input to next function,
for example:

```bash
$ echo '[1]' | R append 2 each
1
2
```

Functions
=========

append
------
append <obj>  appends object to list

Example:
```bash
 $ echo '[1]' | R append 2
 [1,2]
```

concat
------
concat [list]  concatenates two lists

Example:
```bash
 $ echo '[1,2]' | R concat '[3,4]'
 [1,2,3,4]
```

each
----
each          prints each list element in new line

Example:
```bash
 $ echo '[1,2,3]' | R each
 1
 2
 3
```

eq
--
eq  <obj>     compares stdin with first argument for equality

Example:
```bash
 $ echo '{"a":{"b":1}}' | R eq '{"a":{"b":1}}'
```

filter
------
filter <func> returns list of objects matching predicate

Example:
```bash
 $ echo '[{"a":1, "b":2}]' | R filter where '{"a":1}'
 [{"a":1,"b":2}]
```

find
----
find <func>   first object from list matching predicate

Example:
```bash
 $ echo '[{"a":1, "b":2}]' | R find where '{"a":1}'
 {"a":1,"b":2}
 $ echo '[{"a":1, "b":2}]' | R find where '{"a":2}'
 $
```

head
----
head          first element of a list

Example:
```bash
 $ echo '[1,2,3,4]' | R head
 1
```

help
----
help          prints usage details


keys
----
keys          returns object property names

Example:
```bash
 echo '{"a":1,"b":2}' | R keys
 ["a","b"]
```

map
---
map <func>    maps list elements using func

Example:
```bash
 $ echo '[{"a":1},{"a":2}]' | R map path a
 [1,2]
```

mixin
-----
mixin <obj>   adds obj properties into input object

Example:
```bash
 $ echo '{"a":1, "b":2}' | R mixin '{"b": 5,"c":3}'
 {"a":1,"b":5,"c":3}
```

not
---
not <func>    inverts boolean result of following function

Example:
```bash
 $ echo '0' | R not eq '1'
```

omit
----
omit [list]   returns object without specified properties

Example:
```bash
 $ echo '{"a":1,"b":2}' | R omit a
 {"b":2}
```

path
----
path <path>   returns object at period-delimited path

Example:
```bash
 $ echo '{"a":{"b":true}}' | R path a.b
 true
```

pick
----
pick [list]   returns object with only specified properties

Example:
```bash
 $ echo '{"a":1,"b":2}' | R pick a
 {"a":1}
```

tail
----
tail          all but first elements of a list

Example:
```bash
 $ echo '[1,2,3,4]' | R tail
 [2,3,4]
```

values
------
values        returns list of object values

Example:
```bash
 $ echo '{"a":1,"b":2}' | R values
 [1,2]
```

version
-------
version       prints R version


where
-----
where <obj>   checks if object matches spec obj

Example:
```bash
 $ echo '{"a":1, "b":2}' | R where '{"a": 1}'
 {"a":1, "b":2}
 $ echo '{"a":1, "b":2}' | R where '{"a": b}'
 $
```
