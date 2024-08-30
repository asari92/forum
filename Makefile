# Инициализация базы данных
# initDB:
# 	# Создание базы данных и выполнение начального скрипта
# 	sqlite3 forum.db < ./docs/new_forum.sql
	
# 	# Включение поддержки внешних ключей и выполнение тестовых данных
# 	sqlite3 forum.db < ./docs/testdata.sql

# Генерация самоподписанных сертификатов
generateCerts:
	# Создание директории для хранения сертификатов
	mkdir -p tls
	
	# Генерация приватного ключа и самоподписанного сертификата
	openssl req -x509 -newkey rsa:4096 -keyout tls/key.pem -out tls/cert.pem -days 365 -nodes -subj "/CN=localhost"

# Команда для полной инициализации проекта
init: generateCerts #initDB 
	@echo "Initialization complete!"
