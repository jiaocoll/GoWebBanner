package runner

import (
	"GoWebBanner/pkg/cmsparse"
	"GoWebBanner/pkg/options"
)

type Runner struct {
	Options *options.Options
	Path string
	Cmsbanner []cmsparse.Cms
}


func NewRunner(options *options.Options)(*Runner){
	runner := &Runner{
		Options: options,
	}
	return runner
}
