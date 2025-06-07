.PHONY: build

all: build publish

tidy:
	@gum spin --show-error --title "[VERIFIER] Running go mod tidy..." -- sh -c "go mod tidy"
	@gum log --level info --time kitchen "[VERIFIER] Completed go mod tidy successfully."

build: tidy
	@tinygo build -o verifier.wasm -target wasip1 -buildmode=c-shared main.go
	@gum log --level info --time kitchen "[VERIFIER] Completed tinygo build successfully."

publish: build
	@gum spin --show-error --title "[VERIFIER] Uploading verifier.wasm to r2" -- sh -c "rclone copy ./verifier.wasm r2:cdn/bin/"
	@gum log --level info --time kitchen "[VERIFIER] Completed rclone upload successfully."


