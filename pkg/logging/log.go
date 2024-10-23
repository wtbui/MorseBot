package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

func InitLogger(verbose bool) {
	// Create a new production logger
    encoderConfig := zapcore.EncoderConfig{
        MessageKey: "message",                         // Key for the actual log message
        LevelKey:   "level",                           // Key for the log level (INFO, ERROR, etc.)
        EncodeLevel: zapcore.CapitalLevelEncoder,      // Encode log level as all caps (e.g., INFO, ERROR)
  // Optional: log caller information in short form
    }

	// Create a new core with the custom encoder
	encoder := zapcore.NewConsoleEncoder(encoderConfig) // Use ConsoleEncoder or JSONEncoder
	
	zapLevel := zap.InfoLevel
	if verbose {
		zapLevel = zap.DebugLevel	
	}
	core := zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapLevel)

	logger := zap.New(core)
	sugaredLogger := logger.Sugar()

	defer logger.Sync() // Ensures any buffered log entries are written

	// Using the logger
	logger.Info("Zap logger initialized",
			zap.String("module", "main"),
			zap.Int("version", 1),
	)
		
	zap.ReplaceGlobals(sugaredLogger.Desugar())
}	
