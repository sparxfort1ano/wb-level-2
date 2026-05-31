export PROJECT_ROOT=${shell pwd}


run: 
	@export LOGGER_FOLDER=${PROJECT_ROOT}/out/logs \
	go mod tidy && \
	go run ${PROJECT_ROOT}/cmd/calendar/main.go

logs-cleanup:
	@read -p "Очистить все log файлы окружения? Опасность утери логов. [y/N]: " ans; \
	if [ "$$ans" = "y"  ]; then \
		rm -rf ${PROJECT_ROOT}/out/logs && \
		echo "Файлы логов очищены"; \
	else \
		echo "Очистка логов отменена"; \
	fi