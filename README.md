# HTTPProxy
Самостоятельная работа в рамках 3-го семестра программы по Веб-разработке Образовательного центра VK x МГТУ им.Н.Э.Баумана (ex. "Технопарк") по дисциплине "Безопасность веб-приложений".
<br/>
Ссылка на ДЗ: https://docs.google.com/document/d/1b_ORwryxU-Gx5T1pJrbC1LDkzT-4Vz8JBVUyWRfazlY/edit

## Getting started
### Сгенерируйте ssl сертификаты:
```
./scripts/gen_ca.sh
./scripts/gen_cert.sh
./scripts/install_ca.sh
```
### Запустите приложение в Docker:
```
docker-compose build
docker-compose up
```

### Как это работает
Сервер Proxy запущен на порту 8080. <br/>
Вы можете воспользоваться api в [Swagger](http://localhost:8000/swagger/index.html) для тестирования запросов.<br/>
В качестве БД было выбрано MongoDB

### Что было реализовано:
1) Проксирование HTTP запросов. Проверка:
```
curl -x http://127.0.0.1:8080 http://mail.ru
```
2) Проксирование HTTPS запросов. Проверка:
```
curl -x http://127.0.0.1:8080 https://mail.ru
```
3) Повторная отправка проксированных запросов. Проверка:
    * Выполните команду: ```curl -x http://127.0.0.1:8080 http://mail.ru```
    * Перейдите в [API](http://localhost:8000/swagger/index.html) и выполните Get запрос на /requests
    * Скопируйте из ответа ID запроса и Отправьте его с помощью POST запроса на /repeat/{id}
