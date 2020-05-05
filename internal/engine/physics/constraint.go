package physics

type Constraint interface {
	Reset()
	ApplyImpulse()
	ApplyNudge()
}

var _ Constraint = NilConstraint{}

type NilConstraint struct{}

func (NilConstraint) Reset() {}

func (NilConstraint) ApplyImpulse() {}

func (NilConstraint) ApplyNudge() {}
