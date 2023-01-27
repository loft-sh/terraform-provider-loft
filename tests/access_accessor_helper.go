package tests

import (
	"fmt"
	v1 "github.com/loft-sh/api/v2/pkg/apis/storage/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func hasUser(user string) func(obj client.Object) error {
	return func(obj client.Object) error {
		accessAccessor, ok := obj.(v1.AccessAccessor)
		if !ok {
			return fmt.Errorf("object does not have an owner")
		}

		owner := accessAccessor.GetOwner()
		if owner.User != user {
			return fmt.Errorf(
				"%s: User didn't match %q, got %#v",
				obj.GetName(),
				user,
				owner.User)
		}

		return nil
	}
}

func hasTeam(team string) func(obj client.Object) error {
	return func(obj client.Object) error {
		accessAccessor, ok := obj.(v1.AccessAccessor)
		if !ok {
			return fmt.Errorf("object does not have an owner")
		}

		owner := accessAccessor.GetOwner()
		if owner.Team != team {
			return fmt.Errorf(
				"%s: Team didn't match %q, got %#v",
				obj.GetName(),
				team,
				owner.Team)
		}

		return nil
	}
}
