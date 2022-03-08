default: testacc

# Run acceptance tests
.PHONY: testacc
testacc:
	TF_ACC=true go test ./... -v $(TESTARGS) -timeout 120m
