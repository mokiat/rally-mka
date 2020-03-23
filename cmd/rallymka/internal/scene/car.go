package scene

import (
	"fmt"

	"github.com/mokiat/go-whiskey/math"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/stream"
	"github.com/mokiat/rally-mka/internal/engine/collision"
)

const gravity = -0.05 // original: -0.003

func NewCar(stage *Stage, model *stream.Model, modelMatrix math.Mat4x4) *Car {
	bodyNode, _ := model.FindNode("body")
	flWheelNode, _ := model.FindNode("wheel_front_left")
	frWheelNode, _ := model.FindNode("wheel_front_right")
	blWheelNode, _ := model.FindNode("wheel_back_left")
	brWheelNode, _ := model.FindNode("wheel_back_right")

	return &Car{
		stage: stage,
		node:  bodyNode,

		Model:       model,
		ModelMatrix: modelMatrix,
		position:    math.MakeVec3(modelMatrix.M14, modelMatrix.M24, modelMatrix.M34),
		orientation: Orientation{
			VectorX: math.BaseVec3X(),
			VectorY: math.BaseVec3Y(),
			VectorZ: math.BaseVec3Z(),
		},
		flWheel: NewWheel(flWheelNode, true),
		frWheel: NewWheel(frWheelNode, true),
		blWheel: NewWheel(blWheelNode, false),
		brWheel: NewWheel(brWheelNode, false),
	}
}

type Car struct {
	stage *Stage
	node  *stream.Node

	ModelMatrix math.Mat4x4
	Model       *stream.Model

	position            math.Vec3
	orientation         Orientation
	velocity            math.Vec3
	acceleration        math.Vec3
	angularVelocity     math.Vec3
	angularAcceleration math.Vec3

	flWheel *Wheel
	frWheel *Wheel
	blWheel *Wheel
	brWheel *Wheel

	WheelAngle      float32
	Acceleration    float32
	HandbrakePulled bool
}

func (c Car) Position() math.Vec3 {
	return math.MakeVec3(c.ModelMatrix.M14, c.ModelMatrix.M24, c.ModelMatrix.M34)
}

func (c *Car) Update() {
	c.acceleration = math.NullVec3()
	c.angularAcceleration = math.NullVec3()

	c.acceleration = c.acceleration.IncCoords(0.0, gravity, 0.0)

	c.flWheel.RotationAngle += 20.0
	c.frWheel.RotationAngle += 20.0
	c.blWheel.RotationAngle += 20.0
	c.brWheel.RotationAngle += 20.0

	c.flWheel.TurnAngle = c.WheelAngle
	c.frWheel.TurnAngle = c.WheelAngle

	c.flWheel.updateBasics(c)
	c.frWheel.updateBasics(c)
	c.blWheel.updateBasics(c)
	c.brWheel.updateBasics(c)

	c.flWheel.Update(c)
	c.frWheel.Update(c)
	c.blWheel.Update(c)
	c.brWheel.Update(c)

	c.velocity = c.velocity.IncVec3(c.acceleration)
	c.velocity = c.velocity.Mul(0.999) // friction // FIXME
	c.angularVelocity = c.angularVelocity.IncVec3(c.angularAcceleration)
	c.angularVelocity = c.angularVelocity.Mul(0.99) // friction // FIXME

	c.move(c.velocity, c.angularVelocity)

	c.ModelMatrix = math.VectorMat4x4(
		c.orientation.VectorX,
		c.orientation.VectorY,
		c.orientation.VectorZ,
		c.position,
	)
}

func (c *Car) move(translation math.Vec3, rotation math.Vec3) {
	c.position = c.position.IncVec3(c.velocity)
	c.orientation.Rotate(c.angularVelocity)

	c.flWheel.updateBasics(c)
	c.frWheel.updateBasics(c)
	c.blWheel.updateBasics(c)
	c.brWheel.updateBasics(c)

	flDisplace, flTouching := c.checkWheel(c.flWheel)
	frDisplace, frTouching := c.checkWheel(c.frWheel)
	blDisplace, blTouching := c.checkWheel(c.blWheel)
	brDisplace, brTouching := c.checkWheel(c.brWheel)

	if flTouching || frTouching || blTouching || brTouching {
		var maxDisplace math.Vec3
		if flDisplace.LengthSquared() > maxDisplace.LengthSquared() {
			maxDisplace = flDisplace
		}
		if frDisplace.LengthSquared() > maxDisplace.LengthSquared() {
			maxDisplace = frDisplace
		}
		if blDisplace.LengthSquared() > maxDisplace.LengthSquared() {
			maxDisplace = blDisplace
		}
		if brDisplace.LengthSquared() > maxDisplace.LengthSquared() {
			maxDisplace = brDisplace
		}
		c.position = c.position.IncVec3(maxDisplace)
		if c.velocity.Y < 0 {
			c.velocity.Y = 0 // FIXME: should be constrained based on displace vector
		}
		c.angularVelocity.X = 0 // FIXME: should be constrained based on displace vector
		c.angularVelocity.Z = 0 // FIXME: should be constrained based on displace vector
	}
}

