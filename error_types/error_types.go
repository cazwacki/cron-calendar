package errortypes

import "errors"

type DatabaseError struct {
	error
}

type ValidationError struct {
	error
}

type AuthorizationError struct {
	error
}

func GenerateValidationError(message string) ValidationError {
	return ValidationError{
		error: errors.New(message),
	}
}

func GenerateDatabaseError(message string) DatabaseError {
	return DatabaseError{
		error: errors.New(message),
	}
}

func GenerateAuthorizationError(message string) AuthorizationError {
	return AuthorizationError{
		error: errors.New(message),
	}
}
