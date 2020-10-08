Generate:
	@echo 'Build GRPC'
	protoc -I generate/gatewary/ generate/gatewary/*.proto --go_out=plugins=grpc:generate/gatewary/.
