.PHONY: test clean

build-agent:
	@echo "🚀 Building agent..."
	go build -o cmd/agent/agent cmd/agent/*.go

build-server:
	@echo "🚀 Building server..."
	go build -o cmd/server/server cmd/server/*.go

build: build-agent build-server

test_iter1:
	@echo "🧪 Running tests iter1"
	metricstest -test.v -test.run=^TestIteration1$$ -binary-path=cmd/server/server

test_iter2:
	@echo "🧪 Running tests iter2"
	metricstest -test.v -test.run=^TestIteration2[AB]*$$ -source-path=. -agent-binary-path=cmd/agent/agent

test_iter3:
	@echo "🧪 Running tests iter3"
	metricstest -test.v -test.run=^TestIteration3[AB]*$$ -source-path=. -agent-binary-path=cmd/agent/agent -binary-path=cmd/server/server

test_iter4:
	@echo "🧪 Running tests iter4"
	metricstest -test.v -test.run=^TestIteration4$$ -agent-binary-path=cmd/agent/agent -binary-path=cmd/server/server -server-port=3333 -source-path=.

test_iter5:
	@echo "🧪 Running tests iter5"
	metricstest -test.v -test.run=^TestIteration5$$ -agent-binary-path=cmd/agent/agent -binary-path=cmd/server/server -server-port=3333 -source-path=.

test_iter6:
	@echo "🧪 Running tests iter6"
	metricstest -test.v -test.run=^TestIteration6$$ -agent-binary-path=cmd/agent/agent -binary-path=cmd/server/server -server-port=3343 -source-path=.

test_iter7:
	@echo "🧪 Running tests iter7"
	metricstest -test.v -test.run=^TestIteration7$$ -agent-binary-path=cmd/agent/agent -binary-path=cmd/server/server -server-port=3333 -source-path=.

test_iter8:
	@echo "🧪 Running tests iter8"
	metricstest -test.v -test.run=^TestIteration8$$ -agent-binary-path=cmd/agent/agent -binary-path=cmd/server/server -server-port=3333 -source-path=.

test_iter9:
	@echo "🧪 Running tests iter8"
	metricstest -test.v -test.run=^TestIteration9$$ -agent-binary-path=cmd/agent/agent -binary-path=cmd/server/server -file-storage-path=service.json -server-port=3333 -source-path=.

test_iter10:
	@echo "🧪 Running tests iter8"
	metricstest -test.v -test.run=^TestIteration10[AB]$$ -agent-binary-path=cmd/agent/agent -binary-path=cmd/server/server -database-dsn='postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable' -server-port=3333 -source-path=.

clean:
	rm -rf cmd/agent/agent cmd/server/server

rebuild: clean build

test: build test_iter1 test_iter2 test_iter3 test_iter4 test_iter5 test_iter6 test_iter7 test_iter8 test_iter9 test_iter10

test-clean: test clean