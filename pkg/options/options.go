package options

import (
	"flag"
	"fmt"
	"github.com/fatih/color"
)

type Options struct {
	Url string
	TargetFile string
	Rate int
}




func ParseOptions()*Options {
	options := &Options{}
	flag.StringVar(&options.Url,"u","","要扫描的目标url,例如:http://www.example.com")
	flag.StringVar(&options.TargetFile,"f","","存放目标的文件,一个目标一行")
	flag.IntVar(&options.Rate,"r",20,"并发数")
	flag.Usage = usage
	flag.Parse()

	return options
}

func usage(){
	fmt.Fprintf(color.Output,color.CyanString(`Go语言指纹识别工具
Options:
`))
	flag.PrintDefaults()
}

