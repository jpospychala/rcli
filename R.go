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

func values(v interface{}, params []string) interface{} {
	doc, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}
	vals := []interface{}{}
	for _, v := range doc {
		vals = append(vals, v)
	}
	return vals
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

func omit(v interface{}, params []string) interface{} {
	doc, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}
	for k, _ := range doc {
		omit := false
		for _, p := range params {
			if k == p {
				omit = true
			}
		}
		if omit {
			delete(doc, k)
		}
	}
	return doc
}

func where(v interface{}, params []string) interface{} {
	v2 := unmarshal(params[0])
	keys := keys(v2, nil).([]string)
	vpicked := pick(v, keys)
	if deepequal(vpicked, v2) {
		return v
	} else {
		return nil
	}
}

func mixin(v interface{}, params []string) interface{} {
	doc, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}
	v2 := unmarshal(params[0])
	v2map, ok := v2.(map[string]interface{})
	if !ok {
		return nil
	}
	for k, v := range v2map {
		doc[k] = v
	}
	return doc

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
	case json.Number:
		v2num, ok := v2.(json.Number)
		if !ok {
			return false
		}
		return v1.String() == v2num.String()
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

func append_to_list(v interface{}, params []string) interface{} {
	list, ok := v.([]interface{})
	if !ok {
		return nil
	}
	v2 := unmarshal(params[0])
	return append(list, v2)
}

func each(v interface{}, params []string) interface{} {
	list, ok := v.([]interface{})
	if !ok {
		return nil
	}
	for _, item := range list {
		marshal(item, false)
	}
	return nil
}

func obj_map(v interface{}, params []string) interface{} {
	list, ok := v.([]interface{})
	if !ok {
		return nil
	}
	ret := []interface{}{}
	for _, item := range list {
		out := dispatch(item, params)
		ret = append(ret, out)
	}
	return ret
}

func filter(v interface{}, params []string) interface{} {
	list, ok := v.([]interface{})
	if !ok {
		return nil
	}
	matches := []interface{}{}
	for _, item := range list {
		match := dispatch(item, params)
		if match != nil {
			matches = append(matches, match)
		}
	}
	return matches
}

func find(v interface{}, params []string) interface{} {
	list, ok := v.([]interface{})
	if !ok {
		return nil
	}
	for _, item := range list {
		match := dispatch(item, params)
		if match != nil {
			return match
		}
	}
	return nil
}

func help(doc interface{}, params []string) interface{} {
	cmds := allCmds()
	fmt.Println("Usage: R <command> [arguments...]")
	if len(params) > 0 {
		for _, cmd := range cmds {
			if cmd.name == params[0] {
				fmt.Println(cmd.descr)
				if cmd.example != "" {
					fmt.Println("\nExample:\n", cmd.example)
				}
			}
		}
	} else {
		for _, cmd := range cmds {
			fmt.Println(cmd.descr)
		}
	}
	return nil
}

type cmd struct {
	name    string
	run     func(interface{}, []string) interface{}
	stdin   bool
	descr   string
	example string
}

func allCmds() []cmd {
	return []cmd{
		// misc
		{"help", help, false,
			"help          prints usage details",
			""},
		// object
		{"path", path, true,
			"path <path>   returns object at period-delimited path",
			` $ echo '{"a":{"b":true}}' | R path a.b
 true`},
		{"keys", keys, true,
			"keys          returns object property names",
			` echo '{"a":1,"b":2}' | R keys
 ["a","b"]`},
		{"pick", pick, true,
			"pick [list]   returns object with only specified properties",
			` $ echo '{"a":1,"b":2}' | R pick a
 {"a":1}`},
		{"omit", omit, true,
			"omit [list]   returns object without specified properties",
			` $ echo '{"a":1,"b":2}' | R omit a
 {"b":2}`},
		{"values", values, true,
			"values        returns list of object values",
			` $ echo '{"a":1,"b":2}' | R values
 [1,2]`},
		{"where", where, true,
			"where <obj>   checks if object matches spec obj",
			` $ echo '{"a":1, "b":2}' | R where '{"a": 1}'
 {"a":1, "b":2}
 $ echo '{"a":1, "b":2}' | R where '{"a": b}'
 $`},
		{"mixin", mixin, true,
			"mixin <obj>   adds obj properties into input object",
			` $ echo '{"a":1, "b":2}' | R mixin '{"b": 5,"c":3}'
 {"a":1,"b":5,"c":3}`},
		// logic
		{"eq", eq, true,
			"eq  <obj>     compares stdin with first argument for equality",
			` $ echo '{"a":{"b":1}}' | R eq '{"a":{"b":1}}'`},
		{"not", not, false,
			"not <func>    inverts boolean result of following function",
			` $ echo '0' | R not eq '1'`},
		// array
		{"head", head, true,
			"head          first element of a list",
			` $ echo '[1,2,3,4]' | R head
 1`},
		{"tail", tail, true,
			"tail          all but first elements of a list",
			` $ echo '[1,2,3,4]' | R tail
 [2,3,4]`},
		{"each", each, true,
			"each          prints each list element in new line",
			` $ echo '[1,2,3]' | R each
 1
 2
 3`},
		{"map", obj_map, true,
			"map <func>    maps list elements using func",
			` $ echo '[{"a":1},{"a":2}]' | R map path a
 [1,2]`},
		{"append", append_to_list, true,
			"append <obj>  appends object to list",
			` $ echo '[1]' | R append 2
 [1,2]`},
		{"filter", filter, true,
			"filter <func> returns list of objects matching predicate",
			` $ echo '[{"a":1, "b":2}]' | R filter where '{"a":1}'
 [{"a":1,"b":2}]`},
		{"find", find, true,
			"find <func>   first object from list matching predicate",
			` $ echo '[{"a":1, "b":2}]' | R find where '{"a":1}'
 {"a":1,"b":2}
 $ echo '[{"a":1, "b":2}]' | R find where '{"a":2}'
 $`},
	}
}

func dispatch(v interface{}, params []string) interface{} {
	var cmdName string
	var args []string
	if len(params) > 0 {
		cmdName = params[0]
	}
	if len(params) > 1 {
		args = params[1:]
	}
	cmds := allCmds()
	cmd := cmds[0]
	for _, i := range cmds {
		if i.name == cmdName {
			cmd = i
		}
	}
	if cmd.stdin == true {
		dec := json.NewDecoder(os.Stdin)
		dec.UseNumber()
		dec.Decode(&v)
	}
	return cmd.run(v, args)
}

func unmarshal(param string) interface{} {
	var v2 interface{}
	r := strings.NewReader(param)
	dec := json.NewDecoder(r)
	dec.UseNumber()
	err := dec.Decode(&v2)
	if err != nil {
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
