Я не смог настроить Базу данных в kubernetes, по скольку необходим ssl сертификат, который не был предоставлен.  
По запросу https://cnrprod1725729288-team-77382-32918.avito2024.codenrock.com/api/tenders/new  
Ошибка (rc1b-5xmqy6bq501kls4m.mdb.yandexcloud.net): server error: ERROR: odyssey: c9d3568799253: SSL is required (SQLSTATE 08P01)  

По запросу localhost:8080/api/tenders/new 
Все работает идеально.
### Запуск:
    docker compose up -d
    make all

### Postman Workspace:
    https://www.postman.com/rryowa/avito/collection/zj01s72/tender-management-api

### Описание конфигурации линтера
    zadanie-6105/.golangci.yml

### Не реализовано:
    Расширенный процесс согласования.