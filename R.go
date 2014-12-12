package main

import (
	"container/list"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func path(v interface{}, params *list.List) interface{} {
	if params.Len() == 0 {
		return v
	}
	pathStr := params.Remove(params.Front()).(string)
	path := strings.Split(pathStr, ".")
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

func keys(v interface{}, params *list.List) interface{} {
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

func values(v interface{}, params *list.List) interface{} {
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

func pick(v interface{}, params *list.List) interface{} {
	doc, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}
	out := map[string]interface{}{}
	for k, v := range doc {
		for e := params.Front(); e != nil; e = e.Next() {
			if k == e.Value {
				out[k] = v
			}
		}
	}
	params.Init()
	return out
}

func omit(v interface{}, params *list.List) interface{} {
	doc, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}
	for k, _ := range doc {
		omit := false
		for e := params.Front(); e != nil; e = e.Next() {
			if k == e.Value {
				omit = true
			}
		}
		if omit {
			delete(doc, k)
		}
	}
	params.Init()
	return doc
}

func where(v interface{}, params *list.List) interface{} {
	if params.Len() == 0 {
		return nil
	}
	v2 := unmarshal(params.Remove(params.Front()).(string))
	keys := keys(v2, nil).([]string)
	vpicked := pick(v, asList(keys))
	if deepequal(vpicked, v2) {
		return v
	} else {
		return nil
	}
}

func mixin(v interface{}, params *list.List) interface{} {
	doc, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}
	if params.Len() == 0 {
		return v
	}
	v2 := unmarshal(params.Remove(params.Front()).(string))
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

func eq(v interface{}, params *list.List) interface{} {
	if params.Len() == 0 {
		return false
	}
	v2 := unmarshal(params.Remove(params.Front()).(string))
	return deepequal(v, v2)
}

func not(v interface{}, params *list.List) interface{} {
	ret := dispatch(v, params)
	retbool, ok := ret.(bool)
	if !ok {
		return nil
	}
	return !retbool
}

func head(v interface{}, params *list.List) interface{} {
	list, ok := v.([]interface{})
	if !ok {
		return nil
	}
	if len(list) == 0 {
		return nil
	}
	return list[0]
}

func tail(v interface{}, params *list.List) interface{} {
	list, ok := v.([]interface{})
	if !ok {
		return nil
	}
	if len(list) <= 1 {
		return nil
	}
	return list[1:]
}

func Append(v interface{}, params *list.List) interface{} {
	if params.Len() == 0 {
		return v
	}
	list, ok := v.([]interface{})
	if !ok {
		return nil
	}
	v2 := unmarshal(params.Remove(params.Front()).(string))
	return append(list, v2)
}

func Concat(in interface{}, params *list.List) interface{} {
	if params.Len() == 0 {
		return in
	}
	list, ok := in.([]interface{})
	if !ok {
		return nil
	}
	arg := unmarshal(params.Remove(params.Front()).(string))
	listToAdd, ok := arg.([]interface{})
	if !ok {
		return nil
	}
	return append(list, listToAdd...)
}

func each(v interface{}, params *list.List) interface{} {
	list, ok := v.([]interface{})
	if !ok {
		return nil
	}
	for _, item := range list {
		marshal(item)
	}
	return nil
}

func obj_map(v interface{}, params *list.List) interface{} {
	list, ok := v.([]interface{})
	if !ok {
		return nil
	}
	ret := []interface{}{}
	for _, item := range list {
		out := dispatch(item, listCopy(params))
		ret = append(ret, out)
	}
	params.Init()
	return ret
}

func filter(v interface{}, params *list.List) interface{} {
	list, ok := v.([]interface{})
	if !ok {
		return nil
	}
	matches := []interface{}{}
	for _, item := range list {
		match := dispatch(item, listCopy(params))
		if match != nil {
			matches = append(matches, match)
		}
	}
	params.Init()
	return matches
}

func find(v interface{}, params *list.List) interface{} {
	list, ok := v.([]interface{})
	if !ok {
		return nil
	}
	for _, item := range list {
		match := dispatch(item, listCopy(params))
		if match != nil {
			params.Init()
			return match
		}
	}
	params.Init()
	return nil
}

func help(doc interface{}, params *list.List) interface{} {
	cmds := allCmds()
	fmt.Println("Usage: R <func> [arguments...]")
	if params.Len() > 0 {
		lookedFunc := params.Remove(params.Front()).(string)
		for _, cmd := range cmds {
			if cmd.name == lookedFunc {
				fmt.Println(cmd.descr)
				if cmd.example != "" {
					fmt.Printf("\nExample:\n%v\n", cmd.example)
				}
			}
		}
	} else {
		fmt.Println("\nAvailable functions:")
		for _, cmd := range cmds {
			fmt.Println(cmd.descr)
		}
	}
	return nil
}

type cmd struct {
	name    string
	run     func(interface{}, *list.List) interface{}
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
		{"append", Append, true,
			"append <obj>  appends object to list",
			` $ echo '[1]' | R append 2
 [1,2]`},
		{"concat", Concat, true,
			"concat [list]  concatenates two lists",
			` $ echo '[1,2]' | R concat '[3,4]'
 [1,2,3,4]`},
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

func dispatch(v interface{}, params *list.List) interface{} {
	for params.Len() > 0 {
		v = run(v, params)
	}
	return v
}

func run(v interface{}, params *list.List) interface{} {
	var cmdName string
	if params.Len() > 0 {
		cmdName = params.Remove(params.Front()).(string)
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
	return cmd.run(v, params)
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

func marshal(v interface{}) {
	_, isString := v.(string)
	if isString {
		fmt.Println(v)
	} else {
		enc := json.NewEncoder(os.Stdout)
		enc.Encode(&v)
	}
}

func exitCode(v interface{}) int {
	if v == nil {
		return 1
	}
	vBool, ok := v.(bool)
	if ok && !vBool {
		return 1
	}
	return 0
}

func asList(arr []string) *list.List {
	l := list.New()
	for _, i := range arr {
		l.PushBack(i)
	}
	return l
}

func listCopy(src *list.List) *list.List {
	l := list.New()
	for e := src.Front(); e != nil; e = e.Next() {
		l.PushBack(e.Value)
	}
	return l
}

func main() {
	args := asList(os.Args[1:])
	if ret := dispatch(nil, args); ret != nil {
		marshal(ret)
		os.Exit(exitCode(ret))
	}
}
