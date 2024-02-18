
GOCMD = go
GOBUILD = $(GOCMD) build
GOTEST = $(GOCMD) test

VERSION_FILE = VERSION
BUILD_FILE = BUILD

RELEASE_DIR = .release

APP_NAME = quickrest

RELEASE_DIR_WINDOWS = $(RELEASE_DIR)/windows_amd64/
RELEASE_DIR_MACOS = $(RELEASE_DIR)/macos_amd64/
RELEASE_DIR_LINUX = $(RELEASE_DIR)/linux_amd64/
RELEASE_DIR_FREEBSD = $(RELEASE_DIR)/freebsd_amd64/
RELEASE_DIR_OPENBSD = $(RELEASE_DIR)/openbsd_amd64/
RELEASE_DIR_SOLARIS = $(RELEASE_DIR)/solaris_amd64/

WINDOWS_ARCH = amd64 386
MACOS_ARCH = amd64 arm64
LINUX_ARCH = amd64 386
FREEBSD_ARCH = amd64 386
OPENBSD_ARCH = amd64 386
SOLARIS_ARCH = amd64


BUILD_ARGS = -ldflags "-X main.Version=$(shell cat $(VERSION_FILE)) -X main.Build=$(shell cat $(BUILD_FILE))"


release: bump-build clean windows macos linux freebsd openbsd solaris

windows:
	for arch in $(WINDOWS_ARCH); do \
		version=$$(cat $(VERSION_FILE)).$$(cat $(BUILD_FILE)); \
		release_name=${APP_NAME}_$${version}; \
		target_name=$${release_name}_windows_$${arch}; \
		build_dir=${RELEASE_DIR}/$${release_name}/$${target_name}; \
		mkdir -p $$build_dir; \
		GOOS=windows GOARCH=$${arch} $(GOBUILD) -o $$build_dir/ $(BUILD_ARGS) -v ./... ; \
		tar -czf ${RELEASE_DIR}/$${release_name}/$${target_name}.tar.gz -C ${RELEASE_DIR}/$${release_name} $${target_name}; \
	done

macos:
	for arch in $(MACOS_ARCH); do \
		version=$$(cat $(VERSION_FILE)).$$(cat $(BUILD_FILE)); \
		release_name=${APP_NAME}_$${version}; \
		target_name=$${release_name}_macos_$${arch}; \
		build_dir=${RELEASE_DIR}/$${release_name}/$${target_name}; \
		mkdir -p $$build_dir; \
		GOOS=darwin GOARCH=$${arch} $(GOBUILD) -o $$build_dir/ $(BUILD_ARGS) -v ./... ; \
		tar -czf ${RELEASE_DIR}/$${release_name}/$${target_name}.tar.gz -C ${RELEASE_DIR}/$${release_name} $${target_name}; \
	done

linux:
	for arch in $(LINUX_ARCH); do \
		version=$$(cat $(VERSION_FILE)).$$(cat $(BUILD_FILE)); \
		release_name=${APP_NAME}_$${version}; \
		target_name=$${release_name}_linux_$${arch}; \
		build_dir=${RELEASE_DIR}/$${release_name}/$${target_name}; \
		mkdir -p $$build_dir; \
		GOOS=linux GOARCH=$${arch} $(GOBUILD) -o $$build_dir/ $(BUILD_ARGS) -v ./... ; \
		tar -czf ${RELEASE_DIR}/$${release_name}/$${target_name}.tar.gz -C ${RELEASE_DIR}/$${release_name} $${target_name}; \
	done

freebsd:
	for arch in $(FREEBSD_ARCH); do \
		version=$$(cat $(VERSION_FILE)).$$(cat $(BUILD_FILE)); \
		release_name=${APP_NAME}_$${version}; \
		target_name=$${release_name}_freebsd_$${arch}; \
		build_dir=${RELEASE_DIR}/$${release_name}/$${target_name}; \
		mkdir -p $$build_dir; \
		GOOS=freebsd GOARCH=$${arch} $(GOBUILD) -o $$build_dir/ $(BUILD_ARGS) -v ./... ; \
		tar -czf ${RELEASE_DIR}/$${release_name}/$${target_name}.tar.gz -C ${RELEASE_DIR}/$${release_name} $${target_name}; \
	done

openbsd:
	for arch in $(OPENBSD_ARCH); do \
		version=$$(cat $(VERSION_FILE)).$$(cat $(BUILD_FILE)); \
		release_name=${APP_NAME}_$${version}; \
		target_name=$${release_name}_openbsd_$${arch}; \
		build_dir=${RELEASE_DIR}/$${release_name}/$${target_name}; \
		mkdir -p $$build_dir; \
		GOOS=openbsd GOARCH=$${arch} $(GOBUILD) -o $$build_dir/ $(BUILD_ARGS) -v ./... ; \
		tar -czf ${RELEASE_DIR}/$${release_name}/$${target_name}.tar.gz -C ${RELEASE_DIR}/$${release_name} $${target_name}; \
	done

solaris:
	for arch in $(SOLARIS_ARCH); do \
		version=$$(cat $(VERSION_FILE)).$$(cat $(BUILD_FILE)); \
		release_name=${APP_NAME}_$${version}; \
		target_name=$${release_name}_solaris_$${arch}; \
		build_dir=${RELEASE_DIR}/$${release_name}/$${target_name}; \
		mkdir -p $$build_dir; \
		GOOS=solaris GOARCH=$${arch} $(GOBUILD) -o $$build_dir/ $(BUILD_ARGS) -v ./... ; \
		tar -czf ${RELEASE_DIR}/$${release_name}/$${target_name}.tar.gz -C ${RELEASE_DIR}/$${release_name} $${target_name}; \
	done

clean:
	rm -rf $(RELEASE_DIR)

test:
	$(GOTEST) -v ./...

bump-build:
	@current_build=$$(cat $(BUILD_FILE)); \
    next_build=$$(echo $$((current_build + 1))); \
    echo $${next_build} > $(BUILD_FILE)

.PHONY: all windows macos linux freebsd openbsd solaris clean test bump-build
