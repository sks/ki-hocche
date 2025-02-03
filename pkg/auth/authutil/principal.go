package authutil

import "github.com/google/uuid"

func DemoUser() Principal {
	return Principal{
		Zone:   uuid.MustParse("00000000-0000-0000-0000-000000000001"),
		UserID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
		Name:   "Demo User",
	}
}

type Principal struct {
	Zone   uuid.UUID
	UserID uuid.UUID
	Name   string
}
