include Makeroutines.mk

VERSION=$(shell git rev-parse HEAD)
DATE=$(shell date +'%Y-%m-%dT%H:%M%:z')
LDFLAGS=-ldflags '-X github.com/ligato/cn-sample-service/vendor/github.com/ligato/cn-infra/core.BuildVersion=$(VERSION) -X github.com/ligato/cn-sample-service/vendor/github.com/ligato/cn-infra/core.BuildDate=$(DATE)'

# run code analysis
define lint_only
    @echo "# running code analysis"
    @./scripts/golint.sh
    @./scripts/govet.sh
    @echo "# done"
endef

# build helloworld only
define build_helloworld_only
    @echo "# building helloworld"
    @cd cmd/helloworld && go build -v ${LDFLAGS}
    @echo "# done"
endef

# build cassandra only
define build_cassandra_only
    @echo "# building cassandra"
    @cd cmd/cassandra && go build -v ${LDFLAGS}
    @echo "# done"
endef

# clean helloworld only
define clean_helloworld_only
    @echo "# cleaning hello world"
    @rm -f cmd/helloworld/helloworld
    @echo "# done"
endef

# clean cassandra only
define clean_cassandra_only
    @echo "# cleaning cassandra"
    @rm -f cmd/cassandra/cassandra
    @echo "# done"
endef

# build all binaries
build:
	$(call build_helloworld_only)
	$(call build_cassandra_only)

# install dependencies
install-dep:
	$(call install_dependencies)

# update dependencies
update-dep:
	$(call update_dependencies)

# run & print code analysis
lint:
	$(call lint_only)

# clean
clean:
	$(call clean_helloworld_only)
	$(call clean_cassandra_only)
	@echo "# cleanup completed"

# run all targets
all:
	$(call lint_only)
	$(call build_helloworld_only)
	$(call build_cassandra_only)

.PHONY: build update-dep install-dep lint clean
