initDB:
	# Создание базы данных и выполнение начального скрипта
	sqlite3 forum.db < ./docs/forum.sql
	
	# Включение поддержки внешних ключей и выполнение тестовых данных
	sqlite3 forum.db < ./docs/testdata.sql
