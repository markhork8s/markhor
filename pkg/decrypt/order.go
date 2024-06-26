package decrypt

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/markhork8s/markhor/pkg"
	"github.com/markhork8s/markhor/pkg/config"
	orderedmap "github.com/wk8/go-ordered-map/v2"
)

func sortJson(jsonData map[string]interface{}, eid slog.Attr, config config.MarkhorSecretsConfig) *orderedmap.OrderedMap[string, interface{}] {
	ordered, err := sortJSONData(jsonData, eid, config)
	if err != nil {
		slog.Debug(fmt.Sprintf("Could not order the JSON: %v. Will keep the default order (alphabetic in k8s)", err), eid)

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

func sortJSONData(jsonData map[string]interface{}, eid slog.Attr, config config.MarkhorSecretsConfig) (*orderedmap.OrderedMap[string, interface{}], error) {

	sortingParams, ok := jsonData[pkg.MARKHORPARAMS_MANIFEST_KEY].(map[string]interface{})
	if !ok {
		return nil, errors.New("missing key markhorParams")
	}

	separator := config.HierarchySeparator.Default
	customSeparator, ok := sortingParams["hierarchySeparator"].(string)
	if ok {
		if config.HierarchySeparator.AllowOverride {
			separator = customSeparator
			if config.HierarchySeparator.WarnOnOverride {
				slog.Warn(fmt.Sprintf("Using custom hierarchy separator: '%s'", customSeparator), eid)
			} else {
				slog.Debug(fmt.Sprintf("Using custom hierarchy separator: '%s'", customSeparator), eid)
			}
		} else {
			slog.Debug(fmt.Sprintf("This MarkhorSecret asked to use a custom hierarchy separator, '%s', but specifying a custom one is disabled in markhor's configuration (hierarchySeparator>allowOverride)", customSeparator), eid)
		}
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

	o, err := parseOrdering(rawOrder, separator)
	if err != nil {
		return nil, errors.New(fmt.Sprint("could not parse the ordering:", err))
	}
	ordering := Ordering{Name: "", Terms: o}

	return sortWithOrdering(jsonData, ordering)
}

func sortWithOrdering(jsonData map[string]interface{}, ordering Ordering) (*orderedmap.OrderedMap[string, interface{}], error) {
	om := orderedmap.New[string, interface{}]()

	for _, k := range ordering.Terms {
		if len(k.Terms) == 0 {
			data := jsonData[k.Name]
			if data == nil {
				return nil, errors.New(fmt.Sprint("no key ", k.Name, " in JSON"))
			}
			om.Set(k.Name, data)
		} else {
			data, ok := jsonData[k.Name].(map[string]interface{})
			if !ok {
				return nil, errors.New(fmt.Sprint("no key ", k.Name, " in JSON"))
			}
			nestedObj, err := sortWithOrdering(data, k)
			if err != nil {
				return nil, err
			}
			om.Set(k.Name, nestedObj)
		}
	}

	return om, nil
}

type Ordering struct {
	Name  string
	Terms []Ordering
}

const duplicatesErrorMessage = "has duplicates in it: "

func parseOrdering(input []string, separator string) ([]Ordering, error) {
	result, err := parseOrderingRec(input, separator)
	if err != nil {
		return nil, err
	}
	name, hasDuplicate := hasDuplicates(result, separator)
	if hasDuplicate {
		return nil, errors.New(fmt.Sprint(duplicatesErrorMessage, name))
	}
	return result, nil
}

func parseOrderingRec(input []string, separator string) ([]Ordering, error) {
	result := make([]Ordering, 0)
	inputLen := len(input)
	for i := 0; i < inputLen; i++ {
		s := input[i]
		parts := strings.SplitN(s, separator, 2)
		if len(parts) > 1 { //It is a nested element
			name := parts[0]
			terms := make([]string, 1)
			terms[0] = parts[1]
			for {
				i++
				if i >= inputLen {
					break
				}
				s2 := input[i]
				parts2 := strings.SplitN(s2, separator, 2)
				if parts2[0] != name {
					i--
					break
				} else if len(parts2) != 2 {
					return nil, fmt.Errorf("found one element %s missing %s%s", s2, separator, parts[1])
				}
				terms = append(terms, parts2[1])
			}
			nestedOrdering, err := parseOrderingRec(terms, separator)
			if err != nil {
				return nil, errors.New(fmt.Sprint(name, separator, err))
			}
			newOrdering := Ordering{
				Name:  name,
				Terms: nestedOrdering,
			}
			result = append(result, newOrdering)
		} else { //We reached a leaf
			newOrdering := Ordering{
				Name:  s,
				Terms: make([]Ordering, 0),
			}
			result = append(result, newOrdering)
		}
	}
	return result, nil
}

// Checks if an ordering has duplicate keys in it
func hasDuplicates(slice []Ordering, separator string) (string, bool) {
	seen := make(map[string]struct{})

	for _, value := range slice {
		name := value.Name
		if _, ok := seen[name]; ok {
			return name, true
		}
		seen[name] = struct{}{}
		n, h := hasDuplicates(value.Terms, separator)
		if h {
			return name + separator + n, true
		}
	}

	return "", false
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
