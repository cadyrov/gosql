domain:
	 export `cat .local_env` && go run ./cmd/main.go --template=./resources/domain.tmpl --result=./domain.go
store:
	export `cat .local_env` && go run ./cmd/main.go --template=./resources/store.tmpl --result=./store.go
service:
	export `cat .local_env` && go run ./cmd/main.go --template=./resources/service.tmpl --result=./service.go
bundle: domain store service