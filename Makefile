RESOURCES_DIR=resources
ASSETS_DIR=assets

SKYBOX_RESOURCES_DIR=$(RESOURCES_DIR)/skyboxes
SKYBOX_ASSETS_DIR=$(ASSETS_DIR)/skyboxes

.PHONY: assets
assets: skyboxes

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
