package entities

import (
	"github.com/mokiat/go-whiskey/math"
	"github.com/mokiat/rally-mka/collision"
	"github.com/mokiat/rally-mka/render"
)

const minTurn = -30
const maxTurn = 30
const maxAcceleration = 0.15
const gravity = 0.15
const rotationFriction = 60.0 / 100.0
const rotationFriction2 = 1.0 - 60.0/100.0
const suspensionLength = 4.0
const speedFriction = 99.0 / 100.0
const speedFriction2 = 1.0 - (99.0 / 100.0)
const wheelFriction = 99.0 / 100.0

func NewCarExtendedModel() *CarExtendedModel {
	return &CarExtendedModel{
		CarModelSimple: NewCarModelSimple(),
	}
}

type CarExtendedModel struct {
	*CarModelSimple
	turn         int
	acceleration float32
}

func (m *CarExtendedModel) DrawMe(renderer *render.Renderer) {
	m.Draw(renderer, float32(m.turn))
}

func (m *CarExtendedModel) Frame(forward, back, left, right, brake bool, gameMap Map) {
	if left && (m.turn < maxTurn) {
		m.turn += 2
	}
	if right && (m.turn > minTurn) {
		m.turn -= 2
	}

	if left == right {
		if m.turn > 0 {
			m.turn--
		}
		if m.turn < 0 {
			m.turn++
		}
	}
	m.acceleration = 0.0
	if forward {
		m.acceleration = 0.15
	}
	if back {
		m.acceleration = -0.1
	}

	m.CheckMove(gameMap, float32(m.turn), m.acceleration, brake)
	m.CheckPoint(gameMap)
}

type Element struct {
	RenderMesh *render.Mesh

	Location math.Vec3
	Position math.Vec3
	Real     math.Vec3

	IsTouched  bool
	WheelAngle float32
	TurnKoef   float32

	CheckX float32
	CheckY float32
	CheckZ float32
}

func NewCarModelSimple() *CarModelSimple {
	return &CarModelSimple{
		Laps:    1,
		vectorX: math.BaseVec3X(),
		vectorY: math.BaseVec3Y(),
		vectorZ: math.BaseVec3Z(),
	}
}

type CarModelSimple struct {
	Position     math.Vec3
	lastPosition math.Vec3

	speed    math.Vec3
	rotation math.Vec3
	rotForce math.Vec3
	force    math.Vec3

	vectorX math.Vec3
	vectorY math.Vec3
	vectorZ math.Vec3

	body    *Element
	wheelFL *Element
	wheelFR *Element
	wheelBL *Element
	wheelBR *Element

	Laps int
}

func (s *CarModelSimple) Load(path string) error {
	model := NewExtendedModel()
	if err := model.Load(path); err != nil {
		return err
	}
	model.EvaluateMinMax()

	s.body = s.getElement(model, "Car")
	s.wheelFL = s.getElement(model, "LF")
	s.wheelFR = s.getElement(model, "RF")
	s.wheelBL = s.getElement(model, "LB")
	s.wheelBR = s.getElement(model, "RB")

	s.wheelFL.CheckX = math.Abs32(model.MaxX-s.wheelFL.Position.X) + 4.0
	s.wheelFL.CheckY = math.Abs32(model.MinY - s.wheelFL.Position.Y)
	s.wheelFL.CheckZ = math.Abs32(model.MinZ-s.wheelFL.Position.Z) + 4.0
	s.wheelFL.TurnKoef = 180.0 / (float32(math.Pi) * s.wheelFL.CheckY)

	s.wheelFR.CheckX = math.Abs32(model.MinX-s.wheelFR.Position.X) + 4.0
	s.wheelFR.CheckY = math.Abs32(model.MinY - s.wheelFR.Position.Y)
	s.wheelFR.CheckZ = math.Abs32(model.MinZ-s.wheelFR.Position.Z) + 4.0
	s.wheelFR.TurnKoef = 180.0 / (float32(math.Pi) * s.wheelFR.CheckY)

	s.wheelBL.CheckX = math.Abs32(model.MaxX-s.wheelBL.Position.X) + 4.0
	s.wheelBL.CheckY = math.Abs32(model.MinY - s.wheelBL.Position.Y)
	s.wheelBL.CheckZ = math.Abs32(model.MaxZ-s.wheelBL.Position.Z) + 4.0
	s.wheelBL.TurnKoef = 180.0 / (float32(math.Pi) * s.wheelBL.CheckY)

	s.wheelBR.CheckX = math.Abs32(model.MinX-s.wheelBR.Position.X) + 4.0
	s.wheelBR.CheckY = math.Abs32(model.MinY - s.wheelBR.Position.Y)
	s.wheelBR.CheckZ = math.Abs32(model.MaxZ-s.wheelBR.Position.Z) + 4.0
	s.wheelBR.TurnKoef = 180.0 / (float32(math.Pi) * s.wheelBR.CheckY)

	return nil
}

