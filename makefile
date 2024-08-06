clean:
	rm proto/helloworld/v1/*.go
	rm proto/helloworld/v1/*.java
	rm proto/helloworld/v1/*.json
	rm proto/helloworld/v1/*.py
	rm proto/helloworld/v1/*.ts

generate:
	buf generate proto/helloworld/v1/*.proto