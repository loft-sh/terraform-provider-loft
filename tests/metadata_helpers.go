package tests

import (
	"fmt"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func hasAnnotation(annotation, value string) func(obj client.Object) error {
	return func(obj client.Object) error {
		if obj.GetAnnotations()[annotation] != value {
			return fmt.Errorf(
				"%s: Annotation '%s' didn't match %q, got %#v",
				obj.GetName(),
				annotation,
				value,
				obj.GetLabels()[annotation])
		}
		return nil
	}
}

func noAnnotation(annotation string) func(obj client.Object) error {
	return func(obj client.Object) error {
		if obj.GetAnnotations()[annotation] != "" {
			return fmt.Errorf(
				"%s: Annotation '%s' should not be present",
				obj.GetName(),
				annotation,
			)
		}
		return nil
	}
}

func hasLabel(label, value string) func(obj client.Object) error {
	return func(obj client.Object) error {
		if obj.GetLabels()[label] != value {
			return fmt.Errorf(
				"%s: Label '%s' didn't match %q, got %#v",
				obj.GetName(),
				label,
				value,
				obj.GetLabels()[label])
		}
		return nil
	}
}

func noLabel(label string) func(obj client.Object) error {
	return func(obj client.Object) error {
		if obj.GetAnnotations()[label] != "" {
			return fmt.Errorf(
				"%s: Label '%s' should not be present",
				obj.GetName(),
				label,
			)
		}
		return nil
	}
}
