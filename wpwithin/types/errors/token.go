package errors

type ErrorType string

const (
	NotFound              ErrorType = "Order not found for ServiceDeliveryToken"
	EmptyKey              ErrorType = "ServiceDeliveryToken key is empty"
	TokenExpired          ErrorType = "ServiceDeliveryToken has expired"
	InvalidKey            ErrorType = "Invalid ServiceDeliveryToken key"
	TooManyUnitsRequested ErrorType = "Too many units requested"
	Success               ErrorType = "Success"
	OrderError            ErrorType = "Order Error"
)
