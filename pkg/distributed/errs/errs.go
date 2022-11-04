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

// WhileGetOutputStorageClaim _
func WhileGetOutputStorageClaim(err error) error {
	return errors.Wrap(err, "Getting output storage claim")
}

// WhileAllocateInputStorageClaim _
func WhileAllocateInputStorageClaim(err error) error {
	return errors.Wrap(err, "Allocating input storage claim")
}

// WhileGetAllInputStorageClaims _
func WhileGetAllInputStorageClaims(err error) error {
	return errors.Wrap(err, "Getting input storage claims")
}

// WhileAllocateOutputStorageClaim _
func WhileAllocateOutputStorageClaim(err error) error {
	return errors.Wrap(err, "Allocating output storage claim")
}

// WhileBuildStorageClaimByURL _
func WhileBuildStorageClaimByURL(err error) error {
	return errors.Wrap(err, "Building storage claim by URL")
}

// WhileCastingProgressStep _
func WhileCastingProgressStep(err error) error {
	return errors.Wrap(err, "Casting progress step")
}
