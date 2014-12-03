package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func path(v interface{}, params []string) interface{} {
	path := strings.Split(params[0], ".")
	for _, seg := range path {
		doc, ok := v.(map[string]interface{})
		if ok {
			v = doc[seg]
		} else {
			doc, ok := v.([]interface{})
			if !ok {
				return nil
			}
			segInt, err := strconv.Atoi(seg)
			if err != nil {
				return nil
			}
			if (segInt < 0) || (segInt >= len(doc)) {
				return nil
			}
			v = doc[segInt]
		}
	}
	return v
}

func keys(v interface{}, params []string) interface{} {
	doc, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}
	keys := []string{}
	for k := range doc {
		keys = append(keys, k)
	}
	return keys
}

func pick(v interface{}, params []string) interface{} {
	doc, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}
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

func deepequal(v interface{}, v2 interface{}) bool {
	switch v1 := v.(type) {
	case int:
		v2int, ok := v2.(int)
		return ok && v1 == v2int
	case bool:
		v2bool, ok := v2.(bool)
		return ok && v1 == v2bool
	case float64:
		v2float, ok := v2.(float64)
		return ok && v1 == v2float
	case string:
		v2float, ok := v2.(string)
		return ok && v1 == v2float
	case []interface{}:
		v2list, ok := v2.([]interface{})
		if ok && len(v1) == len(v2list) {
			for i, val := range v1 {
				if !deepequal(val, v2list[i]) {
					return false
				}
			}
			return true
		}
	case map[string]interface{}:
		v2map, ok := v2.(map[string]interface{})
		if ok && len(v1) == len(v2map) {
			for k, val := range v1 {
				if !deepequal(val, v2map[k]) {
					return false
				}
			}
			return true
		}
	}
	return false
}

func eq(v interface{}, params []string) interface{} {
	v2 := unmarshal(params[0])
	return deepequal(v, v2)
}

func not(v interface{}, params []string) interface{} {
	ret := dispatch(v, params)
	retbool, ok := ret.(bool)
	if !ok {
		return nil
	}
	return !retbool
}

func head(v interface{}, params []string) interface{} {
	list, ok := v.([]interface{})
	if !ok {
		return nil
	}
	if len(list) == 0 {
		return nil
	}
	return list[0]
}

func tail(v interface{}, params []string) interface{} {
	list, ok := v.([]interface{})
	if !ok {
		return nil
	}
	if len(list) <= 1 {
		return nil
	}
	return list[1:]
}

func each(v interface{}, params []string) interface{} {
	list, ok := v.([]interface{})
	if ok {
		for _, item := range list {
			marshal(item, false)
		}
	}
	return nil
}

func help(doc interface{}, params []string) interface{} {
	fmt.Println("Usage: R <command> [arguments...]")
	fmt.Println("path  returns document from specific path")
	fmt.Println("keys  returns keys of JSON document")
	fmt.Println("pick  picks element from JSON")
	fmt.Println("eq    compares stdin with first argument for equality")
	fmt.Println("not   negation")
	fmt.Println("head  head of a list")
	fmt.Println("path  tail of a list")
	fmt.Println("each  prints each list element in new line")
	fmt.Println("help  prints usage details")
	return nil
}

type cmd struct {
	run   func(interface{}, []string) interface{}
	stdin bool
}

func dispatch(v interface{}, params []string) interface{} {
	cmdName := params[0]
	args := params[1:]
	cmds := map[string]cmd{
		// object
		"path": {path, true},
		"keys": {keys, true},
		"pick": {pick, true},
		// logic
		"eq":  {eq, true},
		"not": {not, false},
		// array
		"head": {head, true},
		"tail": {tail, true},
		"each": {each, true},
		// misc
		"help": {help, false},
	}
	if cmd, ok := cmds[cmdName]; ok {
		if cmd.stdin == true {
			dec := json.NewDecoder(os.Stdin)
			dec.Decode(&v)
		}
		return cmd.run(v, args)
	}
	return nil
}

func unmarshal(param string) interface{} {
	var v2 interface{}
	if err := json.Unmarshal([]byte(param), &v2); err != nil {
		return param
	}
	return v2
}

func marshal(v interface{}, canExit bool) {
	_, isString := v.(string)
	if isString {
		fmt.Println(v)
	} else {
		enc := json.NewEncoder(os.Stdout)
		enc.Encode(&v)
		vBool, ok := v.(bool)
		if canExit && ok && !vBool {
			os.Exit(1)
		}
	}
}

func main() {
	if ret := dispatch(nil, os.Args[1:]); ret != nil {
		marshal(ret, true)
	}
}
