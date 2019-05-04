package teaproxy

import (
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/logs"
	"io"
	"net/http"
	"os"
	"regexp"
)

var urlPrefixRegexp = regexp.MustCompile("^(?)(http|https|ftp)://")

func (this *Request) callPage(writer *ResponseWriter, status int) (shouldStop bool) {
	if len(this.pages) == 0 {
		return false
	}
	for _, page := range this.pages {
		if page.Match(status) {
			if urlPrefixRegexp.MatchString(page.URL) {
				err := this.callURL(writer, http.MethodGet, page.URL)
				if err != nil {
					logs.Error(err)
				}
				return true
			} else {
				file := Tea.Root + Tea.DS + page.URL
				fp, err := os.Open(file)
				if err != nil {
					logs.Error(err)
					msg := "404 page not found: '" + page.URL + "'"

					writer.WriteHeader(http.StatusNotFound)
					writer.Write([]byte(msg))
					return true
				}

				// writer.WriteHeader(status)
				// 状态码改成200
				writer.WriteHeader(http.StatusOK)
				_, err = io.Copy(writer, fp)
				if err != nil {
					logs.Error(err)
				}
				fp.Close()
			}

			return true
		}
	}
	return false
}
