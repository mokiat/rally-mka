RESOURCES_DIR=resources
ASSETS_DIR=assets
LEVEL_RESOURCES_DIR=$(RESOURCES_DIR)/levels
LEVEL_ASSETS_DIR=$(ASSETS_DIR)/levels
SHADER_RESOURCES_DIR=$(RESOURCES_DIR)/shaders
PROGRAM_ASSETS_DIR=$(ASSETS_DIR)/programs
MODEL_RESOURCES_DIR=$(RESOURCES_DIR)/models
MODEL_ASSETS_DIR=$(ASSETS_DIR)/models
MESH_RESOURCES_DIR=$(RESOURCES_DIR)/meshes
MESH_ASSETS_DIR=$(ASSETS_DIR)/meshes
CUBE_TEX_RESOURCES_DIR=$(RESOURCES_DIR)/textures/cube
CUBE_TEX_ASSETS_DIR=$(ASSETS_DIR)/textures/cube
TWOD_TEX_RESOURCES_DIR=$(RESOURCES_DIR)/textures/twod
TWOD_TEX_ASSETS_DIR=$(ASSETS_DIR)/textures/twod

.PHONY: assets
assets: levels programs models meshes cubetextures twodtextures

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

.PHONY: programs
programs: \
	$(PROGRAM_ASSETS_DIR) \
	$(PROGRAM_ASSETS_DIR)/diffuse.dat \
	$(PROGRAM_ASSETS_DIR)/skybox.dat

$(PROGRAM_ASSETS_DIR):
	mkdir -p "$(PROGRAM_ASSETS_DIR)"

$(PROGRAM_ASSETS_DIR)/diffuse.dat: \
	$(SHADER_RESOURCES_DIR)/diffuse.vert \
	$(SHADER_RESOURCES_DIR)/diffuse.frag
	rallygen program $+ $@

$(PROGRAM_ASSETS_DIR)/skybox.dat: \
	$(SHADER_RESOURCES_DIR)/skybox.vert \
	$(SHADER_RESOURCES_DIR)/skybox.frag
	rallygen program $+ $@

.PHONY: models
models: \
	$(MODEL_ASSETS_DIR) \
	$(MODEL_ASSETS_DIR)/tree.dat \
	$(MODEL_ASSETS_DIR)/lamp.dat \
	$(MODEL_ASSETS_DIR)/finish.dat \
	$(MODEL_ASSETS_DIR)/suv.dat \
	$(MODEL_ASSETS_DIR)/hatch.dat \
	$(MODEL_ASSETS_DIR)/truck.dat

$(MODEL_ASSETS_DIR):
	mkdir -p "$(MODEL_ASSETS_DIR)"

$(MODEL_ASSETS_DIR)/tree.dat: \
	$(MODEL_RESOURCES_DIR)/tree.json
	rallygen model $+ $@

$(MODEL_ASSETS_DIR)/lamp.dat: \
	$(MODEL_RESOURCES_DIR)/lamp.json
	rallygen model $+ $@

$(MODEL_ASSETS_DIR)/finish.dat: \
	$(MODEL_RESOURCES_DIR)/finish.json
	rallygen model $+ $@

$(MODEL_ASSETS_DIR)/hatch.dat: \
	$(MODEL_RESOURCES_DIR)/hatch.json
	rallygen model $+ $@

$(MODEL_ASSETS_DIR)/suv.dat: \
	$(MODEL_RESOURCES_DIR)/suv.json
	rallygen model $+ $@

$(MODEL_ASSETS_DIR)/truck.dat: \
	$(MODEL_RESOURCES_DIR)/truck.json
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

.PHONY: cubetextures
cubetextures:	\
	$(CUBE_TEX_ASSETS_DIR) \
	$(CUBE_TEX_ASSETS_DIR)/city.dat

$(CUBE_TEX_ASSETS_DIR):
	mkdir -p "$(CUBE_TEX_ASSETS_DIR)"

$(CUBE_TEX_ASSETS_DIR)/city.dat: \
	$(CUBE_TEX_RESOURCES_DIR)/city_front.png \
	$(CUBE_TEX_RESOURCES_DIR)/city_back.png \
	$(CUBE_TEX_RESOURCES_DIR)/city_left.png \
	$(CUBE_TEX_RESOURCES_DIR)/city_right.png \
	$(CUBE_TEX_RESOURCES_DIR)/city_top.png \
	$(CUBE_TEX_RESOURCES_DIR)/city_bottom.png
	rallygen cubetex --dimension 512 $+ $@

