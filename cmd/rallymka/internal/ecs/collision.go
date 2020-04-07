package ecs

type CollisionComponent struct {
	RestitutionCoef float32
	CollisionShape  interface{}
}
