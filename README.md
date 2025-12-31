# Тестовое задания на стажировку

## Запуск приложения

Находясь в директории в терминале прописать:

### Без docker:

```
go run cmd/app/main.go
```

### Через docker:

```
docker build -t image_name:image_tag .
docker run -p 8080:8080 image_name:image_tag
```

(для теста image_name и image_tag можно указать любые)

#### для завершения работы нажмите ctrl+c

API будет доступно по адресу http://localhost:8080/todos и http://localhost:8080/todos/{id} соответственно.
