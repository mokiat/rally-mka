RESOURCES_DIR=resources
ASSETS_DIR=assets

SKYBOX_RESOURCES_DIR=$(RESOURCES_DIR)/skyboxes
SKYBOX_ASSETS_DIR=$(ASSETS_DIR)/skyboxes
SHADER_RESOURCES_DIR=$(RESOURCES_DIR)/shaders
PROGRAM_ASSETS_DIR=$(ASSETS_DIR)/programs
MESH_RESOURCES_DIR=$(RESOURCES_DIR)/meshes
MESH_ASSETS_DIR=$(ASSETS_DIR)/meshes

.PHONY: assets
assets: programs meshes skyboxes

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

.PHONY: skyboxes
skyboxes: \
	$(SKYBOX_ASSETS_DIR) \
	$(SKYBOX_ASSETS_DIR)/city.dat

$(SKYBOX_ASSETS_DIR):
	mkdir -p "$(SKYBOX_ASSETS_DIR)"

$(SKYBOX_ASSETS_DIR)/city.dat: \
	$(SKYBOX_RESOURCES_DIR)/city/front.png \
	$(SKYBOX_RESOURCES_DIR)/city/back.png \
	$(SKYBOX_RESOURCES_DIR)/city/left.png \
	$(SKYBOX_RESOURCES_DIR)/city/right.png \
	$(SKYBOX_RESOURCES_DIR)/city/top.png \
	$(SKYBOX_RESOURCES_DIR)/city/bottom.png
	rallygen generate-cubemap $+ $@
