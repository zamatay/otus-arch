.DEFAULT_GOAL := gen

gen_rest:
	protoc --go_out=. --go_opt=paths=import \
	--go-grpc_out=. --go-grpc_opt=paths=import  \
	--grpc-gateway_out=. --grpc-gateway_opt generate_unbound_methods=true --openapiv2_out . \
	proto/monolit.proto

	protoc --go_out=. --go_opt=paths=import \
	--go-grpc_out=. --go-grpc_opt=paths=import  \
	--grpc-gateway_out=. --grpc-gateway_opt generate_unbound_methods=true --openapiv2_out . \
	proto/counter.proto
