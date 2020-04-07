package constraint

import (
	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/ecs"
)

type Rope struct {
	ecs.NilConstraint
	First        *ecs.Entity
	FirstAnchor  sprec.Vec3
	Second       *ecs.Entity
	SecondAnchor sprec.Vec3
	Length       float32
}

func (r Rope) ApplyCorrectionImpulses() {
	firstTransformComp := r.First.Transform
	secondTransformComp := r.Second.Transform

	firstAnchorRelativePosition := firstTransformComp.Orientation.MulVec3(r.FirstAnchor)
	secondAnchorRelativePosition := secondTransformComp.Orientation.MulVec3(r.SecondAnchor)

	firstAnchorPosition := sprec.Vec3Sum(firstTransformComp.Position, firstAnchorRelativePosition)
	secondAnchorPosition := sprec.Vec3Sum(secondTransformComp.Position, secondAnchorRelativePosition)

	deltaPosition := sprec.Vec3Diff(secondAnchorPosition, firstAnchorPosition)
	jacobian := sprec.UnitVec3(deltaPosition) // FIXME: Handle if deltaPosition == 0

	firstMotionComp := r.First.Motion
	secondMotionComp := r.Second.Motion

	firstPointVelocity := sprec.Vec3Sum(firstMotionComp.Velocity, sprec.Vec3Cross(firstMotionComp.AngularVelocity, firstAnchorRelativePosition))
	secondPointVelocity := sprec.Vec3Sum(secondMotionComp.Velocity, sprec.Vec3Cross(secondMotionComp.AngularVelocity, secondAnchorRelativePosition))
	deltaVelocity := sprec.Vec3Diff(secondPointVelocity, firstPointVelocity)

	firstEffectiveMass := 1.0 / ((1.0 / firstMotionComp.Mass) + sprec.Vec3Dot(sprec.Mat3Vec3Prod(sprec.InverseMat3(firstMotionComp.MomentOfInertia), sprec.Vec3Cross(firstAnchorRelativePosition, jacobian)), sprec.Vec3Cross(firstAnchorRelativePosition, jacobian)))
	secondEffectiveMass := 1.0 / ((1.0 / secondMotionComp.Mass) + sprec.Vec3Dot(sprec.Mat3Vec3Prod(sprec.InverseMat3(secondMotionComp.MomentOfInertia), sprec.Vec3Cross(secondAnchorRelativePosition, jacobian)), sprec.Vec3Cross(secondAnchorRelativePosition, jacobian)))
	totalMass := firstEffectiveMass * secondEffectiveMass / (firstEffectiveMass + secondEffectiveMass)

	impulseStrength := totalMass * sprec.Vec3Dot(jacobian, deltaVelocity)
	impulse := sprec.Vec3Prod(jacobian, impulseStrength)
	firstMotionComp.ApplyOffsetImpulse(firstAnchorRelativePosition, impulse)
	secondMotionComp.ApplyOffsetImpulse(secondAnchorRelativePosition, sprec.InverseVec3(impulse))
}

func (r Rope) ApplyCorrectionTranslations() {
	firstTransformComp := r.First.Transform
	secondTransformComp := r.Second.Transform

	firstAnchorRelativePosition := firstTransformComp.Orientation.MulVec3(r.FirstAnchor)
	secondAnchorRelativePosition := secondTransformComp.Orientation.MulVec3(r.SecondAnchor)

	firstAnchorPosition := sprec.Vec3Sum(firstTransformComp.Position, firstAnchorRelativePosition)
	secondAnchorPosition := sprec.Vec3Sum(secondTransformComp.Position, secondAnchorRelativePosition)

	deltaPosition := sprec.Vec3Diff(secondAnchorPosition, firstAnchorPosition)
	jacobian := sprec.UnitVec3(deltaPosition) // FIXME: Handle if deltaPosition == 0

	firstMotionComp := r.First.Motion
	secondMotionComp := r.Second.Motion

	firstEffectiveMass := 1.0 / ((1.0 / firstMotionComp.Mass) + sprec.Vec3Dot(sprec.Mat3Vec3Prod(sprec.InverseMat3(firstMotionComp.MomentOfInertia), sprec.Vec3Cross(firstAnchorRelativePosition, jacobian)), sprec.Vec3Cross(firstAnchorRelativePosition, jacobian)))
	secondEffectiveMass := 1.0 / ((1.0 / secondMotionComp.Mass) + sprec.Vec3Dot(sprec.Mat3Vec3Prod(sprec.InverseMat3(secondMotionComp.MomentOfInertia), sprec.Vec3Cross(secondAnchorRelativePosition, jacobian)), sprec.Vec3Cross(secondAnchorRelativePosition, jacobian)))
	totalMass := firstEffectiveMass * secondEffectiveMass / (firstEffectiveMass + secondEffectiveMass)

	nudgeStrength := totalMass * sprec.Vec3Dot(jacobian, sprec.ResizedVec3(deltaPosition, deltaPosition.Length()-r.Length))
	nudge := sprec.Vec3Prod(jacobian, nudgeStrength)

	firstTransformComp.Translate(sprec.Vec3Quot(nudge, firstMotionComp.Mass))
	secondTransformComp.Translate(sprec.InverseVec3(sprec.Vec3Quot(nudge, secondMotionComp.Mass)))

	firstTransformComp.Rotate(sprec.Mat3Vec3Prod(sprec.InverseMat3(firstMotionComp.MomentOfInertia), sprec.Vec3Cross(firstAnchorRelativePosition, nudge)))
	secondTransformComp.Rotate(sprec.InverseVec3(sprec.Mat3Vec3Prod(sprec.InverseMat3(secondMotionComp.MomentOfInertia), sprec.Vec3Cross(secondAnchorRelativePosition, nudge))))
}

// func (r Rope) ApplyCorrectionTranslations() {
// 	// firstTransformComp := r.First.Transform
// 	// secondTransformComp := r.Second.Transform
// 	// deltaPosition := sprec.Vec3Diff(secondTransformComp.Position, firstTransformComp.Position)
// 	// jacobian := sprec.UnitVec3(deltaPosition)
// 	// posError := (deltaPosition.Length() - r.Length)

// 	// firstMotionComp := r.First.Motion
// 	// secondMotionComp := r.Second.Motion

// 	// firstCorrectionAmount := posError * secondMotionComp.Mass / (firstMotionComp.Mass + secondMotionComp.Mass)
// 	// secondCorrectionAmount := -posError * firstMotionComp.Mass / (firstMotionComp.Mass + secondMotionComp.Mass)

// 	// firstTransformComp.Translate(sprec.Vec3Prod(jacobian, firstCorrectionAmount))
// 	// secondTransformComp.Translate(sprec.Vec3Prod(jacobian, secondCorrectionAmount))

// 	// // fmt.Printf("first comp position: %#v\n", firstTransformComp.Position)
// 	// // fmt.Printf("second comp position: %#v\n", secondTransformComp.Position)
// 	// // fmt.Println("translating...")
// }
