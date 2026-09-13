package domain

type Role string

const (
	Owner        Role = "owner"
	Collaborator Role = "collaborator"
	Guest        Role = "guest"
)

type Rights struct {
	Read   bool
	Write  bool
	Manage bool
}

func (r Role) Rights() Rights {
	switch r {
	case Owner:
		return Rights{Read: true, Write: true, Manage: true}
	case Collaborator:
		return Rights{Read: true, Write: true}
	case Guest:
		return Rights{Read: true}
	default:
		return Rights{}
	}
}
