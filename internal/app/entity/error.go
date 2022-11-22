package entity

const (
	ErrInvalidInput            AppError = "invalid input"
	ErrResourceIsEmpty         AppError = "resource is empty"
	ErrResourceNotFound        AppError = "resource not found"
	ErrResourceHasExisted      AppError = "resource has existed"
	ErrPermissionDenied        AppError = "permission denied"
	ErrRuntimePanic            AppError = "runtime panic"
	ErrInternal                AppError = "internal"
	ErrUnauthorized            AppError = "unauthorized"
	ErrTransactionNotCompleted AppError = "transation not completed"
)

type AppError string

func (e AppError) Error() string {
	return string(e)
}