func (s *CarModelSimple) Generate() {
	s.body.RenderMesh.Generate()
	s.wheelFL.RenderMesh.Generate()
	s.wheelFR.RenderMesh.Generate()
	s.wheelBL.RenderMesh.Generate()
	s.wheelBR.RenderMesh.Generate()
}

func (s *CarModelSimple) Draw(renderer *render.Renderer, turn float32) {
	carModelMatrix := math.VectorMat4x4(
		s.vectorX.Resize(1.0),
		s.vectorY.Resize(1.0),
		s.vectorZ.Resize(1.0),
		math.NullVec3(),
	)

	// Car Body
	{
		oldModelMatrix := renderer.ModelMatrix()
		carBodyMatrix := carModelMatrix
		carBodyMatrix = carBodyMatrix.MulMat4x4(math.TranslationMat4x4(s.body.Location.X, s.body.Location.Y, s.body.Location.Z))
		renderer.SetModelMatrix(carBodyMatrix)
		renderer.Render(s.body.RenderMesh, renderer.TextureMaterial())
		renderer.SetModelMatrix(oldModelMatrix)
	}

	// Front Left Wheel
	{
		oldModelMatrix := renderer.ModelMatrix()
		modelMatrix := carModelMatrix
		modelMatrix = modelMatrix.MulMat4x4(math.TranslationMat4x4(s.wheelFL.Real.X, s.wheelFL.Real.Y, s.wheelFL.Real.Z))
		modelMatrix = modelMatrix.MulMat4x4(math.RotationMat4x4(turn, 0.0, 1.0, 0.0))
		modelMatrix = modelMatrix.MulMat4x4(math.RotationMat4x4(s.wheelFL.WheelAngle, 1.0, 0.0, 0.0))
		renderer.SetModelMatrix(modelMatrix)
		renderer.Render(s.wheelFL.RenderMesh, renderer.TextureMaterial())
		renderer.SetModelMatrix(oldModelMatrix)
	}

	// Front Right Wheel
	{
		oldModelMatrix := renderer.ModelMatrix()
		modelMatrix := carModelMatrix
		modelMatrix = modelMatrix.MulMat4x4(math.TranslationMat4x4(s.wheelFR.Real.X, s.wheelFR.Real.Y, s.wheelFR.Real.Z))
		modelMatrix = modelMatrix.MulMat4x4(math.RotationMat4x4(turn, 0.0, 1.0, 0.0))
		modelMatrix = modelMatrix.MulMat4x4(math.RotationMat4x4(s.wheelFR.WheelAngle, 1.0, 0.0, 0.0))
		renderer.SetModelMatrix(modelMatrix)
		renderer.Render(s.wheelFR.RenderMesh, renderer.TextureMaterial())
		renderer.SetModelMatrix(oldModelMatrix)
	}

	// Rear Left Wheel
	{
		oldModelMatrix := renderer.ModelMatrix()
		modelMatrix := carModelMatrix
		modelMatrix = modelMatrix.MulMat4x4(math.TranslationMat4x4(s.wheelBL.Real.X, s.wheelBL.Real.Y, s.wheelBL.Real.Z))
		modelMatrix = modelMatrix.MulMat4x4(math.RotationMat4x4(s.wheelBL.WheelAngle, 1.0, 0.0, 0.0))
		renderer.SetModelMatrix(modelMatrix)
		renderer.Render(s.wheelBL.RenderMesh, renderer.TextureMaterial())
		renderer.SetModelMatrix(oldModelMatrix)
	}

	// Rear Right Wheel
	{
		oldModelMatrix := renderer.ModelMatrix()
		modelMatrix := carModelMatrix
		modelMatrix = modelMatrix.MulMat4x4(math.TranslationMat4x4(s.wheelBR.Real.X, s.wheelBR.Real.Y, s.wheelBR.Real.Z))
		modelMatrix = modelMatrix.MulMat4x4(math.RotationMat4x4(s.wheelBR.WheelAngle, 1.0, 0.0, 0.0))
		renderer.SetModelMatrix(modelMatrix)
		renderer.Render(s.wheelBR.RenderMesh, renderer.TextureMaterial())
		renderer.SetModelMatrix(oldModelMatrix)
	}
}

