RCLI: Ramda for Command Line
======================
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
