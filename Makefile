generate-distributed-client:
	oapi-codegen --config pkg/distributed/schema/generator_config.yml pkg/distributed/schema/swagger.yml

generate-distributed-mocks:
	mockgen -source=pkg/distributed/models/interfaces.go -destination=pkg/distributed/models/mocks/mocks.gen.go IDealer
