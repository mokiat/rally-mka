package physics

import "github.com/mokiat/gomath/sprec"

type MatchRotationConstraint struct {
	NilConstraint
	FirstBody  *Body
	SecondBody *Body
}

func (c MatchRotationConstraint) ApplyImpulse() {
	c.yConstraint().ApplyImpulse()
	c.zConstraint().ApplyImpulse()
}

func (c MatchRotationConstraint) ApplyNudge() {
	c.yConstraint().ApplyNudge()
	c.zConstraint().ApplyNudge()
}

func (c MatchRotationConstraint) yConstraint() MatchAxisConstraint {
	return MatchAxisConstraint{
		FirstBody:      c.FirstBody,
		FirstBodyAxis:  sprec.BasisYVec3(),
		SecondBody:     c.SecondBody,
		SecondBodyAxis: sprec.BasisYVec3(),
	}
}

func (c MatchRotationConstraint) zConstraint() MatchAxisConstraint {
	return MatchAxisConstraint{
		FirstBody:      c.FirstBody,
		FirstBodyAxis:  sprec.BasisZVec3(),
		SecondBody:     c.SecondBody,
		SecondBodyAxis: sprec.BasisZVec3(),
	}
}
