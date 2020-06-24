RESOURCES_DIR=resources
ASSETS_DIR=assets
LEVEL_RESOURCES_DIR=$(RESOURCES_DIR)/levels
LEVEL_ASSETS_DIR=$(ASSETS_DIR)/levels
MODEL_RESOURCES_DIR=$(RESOURCES_DIR)/models
MODEL_ASSETS_DIR=$(ASSETS_DIR)/models
MESH_RESOURCES_DIR=$(RESOURCES_DIR)/meshes
MESH_ASSETS_DIR=$(ASSETS_DIR)/meshes

.PHONY: assets
assets: levels models meshes
	go run cmd/rallypack/main.go

.PHONY: levels
levels: \
	$(LEVEL_ASSETS_DIR) \
	$(LEVEL_ASSETS_DIR)/forest.dat \
	$(LEVEL_ASSETS_DIR)/highway.dat \
	$(LEVEL_ASSETS_DIR)/playground.dat

$(LEVEL_ASSETS_DIR):
	mkdir -p $(LEVEL_ASSETS_DIR)

$(LEVEL_ASSETS_DIR)/forest.dat: \
	$(LEVEL_RESOURCES_DIR)/forest.json
	rallygen level $+ $@

$(LEVEL_ASSETS_DIR)/highway.dat: \
	$(LEVEL_RESOURCES_DIR)/highway.json
	rallygen level $+ $@

$(LEVEL_ASSETS_DIR)/playground.dat: \
	$(LEVEL_RESOURCES_DIR)/playground.json
	rallygen level $+ $@

.PHONY: models
models: \
	$(MODEL_ASSETS_DIR) \
	$(MODEL_ASSETS_DIR)/tree.dat \
	$(MODEL_ASSETS_DIR)/finish.dat

$(MODEL_ASSETS_DIR):
	mkdir -p "$(MODEL_ASSETS_DIR)"

$(MODEL_ASSETS_DIR)/tree.dat: \
	$(MODEL_RESOURCES_DIR)/tree.json
	rallygen model $+ $@

$(MODEL_ASSETS_DIR)/finish.dat: \
	$(MODEL_RESOURCES_DIR)/finish.json
	rallygen model $+ $@

.PHONY: meshes
meshes: \
	$(MESH_ASSETS_DIR) \
	$(MESH_ASSETS_DIR)/quad.dat \
	$(MESH_ASSETS_DIR)/skybox.dat

$(MESH_ASSETS_DIR):
	mkdir -p "$(MESH_ASSETS_DIR)"

$(MESH_ASSETS_DIR)/quad.dat: \
	$(MESH_RESOURCES_DIR)/quad.json
	rallygen mesh $+ $@

$(MESH_ASSETS_DIR)/skybox.dat: \
	$(MESH_RESOURCES_DIR)/skybox.json
	rallygen mesh $+ $@
