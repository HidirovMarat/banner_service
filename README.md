Язык сервиса: Go. 
База данных:PostgreSQL.
---
Код запускается с командой 
go run cmd\banner-sercive/main.go
--
указываете переменую окружения 
CONFIG_PATH = ".config/local.yaml"
----
В local.yaml указывает параметр подключения к бд - "storage_path" и секретный ключ для token - "signingKey"
--
принимает:
/banner - post, get
/banner/{id} - delete
/user_banner - get
--



