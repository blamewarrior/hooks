PACKAGES := $$(go list ./... | grep -v /vendor/ | grep -v /cmd/)

test:
	@echo "Running tests..."
	DB_USER=postgres DB_NAME=bw_users_test go test $(PACKAGES)
