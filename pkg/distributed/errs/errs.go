package errs

import "github.com/pkg/errors"

// WhileSerializeRequest _
func WhileSerializeRequest(err error) error {
	return errors.Wrap(err, "Serializing request")
}

// WhileDeserializeRequest _
func WhileDeserializeRequest(err error) error {
	return errors.Wrap(err, "Deserializing request")
}

// WhileSerializeResponse _
func WhileSerializeResponse(err error) error {
	return errors.Wrap(err, "Serializing response")
}

// WhileDeserializeResponse _
func WhileDeserializeResponse(err error) error {
	return errors.Wrap(err, "Deserializing response")
}

// WhilePerformRequest _
func WhilePerformRequest(err error) error {
	return errors.Wrap(err, "Performing request")
}

// WhileHandleRequest _
func WhileHandleRequest(err error) error {
	return errors.Wrap(err, "Handling request")
}
