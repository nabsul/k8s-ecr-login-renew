package k8s

import (
	"context"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"regexp"
	"strings"
)

func GetNamespaces(includeValue, excludeValue string) ([]string, error) {
	allNamespaces, err := getAllNamespaces()
	if err != nil {
		return nil, err
	}

	includeList, err := getNamespaceList(includeValue, "default", allNamespaces)
	if err != nil {
		return nil, err
	}

	excludeList, err := getNamespaceList(excludeValue, "", allNamespaces)
	if err != nil {
		return nil, err
	}

	lookup := map[string]bool{}
	for _, xVal := range excludeList {
		lookup[xVal] = true
	}

	result := make([]string, 0)
	for _, iVal := range includeList {
		_, found := lookup[iVal]
		if found {
			result = append(result, iVal)
		}
	}

	return result, nil
}

func getNamespaceList(value, defaultValue string, allNamespaces []string) ([]string, error) {
	if value == "" {
		value = defaultValue
	}

	value = formatNamespaceList(value)
	if value == "" {
		return make([]string, 0, 0), nil
	}

	list := strings.Split(value, ",")
	single := make([]string, 0, len(list))
	regex := make([]*regexp.Regexp, 0, len(list))

	for _, val := range list {
		if hasWildCard(val) {
			r, err := getRegex(val)
			if err != nil {
				return nil, err
			}
			regex = append(regex, r)
		} else {
			single = append(single, val)
		}
	}

	matchedNamespaces, err := findNamespaces(regex, allNamespaces)
	if err != nil {
		return nil, err
	}

	return unique(append(single, matchedNamespaces...)), nil
}

var namespaceWhitespace = []string{" ", "\r", "\t", "\v"}
var namespaceSeparators = []string{"\n", ";"}

func formatNamespaceList(namespaceList string) string {
	formattedNamespaceList := namespaceList

	for _, c := range namespaceWhitespace {
		formattedNamespaceList = strings.ReplaceAll(formattedNamespaceList, c, "")
	}

	for _, c := range namespaceSeparators {
		formattedNamespaceList = strings.ReplaceAll(formattedNamespaceList, c, ",")
	}

	return formattedNamespaceList
}

func unique(values []string) []string {
	result := make([]string, 0, len(values))
	check := map[string]bool{}
	for _, val := range values {
		_, ok := check[val]
		if !ok {
			check[val] = true
			result = append(result, val)
		}
	}
	return result
}

func hasWildCard(val string) bool {
	for _, r := range []rune{'*', '?'} {
		if strings.ContainsRune(val, r) {
			return true
		}
	}
	return false
}

func getRegex(val string) (*regexp.Regexp, error) {
	regex := strings.Replace(val, "*", ".*", -1)
	regex = strings.Replace(regex, "?", ".", -1)
	regex = "^" + regex + "$"
	return regexp.Compile(regex)
}

func findNamespaces(regex []*regexp.Regexp, allNamespaces []string) ([]string, error) {
	if len(regex) == 0 {
		return nil, nil
	}

	result := make([]string, 0, len(allNamespaces))
	for _, ns := range allNamespaces {
		if isAnyMatch(ns, regex) {
			result = append(result, ns)
		}
	}

	return result, nil
}

func isAnyMatch(ns string, regexes []*regexp.Regexp) bool {
	for _, r := range regexes {
		if r.MatchString(ns) {
			return true
		}
	}

	return false
}

func getAllNamespaces() ([]string, error) {
	var result []string

	client, err := GetClient()
	if nil != err {
		return nil, err
	}

	opts := metaV1.ListOptions{}
	first := true

	for first || opts.Continue != "" {
		first = false
		res, err := client.CoreV1().Namespaces().List(context.TODO(), opts)
		if nil != err {
			return nil, err
		}

		opts.Continue = res.Continue
		newNames := make([]string, len(res.Items))
		for i, item := range res.Items {
			newNames[i] = item.Name
		}

		result = append(result, newNames...)
	}

	return result, nil
}
