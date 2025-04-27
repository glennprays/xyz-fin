
run-dev:
	@echo "Starting all Docker containers for development mode..."
	@sh docker-dev up -d

stop-dev:
	@echo "Stopping all Docker containers..."
	@sh docker-dev down

swagger:
	@echo "Restarting Swagger UI..."
	@sh docker-dev restart swagger-ui
