package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

func InitLogger(verbose bool) {
    encoderConfig := zapcore.EncoderConfig{
        MessageKey: "message",                        
        LevelKey:   "level",                           
        EncodeLevel: zapcore.CapitalLevelEncoder,     
    }

	encoder := zapcore.NewConsoleEncoder(encoderConfig) // Use ConsoleEncoder or JSONEncoder
	
	zapLevel := zap.InfoLevel
	if verbose {
		zapLevel = zap.DebugLevel	
	}
	core := zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapLevel)

	logger := zap.New(core)
	sugaredLogger := logger.Sugar()

	defer logger.Sync() 

	logger.Info("Zap logger initialized",
		zap.String("module", "main"),
		zap.Int("version", 1),
	)
		
	zap.ReplaceGlobals(sugaredLogger.Desugar())
}	
