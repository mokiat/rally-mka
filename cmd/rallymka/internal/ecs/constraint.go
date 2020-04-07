package ecs

type Constraint interface {
	ApplyForces()
	ApplyCorrectionForces()
	ApplyCorrectionImpulses()
	ApplyCorrectionTranslations()
}

var _ Constraint = NilConstraint{}

type NilConstraint struct{}

func (NilConstraint) ApplyForces() {}

func (NilConstraint) ApplyCorrectionForces() {}

func (NilConstraint) ApplyCorrectionImpulses() {}

func (NilConstraint) ApplyCorrectionTranslations() {}
