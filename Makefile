PACKAGES := $$(go list ./... | grep -v /vendor/)

test:
	@echo "Running tests..."
	go test $(PACKAGES)
