package httpclient

import "errors"

func wrapErr(err error, message string) error {
	if err == nil {
		return nil
	}
	return errors.New(message + ": " + err.Error())
}
