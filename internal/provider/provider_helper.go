package provider

import (
	"fmt"
	"net/url"
	"strings"
)

func attributesToMap(rawMap map[string]interface{}) (map[string]string, error) {
	strMap := map[string]string{}
	for k, v := range rawMap {
		str, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("non-string value used in map")
		}
		strMap[k] = str
	}
	return strMap, nil
}

func mapToAttributes(rawMap map[string]string) (map[string]interface{}, error) {
	attr := map[string]interface{}{}
	for k, v := range rawMap {
		attr[k] = v
	}
	return attr, nil
}

func removeInternalKeys(m map[string]string, d map[string]interface{}) map[string]string {
	for k := range m {
		if isInternalKey(k) && !isKeyInMap(k, d) {
			delete(m, k)
		}
	}
	return m
}

func isKeyInMap(key string, d map[string]interface{}) bool {
	if d == nil {
		return false
	}
	for k := range d {
		if k == key {
			return true
		}
	}
	return false
}

func isInternalKey(annotationKey string) bool {
	u, err := url.Parse("//" + annotationKey)
	if err != nil {
		return false
	}

	// allow user specified application specific keys
	if u.Hostname() == "app.kubernetes.io" {
		return false
	}

	// internal *.kubernetes.io keys
	if strings.HasSuffix(u.Hostname(), "kubernetes.io") {
		return true
	}

	// Specific to DaemonSet annotations, generated & controlled by the server.
	if strings.Contains(annotationKey, "deprecated.daemonset.templates.generation") {
		return true
	}

	// internal *.loft.sh keys
	if strings.HasSuffix(u.Hostname(), "loft.sh") {
		return true
	}

	return false
}

func getAddedModifiedAndDeleted(old map[string]interface{}, new map[string]interface{}) (map[string]interface{}, map[string]interface{}, map[string]interface{}, error) {
	added := map[string]interface{}{}
	modified := map[string]interface{}{}
	deleted := map[string]interface{}{}

	for k, v := range new {
		_, oldHasKey := old[k]
		if !oldHasKey {
			added[k] = v
		}

		if oldHasKey {
			modified[k] = v
		}
	}

	for k, v := range old {
		_, newHasKey := new[k]
		if !newHasKey {
			deleted[k] = v
		}
	}

	return added, modified, deleted, nil
}
