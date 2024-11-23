# Docker_Lab

## Описание приложения

Сервер на Go, реализующий API для создания тендеров и предложений для них. В качестве СУБД используется Postgres

## Запуск

`docker compose up`

##  Что сделано по лабораторной

- Написан Dockerfile для создания легковесного образа сервера
  - Используется мультстейдж сборка, полученный на первом этапе бинарник помещается в чистый alpine

- Написан docker-compose файл, собирающий приложение из двух контейнеров (сервер и БД)
  - В качестве образа для БД использовался также легковесный, основанный на alpine
  - Для сохранения данных БД между перезапусками контейнера используется volume
  - Сервер поднимается после БД (использовался `depends on`)

- Вся конфигурация производилась через переменные окружения, которые сначала прокидывались в compose, а оттуда по контейнерам

- Реализован механизм автомиграций при запуске контейнера с сервером через golang-migrate

- Сервер подключается к БД по внутренней сети докера, поэтому в переменной POSTGRES_CONN вместо localhost'а стоит psql (при желании, можно закрыть доступ к базе извне, если не прокидывать порты наружу)

- Контейнеры перезапускаются, пока не остановишь их вручную (настройка `restart`)

## Работоспособность

Чтобы убедиться, что все работает можно:
- После запуска дернуть к примеру `http://localhost:8080/api/ping`, если запускать с .env из репозитория
- Подключиться к БД и глянуть, что миграции накатились

## Авторы
- Шамхалов Якуб гр. 2300
- Жохов Кирилл гр. 2300

