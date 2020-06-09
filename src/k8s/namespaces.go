package k8s

import (
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	. "regexp"
	"strings"
)

func GetNamespaces(envVar string) ([]string, error) {
	list := strings.Split(envVar, ",")
	single := make([]string, 0, len(list))
	regex := make([]*Regexp, 0, len(list))

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

	matchedNamespaces, err := findNamespaces(regex)
	if err != nil {
		return nil, err
	}

	return append(single, matchedNamespaces...), nil
}

func hasWildCard(val string) bool {
	for _, r := range []rune{'*', '?'} {
		if strings.ContainsRune(val, r) {
			return true
		}
	}
	return false
}

func getRegex(val string) (*Regexp, error) {
	regex := strings.Replace(val, "*", ".*", -1)
	regex = strings.Replace(regex, "?", ".", -1)
	regex = "^" + regex + "$"
	return Compile(regex)
}

func findNamespaces(regex []*Regexp) ([]string, error) {
	if len(regex) == 0 {
		return nil, nil
	}

	namespaces, err := getAllNamespaces()
	if err != nil {
		return nil, err
	}

	result := make([]string, 0, len(namespaces))
	for _, ns := range namespaces {
		if isAnyMatch(ns, regex) {
			result = append(result, ns)
		}
	}

	return result, nil
}

func isAnyMatch(ns string, regexes []*Regexp) bool {
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
		res, err := client.CoreV1().Namespaces().List(opts)
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
