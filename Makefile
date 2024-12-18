# Инициализация базы данных
.PHONY: initDB
initDB:
	# Создание базы данных и выполнение начального скрипта
	sqlite3 forum.db < ./schema/forum.sql
	
# Генерация самоподписанных сертификатов
.PHONY: generateCerts
generateCerts:
	# Создание директории для хранения сертификатов
	mkdir -p tls
	
	# Генерация приватного ключа и самоподписанного сертификата
	openssl req -x509 -newkey rsa:4096 -keyout tls/key.pem -out tls/cert.pem -days 365 -nodes -subj "/CN=localhost"

# Команда для полной инициализации проекта
.PHONY: init
init: generateCerts initDB
	@echo "Initialization complete!"

IMAGE_NAME = forum
CONTAINER_NAME = forum-container
PORT = 4000

# Сборка Docker-образа
.PHONY: build
build: init
	@echo "Сборка Docker-образа $(IMAGE_NAME)..."
	@docker build -t $(IMAGE_NAME) .
	@echo "Образ $(IMAGE_NAME) успешно собран!"

# Запуск Docker-контейнера
.PHONY: run
run:
	@echo "Запуск Docker-контейнера $(CONTAINER_NAME)..."
	@docker run --rm --name $(CONTAINER_NAME) -p $(PORT):$(PORT) $(IMAGE_NAME)
	@echo "Контейнер $(CONTAINER_NAME) запущен на порту $(PORT)!"

# Остановка Docker-контейнера (если запущен в фоновом режиме)
.PHONY: stop
stop:
	@echo "Остановка Docker-контейнера $(CONTAINER_NAME)..."
	@docker stop $(CONTAINER_NAME) || echo "Контейнер $(CONTAINER_NAME) не запущен."