.PHONY: twodtextures
twodtextures: \
	$(TWOD_TEX_ASSETS_DIR) \
	$(TWOD_TEX_ASSETS_DIR)/loading.dat \
	$(TWOD_TEX_ASSETS_DIR)/tree.dat \
	$(TWOD_TEX_ASSETS_DIR)/lamp.dat \
	$(TWOD_TEX_ASSETS_DIR)/finish.dat \
	$(TWOD_TEX_ASSETS_DIR)/hatch_body.dat \
	$(TWOD_TEX_ASSETS_DIR)/hatch_wheel.dat \
	$(TWOD_TEX_ASSETS_DIR)/suv_body.dat \
	$(TWOD_TEX_ASSETS_DIR)/suv_wheel.dat \
	$(TWOD_TEX_ASSETS_DIR)/truck_body.dat \
	$(TWOD_TEX_ASSETS_DIR)/truck_wheel.dat \
	$(TWOD_TEX_ASSETS_DIR)/concrete.dat \
	$(TWOD_TEX_ASSETS_DIR)/road.dat \
	$(TWOD_TEX_ASSETS_DIR)/barrier.dat \
	$(TWOD_TEX_ASSETS_DIR)/grass.dat \
	$(TWOD_TEX_ASSETS_DIR)/gravel.dat

$(TWOD_TEX_ASSETS_DIR):
	mkdir -p "$(TWOD_TEX_ASSETS_DIR)"

$(TWOD_TEX_ASSETS_DIR)/loading.dat: \
	$(TWOD_TEX_RESOURCES_DIR)/loading.png
	rallygen twodtex $+ $@

$(TWOD_TEX_ASSETS_DIR)/tree.dat: \
	$(TWOD_TEX_RESOURCES_DIR)/tree.png
	rallygen twodtex $+ $@

$(TWOD_TEX_ASSETS_DIR)/lamp.dat: \
	$(TWOD_TEX_RESOURCES_DIR)/lamp.png
	rallygen twodtex $+ $@

$(TWOD_TEX_ASSETS_DIR)/finish.dat: \
	$(TWOD_TEX_RESOURCES_DIR)/finish.png
	rallygen twodtex $+ $@

$(TWOD_TEX_ASSETS_DIR)/hatch_body.dat: \
	$(TWOD_TEX_RESOURCES_DIR)/hatch_body.png
	rallygen twodtex $+ $@

$(TWOD_TEX_ASSETS_DIR)/hatch_wheel.dat: \
	$(TWOD_TEX_RESOURCES_DIR)/hatch_wheel.png
	rallygen twodtex $+ $@

$(TWOD_TEX_ASSETS_DIR)/suv_body.dat: \
	$(TWOD_TEX_RESOURCES_DIR)/suv_body.png
	rallygen twodtex $+ $@

$(TWOD_TEX_ASSETS_DIR)/suv_wheel.dat: \
	$(TWOD_TEX_RESOURCES_DIR)/suv_wheel.png
	rallygen twodtex $+ $@

$(TWOD_TEX_ASSETS_DIR)/truck_body.dat: \
	$(TWOD_TEX_RESOURCES_DIR)/truck_body.png
	rallygen twodtex $+ $@

$(TWOD_TEX_ASSETS_DIR)/truck_wheel.dat: \
	$(TWOD_TEX_RESOURCES_DIR)/truck_wheel.png
	rallygen twodtex $+ $@

$(TWOD_TEX_ASSETS_DIR)/concrete.dat: \
	$(TWOD_TEX_RESOURCES_DIR)/concrete.png
	rallygen twodtex $+ $@

$(TWOD_TEX_ASSETS_DIR)/road.dat: \
	$(TWOD_TEX_RESOURCES_DIR)/road.png
	rallygen twodtex $+ $@

$(TWOD_TEX_ASSETS_DIR)/barrier.dat: \
	$(TWOD_TEX_RESOURCES_DIR)/barrier.png
	rallygen twodtex $+ $@

$(TWOD_TEX_ASSETS_DIR)/grass.dat: \
	$(TWOD_TEX_RESOURCES_DIR)/grass.png
	rallygen twodtex $+ $@

$(TWOD_TEX_ASSETS_DIR)/gravel.dat: \
	$(TWOD_TEX_RESOURCES_DIR)/gravel.png
	rallygen twodtex $+ $@
