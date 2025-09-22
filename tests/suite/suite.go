package suite

import (
	"context"
	"github.com/evgeniySeleznev/gRPC-auth/internal/config"
	auth_v1 "github.com/evgeniySeleznev/gRPC-auth/pkg/authV1"
	"os"
	"testing"
)

type Suite struct {
	*testing.T                    //потребуется для вызова методов *testing.T внутри Suite
	cfg        *config.Config     //конфигурация приложения
	AuthClient auth_v1.AuthClient //Клиент для взаимодействия с gRPC-сервером
}

func New(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	t.Parallel() // тестим параллельные обращения к сервису

	//key := "CONFIG_PATH"
	//if v := os.Getenv(key); v != "" {
	//	return v
	//} заготовка для получения конфига из переменной окружения
	// не запускаем на гитхабе, поэтому норм

	cfg := config.MustLoadByPath("../config/local_tests.yaml")

	ctx, cancelCtx := context.WithTimeout(context.Background(), cfg.GRPC.Timeout)

	t.Cleanup(func() {
		t.Helper()
		cancelCtx()
	})

	//2:43:00
}
