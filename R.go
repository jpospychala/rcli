package main

import (
	"encoding/json"
	"log"
	"os"
	"strings"
)

func path(v interface{}, params []string) interface{} {
	path := strings.Split(params[0], ".")
	for _, seg := range path {
		doc := v.(map[string]interface{})
		v = doc[seg]
	}
	return v
}

func keys(v interface{}, params []string) interface{} {
	doc := v.(map[string]interface{})
	keys := []string{}
	for k := range doc {
		keys = append(keys, k)
	}
	return keys
}

func pick(v interface{}, params []string) interface{} {
	doc := v.(map[string]interface{})
	out := map[string]interface{}{}
	for k, v := range doc {
		for _, p := range params {
			if k == p {
				out[k] = v
			}
		}
	}
	return out
}

func eq(v interface{}, params []string) interface{} {
	var v2 interface{}
	if err := json.Unmarshal([]byte(params[0]), &v2); err != nil {
		log.Println(err)
		return nil
	}
	return v == v2
}

func head(v interface{}, params []string) interface{} {
	list := v.([]interface{})
	return list[0]
}

func tail(v interface{}, params []string) interface{} {
	list := v.([]interface{})
	return list[1:]
}

func help(doc interface{}, params []string) interface{} {
	return nil
}

type cmd struct {
	handler func(interface{}, []string) interface{}
	doc     string
}

func main() {
	cmdName := os.Args[1]
	args := os.Args[2:]
	dec := json.NewDecoder(os.Stdin)
	enc := json.NewEncoder(os.Stdout)

	cmds := map[string]cmd{
		// object
		"path": {path, "extracts part of JSON document"},
		"keys": {keys, "returns keys of JSON document"},
		"pick": {pick, "picks element from JSON"},
		// logic
		"eq": {eq, "compares stdin with first argument for equality"},
		// array
		"head": {head, "head of a list"},
		"tail": {tail, "tail of a list"},
		// misc
		"help": {help, "prints usage details"},
	}
	var v interface{}
	if err := dec.Decode(&v); err != nil {
		log.Println(err)
		return
	}
	if c, ok := cmds[cmdName]; ok {
		ret := c.handler(v, args)
		if err := enc.Encode(&ret); err != nil {
			log.Println(err)
		}
	}
}
