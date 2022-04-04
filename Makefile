lint:
	gofmt -s -w ./ && golangci-lint run
domain:
	 export `cat .local_env` && go run ./cmd/main.go --template=./resources/psql/domain.tmpl --result=./result/domain.go
store:
	export `cat .local_env` && go run ./cmd/main.go --template=./resources/psql/store.tmpl --result=./result/store.go
service:
	export `cat .local_env` && go run ./cmd/main.go --template=./resources/psql/service.tmpl --result=./result/service.go
bundle: domain store service