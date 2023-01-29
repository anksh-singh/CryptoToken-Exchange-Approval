package common

import (
	"bridge-allowance/utils"
)

// TODO: Move services to a common model
type Services struct {
	Http                     *utils.HttpRequest
	Helper                   *utils.Helpers
}
