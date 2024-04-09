package main

import (
	"errors"
	"fmt"
	"strings"

	orderedmap "github.com/wk8/go-ordered-map/v2"
)

func sortJson(jsonData map[string]interface{}) *orderedmap.OrderedMap[string, interface{}] {
	ordered, err := sortJSONData(jsonData)
	if err != nil {
		fmt.Println("Error ordering the JSON:", err, "\nWill keep the default order (alphabetic in k8s)")

		orderedMap := orderedmap.New[string, interface{}]()

		for key, value := range jsonData {
			orderedMap.Set(key, value)
		}

		return orderedMap
	}
	_, present := ordered.Get("sops")
	if !present {
		ordered.Set("sops", jsonData["sops"])
	}
	return ordered
}

func sortJSONData(jsonData map[string]interface{}) (*orderedmap.OrderedMap[string, interface{}], error) {

	sortingParams, ok := jsonData["markhorParams"].(map[string]interface{})
	if !ok {
		return nil, errors.New("missing key markhorParams")
	}
	separator, ok := sortingParams["hierarchySeparator"].(string)
	if !ok {
		separator = "."
	} else {
		fmt.Println("Info: using custom hierarchy separator: ", separator)
	}
	rawOrderIntf, ok := sortingParams["order"].([]interface{})
	if !ok {
		return nil, errors.New("no order field found")
	}
	rawOrder := make([]string, len(rawOrderIntf))
	for i, v := range rawOrderIntf {
		s, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("order term number %d (0-indexed) is not a string", i)
		}
		rawOrder[i] = s
	}

	ordering := Ordering{name: "", terms: parseOrdering(rawOrder, separator)}

	return sortWithOrdering(jsonData, ordering), nil
}

func sortWithOrdering(jsonData map[string]interface{}, ordering Ordering) *orderedmap.OrderedMap[string, interface{}] {
	om := orderedmap.New[string, interface{}]()

	for _, k := range ordering.terms {
		if len(k.terms) == 0 {
			om.Set(k.name, jsonData[k.name])
		} else {
			data, ok := jsonData[k.name].(map[string]interface{})
			if !ok {
				fmt.Println("Error no key ", k.name, "in JSON")
			}
			nestedObj := sortWithOrdering(data, k)
			om.Set(k.name, nestedObj)
		}
	}

	return om
}

type Ordering struct {
	name  string
	terms []Ordering
}

func parseOrdering(raw []string, separator string) []Ordering {
	result := make([]Ordering, 0)
	lr := len(raw)
	for i := 0; i < lr; i++ {
		s := raw[i]
		parts := strings.SplitN(s, separator, 2)
		if len(parts) > 1 {
			name := parts[0]
			terms := make([]string, 1)
			terms[0] = parts[1]
			for {
				i++
				if i >= lr {
					break
				}
				s2 := raw[i]
				parts := strings.SplitN(s2, separator, 2)
				if parts[0] != name {
					i--
					break
				}
				terms = append(terms, parts[1])
			}
			newOrdering := Ordering{
				name:  name,
				terms: parseOrdering(terms, separator),
			}
			result = append(result, newOrdering)
		} else {
			newOrdering := Ordering{
				name:  s,
				terms: make([]Ordering, 0),
			}
			result = append(result, newOrdering)
		}
	}
	return result
}

// func sortAlphabetically(data map[string]interface{}) map[string]interface{} {
// 	result := make(map[string]interface{})
// 	keys := make([]string, 0, len(data))
// 	for key := range data {
// 		keys = append(keys, key)
// 	}
// 	sort.Strings(keys)

// 	for _, key := range keys {
// 		value := data[key]
// 		childData, ok := value.(map[string]interface{})
// 		if ok {
// 			result[key] = sortAlphabetically(childData)
// 		} else {
// 			result[key] = value
// 		}
// 	}

// 	return result
// }
