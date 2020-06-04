package main

import (
	. "k8s.io/apimachinery/pkg/apis/meta/v1"
	"regexp"
	"strings"
)

func getNamespaces(envVar string) ([]string, error) {
	list := strings.Split(envVar, ",")
	wildcardMatches, err := getWildCardMatches(list)
	if nil != err {
		return nil, err
	}

	list = append(list, wildcardMatches...)
	return list, nil
}

func getWildCardMatches(list []string) ([]string, error) {
	var result []string

	wildcards, err := getWildCardRegexes(list)
	if nil != err || 0 == len(wildcards){
		return result, err
	}

	namespaces, err := getAllNamespaces()
	if nil != err {
		return result, err
	}

	for _, ns := range namespaces {
		if isAnyMatch(ns, wildcards) {
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

func getWildCardRegexes(list []string) ([]*regexp.Regexp, error) {
	var result []*regexp.Regexp
	for _, ns := range list {
		if !strings.Contains(ns, "*") && !strings.Contains(ns, "?") {
			continue
		}

		regex := strings.Replace(ns, "*", ".*", -1)
		regex = strings.Replace(ns, "?", ".", -1)
		regex = "^" + regex + "$"
		r, err := regexp.Compile(regex)
		if nil != err {
			return nil, err
		}

		result = append(result, r)
	}

	return result, nil
}

func getAllNamespaces() ([]string, error) {
	var result []string

	client, err := getClient()
	if nil != err {
		return nil, err
	}

	opts := ListOptions{}
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
