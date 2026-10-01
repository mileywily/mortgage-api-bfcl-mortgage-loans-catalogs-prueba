package config

// Headers salientes que esta API agrega en el wire-up del cliente REST.
// Los headers corporativos del gateway los propaga go-lib-http-fif
// (WithMandatoryHeaders) y no se declaran acá.
const (
	HeaderApplicationID = "X-Application-ID"
	HeaderContentType   = "Content-Type"
)
