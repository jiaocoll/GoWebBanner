package scan

import (
	"GoWebBanner/pkg/cmsparse"
	"GoWebBanner/pkg/request"
	"bufio"
	"crypto/md5"
	"fmt"
	"github.com/fatih/color"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
	"time"
)

type Args struct {
	Url string
	Path string
	Cmsdata []cmsparse.Cms
}

func Wafscan(args Args)bool{
	url := args.Url
	var wafname string
	data,headers,err := request.Get(url)
	if err != nil {
		return false
	}

	waffile, err1 := os.Open("waf.ini")
	if err1 != nil {
		log.Println(err1)
	}
	defer waffile.Close()
	tmps := bufio.NewScanner(waffile)
	for tmps.Scan(){
		tmp := tmps.Text()
		datas := strings.Split(tmp,"|")
		option := datas[1]
		content := datas[3]
		if option == "index"{
			if strings.Contains(string(data),content){
				wafname = datas[0]
				break
			}
		}else {
			k,ok := headers[datas[2]]
			if ok {
				match, _ := regexp.MatchString(content,strings.Join(k,""))
				if match {
					wafname = datas[0]
					break
				}
			}
		}

		if wafname != "" {
			fmt.Fprintln(color.Output,color.HiRedString("[WARRING] ")+color.HiYellowString("URL: ")+url+wafname)
			return true
		}
	}
	return false
}


func Bannerscan(args Args) bool{

	url := args.Url
	path := args.Path
	Cmsdata := args.Cmsdata
	target := url + path
	resp1, err1 := request.Head(target)
	if err1 != nil{
		return false
	}
	resp1bytes, _ := ioutil.ReadAll(resp1.Body)
	defer resp1.Body.Close()

	resp2,_ ,err2 := request.Get(target)
	if err2 != nil{
		time.Sleep(2 * time.Second)
		return false
	}

	for _,cmsinfo := range Cmsdata{
		if strings.Contains(string(resp2),"拦截") || strings.Contains(string(resp1bytes),"拦截"){
			fmt.Fprintln(color.Output,color.HiRedString("[WARING] ")+color.HiYellowString("WAF! ")+target)
			time.Sleep(1 * time.Second)
			return false
		}
		if cmsinfo.Option == "keyword"{
			if strings.Contains(string(resp2),cmsinfo.Content){
				fmt.Fprintln(color.Output,color.HiCyanString("Url:")+url+color.HiCyanString(" Cms:")+cmsinfo.Name+color.HiCyanString(" Path:")+cmsinfo.Path+color.HiCyanString(" Option:")+cmsinfo.Option+color.HiCyanString(" Content:")+cmsinfo.Content)
				return true
			}
		}else if cmsinfo.Option == "md5" {
			md5str := fmt.Sprintf("%x",md5.Sum(resp2))
			if md5str == cmsinfo.Content{
				fmt.Fprintln(color.Output,color.HiCyanString("Url:")+url+color.HiCyanString(" Cms:")+cmsinfo.Name+color.HiCyanString(" Path:")+cmsinfo.Path+color.HiCyanString(" Option:")+cmsinfo.Option+color.HiCyanString(" Content:")+cmsinfo.Content)
				return true
			}
		}
	}
	return false
}