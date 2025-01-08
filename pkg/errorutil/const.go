package errorutil

type OmniaError string

func (e OmniaError) Error() string {
	return string(e)
}

const (
	ErrInvalidMinutes OmniaError = "invalid minutes"
	ErrInvalidMarket  OmniaError = "invalid market"
)
