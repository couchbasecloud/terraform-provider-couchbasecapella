default: testacc

# Run acceptance tests
.PHONY: testacc
testacc:
	TF_ACC=1 go test ./provider/... -v $(TESTARGS) -timeout 120m
