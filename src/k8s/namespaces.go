package k8s

import (
	"context"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"regexp"
	"strings"
)

func GetNamespaces(includeValue, excludeValue string) ([]string, error) {
	var err error
	var allNamespaces []string
	var includeList, excludeList []*regexp.Regexp

	includeList, err = getNamespaceRegexList(includeValue, "default")
	if err != nil {
		return nil, err
	}

	excludeList, err = getNamespaceRegexList(excludeValue, "")
	if err != nil {
		return nil, err
	}

	allNamespaces, err = getAllNamespaces()
	if err != nil {
		return nil, err
	}

	result := make([]string, 0)
	for _, ns := range allNamespaces {
		if isAnyMatch(ns, includeList) && !isAnyMatch(ns, excludeList) {
			result = append(result, ns)
		}
	}

	return result, nil
}

func getNamespaceRegexList(value, defaultValue string) ([]*regexp.Regexp, error) {
	result := make([]*regexp.Regexp, 0, 0)

	if value == "" {
		value = defaultValue
	}

	value = formatNamespaceList(value)
	if value == "" {
		return result, nil
	}

	for _, val := range strings.Split(value, ",") {
		r, err := getRegex(val)
		if err != nil {
			return nil, err
		}
		result = append(result, r)
	}

	return result, nil
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

func getRegex(val string) (*regexp.Regexp, error) {
	regex := strings.Replace(val, "*", ".*", -1)
	regex = strings.Replace(regex, "?", ".", -1)
	regex = "^" + regex + "$"
	return regexp.Compile(regex)
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
