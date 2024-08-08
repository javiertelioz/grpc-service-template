clean:
	@rm proto/helloworld/v1/*.{go,java,py,ts,yaml}
	@rm proto/payments/v1/*.{go,java,py,ts,yaml}
	@echo "Success clean folder"

generate:
	@buf generate proto/helloworld/v1/*.proto
	@buf generate proto/payments/v1/*.proto
	@echo "Success generate"