package handler

const (
	StatusSuccess          = "SUCCESS"
	StatusDeleted          = "DELETED"
	StatusVersionNotFound  = "VERSION_NOT_FOUND"
	StatusVersionsNotFound = "VERSIONS_NOT_FOUND"
)

type HTTPStatus struct {
	Status  string `json:"status,omitempty"`
	Error   string `json:"error,omitempty"`
	Version int    `json:"version,omitempty"`
}
