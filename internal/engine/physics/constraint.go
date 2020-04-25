package physics

type Constraint interface {
	ApplyForce()
	ApplyImpulse()
	ApplyBaumgarte()
	ApplyNudge()
}

var _ Constraint = NilConstraint{}

type NilConstraint struct{}

func (NilConstraint) ApplyForce() {}

func (NilConstraint) ApplyImpulse() {}

func (NilConstraint) ApplyBaumgarte() {}

func (NilConstraint) ApplyNudge() {}
