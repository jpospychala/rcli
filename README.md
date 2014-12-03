RCLI: Ramda for Command Line
============================
JSON manipulation tools for command line, inspired by Ramda.js, a practical functional library inspired by Clojure.

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

Usage
=====

path
----
Outputs JSON document node specified by path. Path is built of dot-separated segments.

```bash
$ echo '{"servers": {"sunshine": {"ip": "127.0.0.1" }}}' | R path servers.sunshine.ip
127.0.0.1

$ echo '[{"ip": "8.8.8.8"}]' | R path 0.ip
8.8.8.8
```

append
------
Outputs an array piped to stdin and adds new element to it. New array element should be
passed as an argument.

```bash
$ echo '[1, 2, 3]' | R append 4
[1,2,3,4]
```
