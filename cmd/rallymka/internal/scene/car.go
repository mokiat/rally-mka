package scene

import (
	"github.com/mokiat/go-whiskey/math"
	"github.com/mokiat/rally-mka/cmd/rallymka/internal/stream"
	"github.com/mokiat/rally-mka/internal/engine/collision"
)

const (
	gravity             = 0.01
	wheelFriction       = 99.0 / 100.0
	speedFriction       = 99.0 / 100.0
	speedFriction2      = 1.0 - (99.0 / 100.0)
	rotationFriction    = 60.0 / 100.0
	rotationFriction2   = 1.0 - 60.0/100.0
	maxSuspensionLength = 0.4
	skidWheelSpeed      = 0.03
)

func NewCar(stage *Stage, model *stream.Model, position math.Vec3) *Car {
	bodyNode, _ := model.FindNode("body")
	flWheelNode, _ := model.FindNode("wheel_front_left")
	frWheelNode, _ := model.FindNode("wheel_front_right")
	blWheelNode, _ := model.FindNode("wheel_back_left")
	brWheelNode, _ := model.FindNode("wheel_back_right")

	car := &Car{
		stage: stage,
		node:  bodyNode,

		position: position,
		orientation: Orientation{
			VectorX: math.BaseVec3X(),
			VectorY: math.BaseVec3Y(),
			VectorZ: math.BaseVec3Z(),
		},
	}
	car.flWheel = NewWheel(car, flWheelNode, true)
	car.frWheel = NewWheel(car, frWheelNode, true)
	car.blWheel = NewWheel(car, blWheelNode, true)
	car.brWheel = NewWheel(car, brWheelNode, true)
	return car
}

type Car struct {
	stage   *Stage
	node    *stream.Node
	flWheel *Wheel
	frWheel *Wheel
	blWheel *Wheel
	brWheel *Wheel

	position        math.Vec3
	orientation     Orientation
	velocity        math.Vec3
	angularVelocity math.Vec3

	WheelAngle      float32
	Acceleration    float32
	HandbrakePulled bool
}

func (c *Car) Position() math.Vec3 {
	return c.position
}

func (c *Car) Update() {
	c.flWheel.turnAngle = c.WheelAngle
	c.frWheel.turnAngle = c.WheelAngle

	c.calculateWheelRotation(c.flWheel)
	c.calculateWheelRotation(c.frWheel)
	c.calculateWheelRotation(c.blWheel)
	c.calculateWheelRotation(c.brWheel)
	c.accelerate()
	c.translate()

	c.updateModelMatrix()
	c.flWheel.updateModelMatrix()
	c.frWheel.updateModelMatrix()
	c.blWheel.updateModelMatrix()
	c.brWheel.updateModelMatrix()
}

func (c *Car) calculateWheelRotation(wheel *Wheel) {
	// Handbrake locks all wheels
	if wheel.grounded && !c.HandbrakePulled {
		wheelOrientation := c.orientation
		wheelOrientation.Rotate(wheelOrientation.VectorY.Resize(wheel.turnAngle * math.Pi / 180))
		wheel.rotate(math.Vec3DotProduct(c.velocity, wheelOrientation.VectorZ))
	}

	// Add some wheel slip if we are accelerating
	if math.Abs32(c.Acceleration) > 0.0001 && wheel.driving {
		wheel.rotate(math.Signum32(c.Acceleration) * skidWheelSpeed)
	}
}

func (c *Car) accelerate() {
	shouldAccelerate :=
		(c.flWheel.driving && c.flWheel.grounded) ||
			(c.frWheel.driving && c.frWheel.grounded) ||
			(c.blWheel.driving && c.blWheel.grounded) ||
			(c.brWheel.driving && c.brWheel.grounded)

	if shouldAccelerate {
		acceleration := c.orientation.VectorZ.Resize(c.Acceleration)
		c.velocity = c.velocity.IncVec3(acceleration)
	}

	shouldDecelerate := c.HandbrakePulled &&
		(c.flWheel.grounded || c.frWheel.grounded || c.blWheel.grounded || c.brWheel.grounded)

	if shouldDecelerate {
		c.velocity = c.velocity.Mul(wheelFriction)
	}

	c.velocity = c.velocity.DecVec3(math.MakeVec3(0.0, gravity, 0.0))
	c.velocity = c.velocity.DecVec3(c.velocity.Mul(speedFriction2))

	// no idea how I originally go to this
	magicVelocity := math.Vec3DotProduct(c.velocity, c.orientation.VectorZ) + math.Vec3DotProduct(c.velocity, c.orientation.VectorX)*math.Sin32(c.WheelAngle*math.Pi/180.0)
	var turnAcceleration float32
	if !c.HandbrakePulled {
		turnAcceleration = c.WheelAngle * magicVelocity / (100.0 + 2.0*magicVelocity*magicVelocity)
	}
	turnAcceleration += c.WheelAngle * c.Acceleration * 0.6
	turnAcceleration /= 2.0

	if c.flWheel.grounded || c.frWheel.grounded {
		c.angularVelocity = c.angularVelocity.IncVec3(c.orientation.VectorY.Resize(turnAcceleration))
	}
	c.angularVelocity = c.angularVelocity.Mul(rotationFriction)
	c.angularVelocity = c.angularVelocity.DecVec3(c.angularVelocity.Mul(rotationFriction2))
}

func (c *Car) translate() {
	c.position = c.position.IncVec3(c.velocity)
	c.orientation.Rotate(c.angularVelocity)
	c.checkCollisions()
}

