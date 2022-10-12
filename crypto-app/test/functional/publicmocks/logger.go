package publicmocks

var EmptyLogger = emptyLogger{}

type emptyLogger struct{}

func (e emptyLogger) Debug(string) {}

func (e emptyLogger) Info(string) {}

func (e emptyLogger) Error(string) {}
