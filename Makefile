generate-distributed-client:
	oapi-codegen --config pkg/distributed/remote/generator_config.yml docs/swagger.yml