func (s *CarModelSimple) CheckMove(gameMap Map, turn, acceleration float32, brake bool) {
	s.lastPosition = s.Position

	directSpeed := acceleration * 20.0
	if s.wheelFL.IsTouched && s.wheelFR.IsTouched {
		var frontWheelSpeed float32
		if brake {
			s.speed = s.speed.Mul(wheelFriction)
			frontWheelSpeed = 0.0
		} else {
			matrix := math.RotationMat4x4(turn, s.vectorY.X, s.vectorY.Y, s.vectorY.Z)
			vectorZ2 := transformVec3(matrix, s.vectorZ)
			frontWheelSpeed = math.Vec3DotProduct(s.speed, vectorZ2)
		}
		if math.Abs32(directSpeed) > math.Abs32(frontWheelSpeed) {
			frontWheelSpeed = directSpeed
		}

		goForwardVector := s.vectorZ.Resize(acceleration)
		s.speed = s.speed.IncVec3(goForwardVector)
		s.wheelFL.WheelAngle += frontWheelSpeed * s.wheelFL.TurnKoef
		s.wheelFR.WheelAngle += frontWheelSpeed * s.wheelFR.TurnKoef
	} else {
		frontWheelSpeed := directSpeed
		if !s.wheelFL.IsTouched && !brake {
			s.wheelFL.WheelAngle += frontWheelSpeed * s.wheelFL.TurnKoef
		}
		if !s.wheelFR.IsTouched && !brake {
			s.wheelFR.WheelAngle += frontWheelSpeed * s.wheelFR.TurnKoef
		}
	}

	backWheelSpeed := math.Vec3DotProduct(s.speed, s.vectorZ)
	if !brake {
		if s.wheelBL.IsTouched {
			s.wheelBL.WheelAngle += backWheelSpeed * s.wheelBL.TurnKoef
		}
		if s.wheelBR.IsTouched {
			s.wheelBR.WheelAngle += backWheelSpeed * s.wheelBR.TurnKoef
		}
	}

	speed2 := backWheelSpeed + math.Vec3DotProduct(s.speed, s.vectorX)*math.Sin32(turn*math.Pi/180.0)
	var doTurn float32
	if !brake {
		doTurn = turn * speed2 / (20.0 + 2.0*speed2*speed2)
	} else {
		doTurn = 0
	}
	doTurn2 := turn * acceleration
	doTurn = (doTurn + doTurn2) / 2.0

	if (s.wheelFL.IsTouched || s.wheelFR.IsTouched) && (math.Abs32(doTurn) > 0.0001) {
		s.rotation = s.rotation.IncVec3(s.vectorY.Resize(doTurn))
		s.rotate(doTurn, s.vectorY)
	}

	s.CheckCollision(gameMap)
}

func (s *CarModelSimple) CheckPoint(gameMap Map) {
}

func (s *CarModelSimple) CheckCollision(gameMap Map) {
	s.force = math.Vec3{
		X: 0.0,
		Y: -gravity,
		Z: 0.0,
	}
	s.speed = s.speed.IncVec3(s.force)
	s.speed = s.speed.DecVec3(s.speed.Mul(speedFriction2))
	s.translate(s.speed)

	s.force = math.NullVec3()
	s.checkCar(gameMap)
	s.translate(s.force)

	dis := s.rotation.Length() / 10.0
	if dis > 0.0000001 {
		s.rotate(dis, s.rotation.Resize(1.0))
	}

	s.rotation = s.rotation.Mul(rotationFriction)
	s.rotation = s.rotation.DecVec3(s.rotation.Mul(rotationFriction2))

	s.force = math.NullVec3()
	s.rotForce = math.NullVec3()
	s.checkWheel(s.wheelFL, gameMap, s.Position)
	s.checkWheel(s.wheelFR, gameMap, s.Position)
	s.checkWheel(s.wheelBL, gameMap, s.Position)
	s.checkWheel(s.wheelBR, gameMap, s.Position)

	s.speed = s.speed.IncVec3(s.force)
	dis = s.rotForce.Length() / 10.0
	if dis > 0.0000001 {
		s.rotForce = s.rotForce.Resize(1.0)
		s.rotate(dis, s.rotForce)
	}
}

func (s *CarModelSimple) checkCar(gameMap Map) {
	s.checkWheelCollision(gameMap, s.wheelFL, s.vectorX.Inverse(), s.vectorZ.Inverse())
	s.checkWheelCollision(gameMap, s.wheelFR, s.vectorX, s.vectorZ.Inverse())
	s.checkWheelCollision(gameMap, s.wheelBL, s.vectorX.Inverse(), s.vectorZ)
	s.checkWheelCollision(gameMap, s.wheelBR, s.vectorX, s.vectorZ)
}