func (c *Car) checkWheel(wheel *Wheel) (math.Vec3, bool) {
	collided := false
	displacement := math.NullVec3()

	wheel.collisionCylinder.touching = false
	wheel.collisionCylinder.force = math.NullVec3()

	// test ground
	p1 := wheel.collisionCylinder.position.IncVec3(wheel.collisionCylinder.orientation.VectorY.Mul(wheel.collisionCylinder.radius))
	p2 := wheel.collisionCylinder.position.DecVec3(wheel.collisionCylinder.orientation.VectorY.Mul(wheel.collisionCylinder.radius))

	result, active := c.stage.CheckCollision(collision.MakeLine(p1, p2))
	if active {
		collided = true
		f := result.Normal().Resize(math.Abs32(result.BottomHeight()))
		displacement = displacement.IncVec3(f)
		wheel.collisionCylinder.touching = true
		wheel.collisionCylinder.force = wheel.collisionCylinder.force.IncVec3(f.Div(8))
	}

	p1 = wheel.collisionCylinder.position.IncVec3(wheel.collisionCylinder.orientation.VectorX.Mul(wheel.collisionCylinder.length / 2))
	p2 = wheel.collisionCylinder.position.DecVec3(wheel.collisionCylinder.orientation.VectorX.Mul(wheel.collisionCylinder.length / 2))
	result, active = c.stage.CheckCollision(collision.MakeLine(p1, p2))
	if active {
		fmt.Println("right collision!")
		collided = true
		f := result.Normal().Resize(math.Abs32(result.BottomHeight()))
		displacement = displacement.IncVec3(f)
		wheel.collisionCylinder.touching = true
		wheel.collisionCylinder.force = wheel.collisionCylinder.force.IncVec3(f.Div(8))
	}
	result, active = c.stage.CheckCollision(collision.MakeLine(p2, p1))
	if active {
		fmt.Println("left collision!")
		collided = true
		f := result.Normal().Resize(math.Abs32(result.BottomHeight()))
		displacement = displacement.IncVec3(f)
		wheel.collisionCylinder.touching = true
		wheel.collisionCylinder.force = wheel.collisionCylinder.force.IncVec3(f.Div(8))
	}

	p1 = wheel.collisionCylinder.position.IncVec3(wheel.collisionCylinder.orientation.VectorZ.Mul(wheel.collisionCylinder.radius))
	p2 = wheel.collisionCylinder.position.DecVec3(wheel.collisionCylinder.orientation.VectorZ.Mul(wheel.collisionCylinder.radius))
	result, active = c.stage.CheckCollision(collision.MakeLine(p1, p2))
	if active {
		fmt.Println("back collision!")
		collided = true
		f := result.Normal().Resize(math.Abs32(result.BottomHeight()))
		displacement = displacement.IncVec3(f)
		wheel.collisionCylinder.touching = true
		wheel.collisionCylinder.force = wheel.collisionCylinder.force.IncVec3(f.Div(8))
	}
	result, active = c.stage.CheckCollision(collision.MakeLine(p2, p1))
	if active {
		fmt.Println("front collision!")
		collided = true
		f := result.Normal().Resize(math.Abs32(result.BottomHeight()))
		displacement = displacement.IncVec3(f)
		wheel.collisionCylinder.touching = true
		wheel.collisionCylinder.force = wheel.collisionCylinder.force.IncVec3(f.Div(8))
	}

	return displacement, collided
}

func NewWheel(node *stream.Node, driving bool) *Wheel {
	return &Wheel{
		// FIXME: Not accurate, take body node into consideration and fix Y
		relativePosition:    math.MakeVec3(node.Matrix.M14-node.Parent.Matrix.M14, 0.0, node.Matrix.M34-node.Parent.Matrix.M34),
		driving:             driving,
		originalModelMatrix: node.Matrix,
		collisionCylinder: Cylinder{
			length: 0.4,
			radius: 0.2,
		},
		Node: node,
	}
}

