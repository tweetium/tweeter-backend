alias gom="migrate -database $DATABASE_URL -path './migrations'"
alias gomcreate="migrate create -ext sql -dir migrations"

alias tests="gotestsum"
alias lints="golangci-lint run ./..."

alias gorun="go run main.go"
