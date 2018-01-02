PACKAGES := $$(go list ./... | grep -v /vendor/ | grep -v /cmd/)

test:
	@echo "Running tests..."
	go test $(PACKAGES)
