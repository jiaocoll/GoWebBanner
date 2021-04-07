package main

import (
	"GoWebBanner/pkg/cmsparse"
	"GoWebBanner/pkg/options"
	runner2 "GoWebBanner/pkg/runner"
	"GoWebBanner/pkg/scan"
	"bufio"
	"fmt"
	"github.com/fatih/color"
	"github.com/panjf2000/ants/v2"
	"log"
	"os"
	"strconv"
	"sync"
	"time"
)

var(
	targets[] string
)


func main(){
	start := time.Now()
	var count1 int
	var count2 int
	options := options.ParseOptions()
	runner := runner2.NewRunner(options)
	sortCoups, webdata := cmsparse.Parsecms("cms.json")
	if runner.Options.TargetFile == ""{
		count1 = 1
		for _, v := range sortCoups{
			args := scan.Args{
				Url: runner.Options.Url,
				Path: v.Path,
				Cmsdata: webdata[v.Path],
			}
			if scan.Wafscan(args){
				time.Sleep(1 * time.Second)
			}
			if scan.Bannerscan(args){
				count2++
				break
			}
		}

	}else {
		targetfile, err := os.OpenFile(runner.Options.TargetFile,os.O_RDONLY,1)
		defer targetfile.Close()
		if err != nil{
			log.Println(err)
		}

		tmps := bufio.NewScanner(targetfile)
		for tmps.Scan(){
			tmp := tmps.Text()
			targets = append(targets, tmp)
		}

		var wg sync.WaitGroup
		p,_ := ants.NewPoolWithFunc(runner.Options.Rate, func(i interface{}) {
			if scan.Wafscan(i.(scan.Args)){
				time.Sleep(1 * time.Second)
			}
			if scan.Bannerscan(i.(scan.Args)){
				count2++
			}
			wg.Done()
		})
		for _,url := range targets{
			count1++
			for _,v := range sortCoups{
				args := scan.Args{
					Url: url,
					Path: v.Path,
					Cmsdata: webdata[v.Path],
				}
				wg.Add(1)
				_ = p.Invoke(args)
			}
		}
	}
	end := time.Since(start)
	fmt.Fprintln(color.Output,color.HiCyanString("探测目标数:")+strconv.Itoa(count1)+color.HiMagentaString(" 识别成功数:")+strconv.Itoa(count2)+color.HiGreenString(" 总用时:")+end.String())

}
