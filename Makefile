# Инициализация базы данных
initDB:
	# Создание базы данных и выполнение начального скрипта
	sqlite3 forum.db < ./docs/new_forum.sql
	
	

# Генерация самоподписанных сертификатов
generateCerts:
	# Создание директории для хранения сертификатов
	mkdir -p tls
	
	# Генерация приватного ключа и самоподписанного сертификата
	openssl req -x509 -newkey rsa:4096 -keyout tls/key.pem -out tls/cert.pem -days 365 -nodes -subj "/CN=localhost"

# Команда для полной инициализации проекта
init: generateCerts initDB
	@echo "Initialization complete!"

IMAGE_NAME=forum

# Сборка образа
build:
	docker build -t $(IMAGE_NAME) .

# Запуск контейнера
run:
	@docker run --rm -p 4000:4000 $(IMAGE_NAME)
