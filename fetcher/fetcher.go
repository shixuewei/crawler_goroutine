package fetcher

import (
	"errors"
	"io/ioutil"
	"net/http"
	"time"
)

// func determineEncoding(r *bufio.Reader) encoding.Encoding {
// 	bytes, err := r.Peek(1024)
// 	if err != nil {
// 		log.Printf("Fetcher err : %v", err)
// 		return unicode.UTF8
// 	}
// 	e, _, _ := charset.DetermineEncoding(bytes, "")
// 	return e
// }

//降速（将一个time.Time数据传入管道中，时间到了以后就会卡住）
//限制下载速度
var rateLimit = time.Tick(5 * time.Millisecond)

func Fetch(url string) ([]byte, error) {
	//弹出以后就会接着执行，不会卡死
	<-rateLimit
	//跳过证书验证
	// tr := &http.Transport{
	// 	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	// }
	// //http cookie接口
	// cookieJar, _ := cookiejar.New(nil)
	// c := &http.Client{
	// 	Jar:       cookieJar,
	// 	Transport: tr,
	// }
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}
	//bodyReader := bufio.NewReader(resp.Body)
	//将gbk编码转变为utf-8
	//e := determineEncoding(bodyReader)
	//utf8Reader := transform.NewReader(resp.Body, e.NewDecoder())

	return ioutil.ReadAll(resp.Body)

}
