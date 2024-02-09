package mover

import "os"

func RunWithAWSProfile[T any](profile string, f func() (T, error)) (T, error) {
	orig := os.Getenv("AWS_PROFILE")
	os.Setenv("AWS_PROFILE", profile)
	defer os.Setenv("AWS_PROFILE", orig)

	return f()
}
