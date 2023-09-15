package constants

type otel struct {
	HttpErrorMessage string
}

var Otel = otel{
	HttpErrorMessage: "http.error_message",
}
