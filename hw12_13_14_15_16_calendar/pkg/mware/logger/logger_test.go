package logger

import (
	"testing"
)

func TestLogger(_ *testing.T) {
	// ctx := context.Background()
	// config, err := configs.New(&ctx, true, "./../configs/cfg.json") // ./configs/cfg.json
	// require.NoError(t, err)

	// encoderCfg := zap.NewProductionEncoderConfig()
	// logLevel := zap.InfoLevel

	// if config.Logger.Level != "" {
	// 	logLevel, err = zapcore.ParseLevel(config.Logger.Level)
	// 	require.NoError(t, err)
	// }

	// logg := zap.New(zapcore.NewCore(
	// 	zapcore.NewJSONEncoder(encoderCfg),
	// 	zapcore.Lock(os.Stdout),
	// 	zap.NewAtomicLevelAt(logLevel),
	// ))

	// f := WithLogging(func(w http.ResponseWriter, r *http.Request) {
	// 	w.WriteHeader(200)
	// }, logg)

	// writer := httptest.NewRecorder()

	// f(writer, httptest.NewRequest(http.MethodGet, "/", strings.NewReader("")))
	// require.Equal(t, t, writer.Code, 200)
}
