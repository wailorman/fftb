generate-distributed-client:
	oapi-codegen --config pkg/distributed/remote/generator_config.yml docs/swagger.yml

generate-distributed-mocks:
	mockgen -source=pkg/distributed/models/interfaces.go -destination=pkg/distributed/models/mocks/mocks.gen.go IDealer
