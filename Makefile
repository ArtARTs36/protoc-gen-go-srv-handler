test:
	mkdir -p ./srv/
	go install ./ && \
        protoc --go_out=. --go_opt=paths=source_relative \
        --go-grpc_out=./ --go-grpc_opt=paths=source_relative \
        --go-srv-handler_out=./srv/ --go-srv-handler_opt=paths=source_relative \
        --go-srv-handler_opt=out_dir=./srv \
        --go-srv-handler_opt=overwrite=true \
        --go-srv-handler_opt=pkg_naming=without_service_suffix \
        --go-srv-handler_opt=srv_naming=just_service \
        userapi/*.proto

lint:
	golangci-lint run --fix
