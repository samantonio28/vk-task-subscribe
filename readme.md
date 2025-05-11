# Сервис подписки Publisher-Subscriber

Сервис предназначен для подписки на события. Работает по gRPC.

## Зависимости

- gRPC
- Logrus
- Cleanenv

## Установка

1. Склонируйте репозиторий
2. Установите зависимости `go mod tidy`
3. Запустите сервис `go run cmd/main.go`

## API

Сервис предоставляет следующие методы:

- `Subscribe` - подписка на события
- `Publish` - публикация события

## Конфигурация

Сервис настроен через файл `config.yaml`.

При написании сервиса я придерживался метода Dependency Injection, когда вызывается 
конструктор с нужными зависимостями:

```go
func NewSubPubService(logger *logger.LogrusLogger, subPub repository.SubPub) (*SubPubService, error) {
    if logger == nil || subPub == nil {
        return nil, status.Error(codes.InvalidArgument, "logger and subPub must not be nil")
    }

    return &SubPubService{
        subPub: subPub,
        Logger: logger,
        subs:   make(map[string]map[subpub.PubSub_SubscribeServer]struct{}), 
    }, nil
}
```

При запуске сервиса происходит инициализация зависимостей:

```go
subPubService, err := usecase.NewSubPubService(logger, subpub)
if err != nil {
    fmt.Printf("Failed to initialize subPubService: %v\n", err)
    return
}
```
