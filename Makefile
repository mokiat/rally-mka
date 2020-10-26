.PHONY: assets
assets:
	go run cmd/rallypack/main.go

.PHONY: play
play:
	go run cmd/rallymka/main.go
