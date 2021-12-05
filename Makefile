VERSION = v0.12.3
IMAGE = iskorotkov/chaos-scheduler

.PHONY: ci
ci: build test-short test build-image push-image deploy

.PHONY: build
build:
	go build ./...

.PHONY: test
test:
	go test -v ./...

.PHONY: test-short
test-short:
	go test -short -v ./...

.PHONY: build-image
build-image:
	docker build -f build/scheduler.dockerfile -t $(IMAGE):$(VERSION) .

.PHONY: push-image
push-image:
	docker push $(IMAGE):$(VERSION)

.PHONY: deploy
deploy:
	kubectl apply -f deploy/scheduler.yaml

.PHONY: undeploy
undeploy:
	kubectl delete -f deploy/scheduler.yaml
