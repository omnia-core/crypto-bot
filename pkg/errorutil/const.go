package errorutil

type OmniaError string

func (e OmniaError) Error() string {
	return string(e)
}

const (
	ErrInvalidMinutes   OmniaError = "invalid minutes"
	ErrInvalidMarket    OmniaError = "invalid market"
	ErrInvalidSide      OmniaError = "invalid side"
	ErrInvalidVolume    OmniaError = "invalid volume"
	ErrInvalidPrice     OmniaError = "invalid price"
	ErrInvalidOrderType OmniaError = "invalid order type"
)
