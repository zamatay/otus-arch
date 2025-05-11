.DEFAULT_GOAL := gen

gen:
	protoc --go_out=. --go_opt=paths=import \
	--go-grpc_out=. --go-grpc_opt=paths=import proto/balances.proto
gen_php:
	protoc --proto_path=proto --php_out=php --php-grpc_out=php balances.proto
gen_rest:
	protoc --go_out=. --go_opt=paths=import \
	--go-grpc_out=. --go-grpc_opt=paths=import  \
	--grpc-gateway_out=. --grpc-gateway_opt generate_unbound_methods=true --openapiv2_out . \
	proto/monolit.proto
