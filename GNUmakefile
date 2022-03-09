default: testacc

# Run acceptance tests
.PHONY: testacc
testacc:
	TF_ACC=1 DEBUG_HTTP=0 go test ./... -v $(TESTARGS) -timeout 120m
