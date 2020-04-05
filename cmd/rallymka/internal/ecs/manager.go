package ecs

// FIXME: This whole manager is suboptimal and open to bugs
// (for example: deleting an entity while iterating)

func NewManager() *Manager {
	return &Manager{
		entities: make([]*Entity, 0),
	}
}

type Manager struct {
	entities []*Entity
}

func (m *Manager) Entities() []*Entity {
	return m.entities
}

func (m *Manager) CreateEntity() *Entity {
	entity := &Entity{}
	m.entities = append(m.entities, entity)
	return entity
}

func (m *Manager) DeleteEntity(entity *Entity) {
	index := -1
	for i, candidate := range m.entities {
		if candidate == entity {
			index = i
			break
		}
	}
	if index != -1 {
		count := len(m.entities)
		m.entities[index] = m.entities[count-1]
		m.entities = m.entities[:count-1]
	}
}
