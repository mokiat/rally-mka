package physics

type Context struct {
	ElapsedSeconds    float32
	ImpulseIterations int
	NudgeIterations   int
}

type Constraint interface {
	Reset()
	ApplyImpulse(ctx Context)
	ApplyNudge(ctx Context)
}

var _ Constraint = NilConstraint{}

type NilConstraint struct{}

func (NilConstraint) Reset() {}

func (NilConstraint) ApplyImpulse(Context) {}

func (NilConstraint) ApplyNudge(Context) {}
