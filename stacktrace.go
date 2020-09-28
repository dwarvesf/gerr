package gerr

type stacktrace struct {
	file      string
	line      int
	function  string
	fullTrace string
}

func newStackTrace(file string, line int, fnName string, fullTrace string) stacktrace {
	return stacktrace{
		file:      file,
		line:      line,
		function:  fnName,
		fullTrace: fullTrace,
	}
}