type Wheel struct {
	relativePosition math.Vec3
	orientation      Orientation

	originalModelMatrix math.Mat4x4
	Node                *stream.Node
	driving             bool
	TurnAngle           float32
	RotationAngle       float32
	grounded            bool

	collisionCylinder Cylinder
}

func (w *Wheel) updateBasics(car *Car) {
	w.orientation = car.orientation
	w.orientation.Rotate(w.orientation.VectorY.Mul(w.TurnAngle * math.Pi / 180.0))
	transformedRelativePosition := car.orientation.MulVec3(w.relativePosition)

	w.collisionCylinder.orientation = w.orientation
	w.collisionCylinder.position = car.position.IncVec3(transformedRelativePosition)
}

func (w *Wheel) Update(car *Car) {
	w.Node.Matrix = math.Mat4x4MulMany(
		w.originalModelMatrix,
		math.RotationMat4x4(w.TurnAngle, 0.0, 1.0, 0.0),
		math.RotationMat4x4(w.RotationAngle, 1.0, 0.0, 0.0),
	)

	if !w.collisionCylinder.touching {
		return
	}

	// TODO
	// if !w.grounded {
	// 	return
	// }

	transformedRelativePosition := car.orientation.MulVec3(w.relativePosition)

	velocity := car.velocity
	velocity = velocity.IncVec3(math.Vec3CrossProduct(
		car.angularVelocity, transformedRelativePosition,
	))

	const wheelFriction = 0.01 // 0.03 // 0.005
	flWheelForce := math.NullVec3()
	// TODO: Take into consideration wheel spin
	flWheelForce = flWheelForce.DecVec3(w.orientation.VectorX.Mul(
		math.Vec3DotProduct(w.orientation.VectorX, velocity) * 0.2, // FIXME
	))
	if w.driving {
		flWheelForce = flWheelForce.IncVec3(w.orientation.VectorZ.Mul(car.Acceleration))
	}
	if flWheelForce.Length() > wheelFriction {
		flWheelForce = flWheelForce.Resize(wheelFriction)
	}
	if w.collisionCylinder.touching {
		flWheelForce = flWheelForce.IncVec3(w.collisionCylinder.force)
	}
	// fmt.Printf("wheel force: %.4f, %.4f, %.4f\n", flWheelForce.X, flWheelForce.Y, flWheelForce.Z)

	car.acceleration = car.acceleration.IncVec3(flWheelForce)
	car.angularAcceleration = car.angularAcceleration.IncVec3(math.Vec3CrossProduct(
		transformedRelativePosition, flWheelForce, // FIXME should use wheel relative position
	).Div(transformedRelativePosition.Length())) // FIXME: Use moment of intertia formulas

	// FIXME: If wheel is sliding, there is no controlled directional component
}

type Cylinder struct {
	position    math.Vec3
	orientation Orientation
	length      float32
	radius      float32
	touching    bool
	force       math.Vec3
}

type Orientation struct {
	VectorX math.Vec3
	VectorY math.Vec3
	VectorZ math.Vec3
}

func (o *Orientation) Rotate(rotation math.Vec3) {
	length := rotation.Length()
	if length < 0.00001 {
		return
	}
	matrix := math.RotationMat4x4(length*(180.0/math.Pi), rotation.X, rotation.Y, rotation.Z)
	orientationMatrix := math.VectorMat4x4(
		o.VectorX,
		o.VectorY,
		o.VectorZ,
		math.NullVec3(),
	)
	result := matrix.MulMat4x4(orientationMatrix)
	o.VectorX = math.MakeVec3(result.M11, result.M21, result.M31)
	o.VectorY = math.MakeVec3(result.M12, result.M22, result.M32)
	o.VectorZ = math.MakeVec3(result.M13, result.M23, result.M33)
}

func (o Orientation) MulVec3(vec math.Vec3) math.Vec3 {
	result := math.VectorMat4x4(
		o.VectorX,
		o.VectorY,
		o.VectorZ,
		math.NullVec3(),
	).MulVec4(math.MakeVec4(vec.X, vec.Y, vec.Z, 1.0))
	return math.MakeVec3(result.X, result.Y, result.Z)
}
