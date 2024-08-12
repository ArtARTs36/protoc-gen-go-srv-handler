test:
	mkdir -p ./srv/
	go install ./cmd/protoc-gen-go-srv-handler && \
        protoc --go_out=. --go_opt=paths=source_relative \
        --go-grpc_out=./ --go-grpc_opt=paths=source_relative \
        --go-srv-handler_out=./srv/ --go-srv-handler_opt=paths=source_relative \
        --go-srv-handler_opt=out_dir=./srv \
        --go-srv-handler_opt=overwrite=true \
        userapi/*.proto
