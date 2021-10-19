package converters

// ToRPCError converts internal error to RPC format
func ToRPCError(err error) error {
	return err
}

// FromRPCError converts RPC error to internal format
func FromRPCError(err error) error {
	return err
}
