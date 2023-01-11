projectName ?= ai-creator

.PHONY: help
help:
	@echo "help:"
	@echo "use [ make docker-build ] docker build image"
	@echo "use [ make docker-run ] docker run builded image with env"

.PHONY: docker-build
docker-build:
	@docker build -t weedge/$(projectName):latest -f build/Dockerfile .

# u maybe see this: https://www.tutorialworks.com/container-networking/
.PHONY: docker-run
docker-run:
#	@make docker-build	
	@docker run -it \
		--name doraemon-ai-creator \
		--env NOLIBOX_API_SK=$(NOLIBOX_API_SK) \
		--env NOLIBOX_API_AK=$(NOLIBOX_API_AK) \
		--env REDIS_TYPE=$(REDIS_TYPE) \
		--env REDIS_URL=$(REDIS_URL) \
		--env REDIS_PASSWORD=$(REDIS_PASSWORD) \
		--env REDIS_USERNAME=$(REDIS_USERNAME) \
		--env TEMPLATE_DIR=$(TEMPLATE_DIR) \
		-p 8123:8123  \
		weedge/$(projectName)