func (s *CarModelSimple) checkWheelCollision(gameMap Map, wheel *Element, dirX math.Vec3, dirZ math.Vec3) {
	pa := wheel.Position.IncVec3(s.Position)
	p2 := pa.DecVec3(dirX.Resize(wheel.CheckX))
	p1 := pa.IncVec3(dirX.Resize(20.0))

	result, active := gameMap.CheckCollisionWall(collision.MakeLine(p1, p2))
	if active {
		f := result.Normal().Resize(math.Abs32(result.BottomHeight()))
		s.translate(f)
		cross := math.Vec3CrossProduct(s.speed, result.Normal())
		f = math.Vec3CrossProduct(result.Normal(), cross).Resize(1.0)
		cross = f
		lng := math.Vec3DotProduct(f, s.speed)
		cross = cross.Resize(lng)
		s.speed = cross.Mul(90.0 / 100.0)
	}

	p2 = pa.IncVec3(dirZ.Resize(math.Abs32(wheel.CheckZ)))
	p1 = pa.DecVec3(dirZ.Resize(20.0))

	result, active = gameMap.CheckCollisionWall(collision.MakeLine(p1, p2))
	if active {
		f := result.Normal().Resize(math.Abs32(result.BottomHeight()))
		s.translate(f)
		cross := math.Vec3CrossProduct(s.speed, result.Normal())
		f = math.Vec3CrossProduct(result.Normal(), cross).Resize(1.0)
		cross = f
		lng := math.Vec3DotProduct(f, s.speed)
		cross = cross.Resize(lng)
		s.speed = cross.Mul(90.0 / 100.0)
	}
}

func (s *CarModelSimple) checkWheel(wheel *Element, gameMap Map, position math.Vec3) {
	wheel.IsTouched = false

	pa := wheel.Position.IncVec3(position)
	p1 := pa.IncVec3(s.vectorY.Resize(1000.0))
	p2 := pa.DecVec3(s.vectorY.Resize(wheel.CheckY + suspensionLength))

	result, active := gameMap.CheckCollisionGround(collision.MakeLine(p1, p2))
	if !active {
		wheel.Real = wheel.Location.DecCoords(0.0, suspensionLength, 0.0)
		return
	}

	dis := result.Intersection().DecVec3(p1).Length()
	if dis > 1000.0+wheel.CheckY {
		wheel.Real = wheel.Location.DecCoords(0.0, dis-(1000.0+wheel.CheckY), 0.0)
	} else {
		wheel.Real = wheel.Location
	}
	wheel.IsTouched = true

	if dis < 1000.0+wheel.CheckY {
		dis2 := math.Vec3DotProduct(result.Normal(), s.speed)
		f := result.Normal().Resize(dis2)
		s.speed = s.speed.DecVec3(f)
		f = s.vectorY.Resize(1000.0 + wheel.CheckY - dis)
		s.translate(f)
	}

	cross := math.Vec3CrossProduct(result.Normal(), wheel.Position.Inverse())
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
		cross = cross.Resize(koef.X * koef.X * (1.0 - (dis-(1000.0+wheel.CheckY))/suspensionLength) * 20)
		s.rotForce = s.rotForce.IncVec3(cross)
		s.rotation = s.rotation.IncVec3(cross)
	}
}

func (s *CarModelSimple) getElement(model *ExtendedModel, name string) *Element {
	searchIndex := model.FindObjectIndex(name, 0, false)
	object := model.GetObjectIndex(searchIndex)

	result := &Element{}
	result.Location = model.ObjectCenter(object)
	result.Position = model.ObjectCenter(object)
	model.DoCenter(object)
	result.RenderMesh = createRenderMesh(model, object)
	result.Real = result.Position
	result.WheelAngle = 0.0
	result.IsTouched = false
	return result
}

func (s *CarModelSimple) translate(translation math.Vec3) {
	s.Position = s.Position.IncVec3(translation)
}

func (s *CarModelSimple) rotate(angle float32, rotationVector math.Vec3) {
	matrix := math.RotationMat4x4(angle, rotationVector.X, rotationVector.Y, rotationVector.Z)

	s.vectorX = transformVec3(matrix, s.vectorX).Resize(1.0)
	s.vectorY = transformVec3(matrix, s.vectorY).Resize(1.0)
	s.vectorZ = transformVec3(matrix, s.vectorZ).Resize(1.0)

	s.body.Position = transformVec3(matrix, s.body.Position)
	s.wheelFL.Position = transformVec3(matrix, s.wheelFL.Position)
	s.wheelFR.Position = transformVec3(matrix, s.wheelFR.Position)
	s.wheelBL.Position = transformVec3(matrix, s.wheelBL.Position)
	s.wheelBR.Position = transformVec3(matrix, s.wheelBR.Position)
}

func transformVec3(matrix math.Mat4x4, vector math.Vec3) math.Vec3 {
	result := matrix.MulVec4(math.MakeVec4(vector.X, vector.Y, vector.Z, 1.0))
	return math.MakeVec3(result.X, result.Y, result.Z)
}
