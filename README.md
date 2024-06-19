<h2>Тестовое задания на стажировку авито</h2>

<h3>Условие: https://github.com/avito-tech/backend-trainee-assignment-2024  </h3>

<h5>Основная информация:</h5>
<ul>
  <li> Язык сервиса: Go.</li>
  <li> База данных:PostgreSQL.</li>
</ul>

<h5>Как запустить:</h5>
<ul>
  <li> Код запускается с командой `go run cmd\banner-sercive/main.go`</li>
  <li> Указываете переменую окружения `CONFIG_PATH = ".config/local.yaml"`</li>
  <li> В local.yaml указывает параметр подключения к бд - "storage_path" и секретный ключ для token - "signingKey"</li>
</ul>

<h5>принимает:</h5>
<ul>
  <li> /banner - post, get</li>
  <li> /banner/{id} - delete</li>
  <li> /user_banner - get</li>
</ul>