func (c *Car) checkCollisions() {
	c.checkWheelSideCollision(c.flWheel, c.orientation.VectorX.Inverse(), c.orientation.VectorZ.Inverse())
	c.checkWheelSideCollision(c.frWheel, c.orientation.VectorX, c.orientation.VectorZ.Inverse())
	c.checkWheelSideCollision(c.blWheel, c.orientation.VectorX.Inverse(), c.orientation.VectorZ)
	c.checkWheelSideCollision(c.brWheel, c.orientation.VectorX, c.orientation.VectorZ)
	c.checkWheelBottomCollision(c.flWheel)
	c.checkWheelBottomCollision(c.frWheel)
	c.checkWheelBottomCollision(c.blWheel)
	c.checkWheelBottomCollision(c.brWheel)
}

func (c *Car) checkWheelSideCollision(wheel *Wheel, dirX math.Vec3, dirZ math.Vec3) {
	pa := wheel.position()

	p2 := pa.DecVec3(dirX.Resize(wheel.length))
	p1 := pa.IncVec3(dirX.Resize(1.2))
	result, active := c.stage.CheckCollision(collision.MakeLine(p1, p2))
	if active {
		f := result.Normal().Resize(math.Abs32(result.BottomHeight()))
		c.position = c.position.IncVec3(f)
	}

	p2 = pa.IncVec3(dirZ.Resize(math.Abs32(wheel.radius)))
	p1 = pa.DecVec3(dirZ.Resize(1.2))
	result, active = c.stage.CheckCollision(collision.MakeLine(p1, p2))
	if active {
		f := result.Normal().Resize(math.Abs32(result.BottomHeight()))
		c.position = c.position.IncVec3(f)
	}
}

func (c *Car) checkWheelBottomCollision(wheel *Wheel) {
	wheel.grounded = false

	pa := c.position.IncVec3(c.orientation.MulVec3(wheel.anchorPosition))
	p1 := pa.IncVec3(c.orientation.VectorY.Resize(3.6))
	p2 := pa.DecVec3(c.orientation.VectorY.Resize(wheel.radius + maxSuspensionLength))

	result, active := c.stage.CheckCollision(collision.MakeLine(p1, p2))
	if !active {
		wheel.suspensionLength = maxSuspensionLength
		return
	}

	dis := result.Intersection().DecVec3(p1).Length()
	if dis > 3.6+wheel.radius {
		wheel.suspensionLength = dis - (3.6 + wheel.radius)
	} else {
		wheel.suspensionLength = 0.0
	}
	wheel.grounded = true

	if dis < 3.6+wheel.radius {
		dis2 := math.Vec3DotProduct(result.Normal(), c.velocity)
		f := result.Normal().Resize(dis2)
		c.velocity = c.velocity.DecVec3(f)
		f = c.orientation.VectorY.Resize(3.6 + wheel.radius - dis)
		c.position = c.position.IncVec3(f)
	}

	relativePosition := wheel.position().DecVec3(c.position)
	cross := math.Vec3CrossProduct(result.Normal(), relativePosition.Inverse())
	cross = cross.Resize(1.0)

	koef := math.NullVec2()
	koef.Y = math.Vec3DotProduct(result.Intersection().Inverse(), result.Intersection())
	tmp := result.Intersection().Inverse().Length()
	koef.X = math.Sqrt32(tmp*tmp - koef.Y*koef.Y)

	if koef.Length() > 0.0000001 {
		koef = koef.Resize(1.0)
	} else {
		koef.X = 1.0
		koef.Y = 0.0
	}

	if math.Abs32(koef.X) > 0.0000001 {
		cross = cross.Resize(koef.X * koef.X * (1.0 - (dis-(3.6+wheel.radius))/maxSuspensionLength) / 10)
		c.angularVelocity = c.angularVelocity.IncVec3(cross)
	}
}

func (c *Car) updateModelMatrix() {
	c.node.Matrix = math.VectorMat4x4(
		c.orientation.VectorX,
		c.orientation.VectorY,
		c.orientation.VectorZ,
		c.position,
	)
}

func NewWheel(car *Car, node *stream.Node, driving bool) *Wheel {
	return &Wheel{
		car:     car,
		node:    node,
		driving: driving,

		length: 0.4, // FIXME: Get from model
		radius: 0.2, // FIXME: Get from model

		anchorPosition:   node.Matrix.Translation().DecVec3(car.node.Matrix.Translation()),
		suspensionLength: 0.0,
	}
}

type Wheel struct {
	car     *Car
	node    *stream.Node
	driving bool

	length float32
	radius float32

	anchorPosition math.Vec3

	suspensionLength float32
	grounded         bool

	turnAngle     float32
	rotationAngle float32
}

func (w *Wheel) position() math.Vec3 {
	return w.car.position.
		IncVec3(w.car.orientation.MulVec3(w.anchorPosition)).
		IncVec3(w.car.orientation.MulVec3(math.MakeVec3(0.0, -w.suspensionLength, 0.0)))
}

func (w *Wheel) rotate(speed float32) {
	rotationCoefficient := float32(180.0 / (math.Pi * w.radius))
	w.rotationAngle += speed * rotationCoefficient
}

func (w *Wheel) updateModelMatrix() {
	w.node.Matrix = math.Mat4x4MulMany(
		math.TranslationMat4x4(w.anchorPosition.X, w.anchorPosition.Y-w.suspensionLength, w.anchorPosition.Z),
		math.RotationMat4x4(w.turnAngle, 0.0, 1.0, 0.0),
		math.RotationMat4x4(w.rotationAngle, 1.0, 0.0, 0.0),
	)
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
