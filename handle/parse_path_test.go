package handle

import (
	"imgGo-Server/global"
	"strings"
	"testing"
)

func TestParsePath(t *testing.T) {
	//conf.Host = "https://imgo.erps.bio:88/"
	//p := "https://imgo.erps.bio:88/2022/09/26/1642834046009692757.png"
	conf := global.Conf
	conf.Host = "http://localhost:8800/"
	p := "http://localhost:8800/2023/08/21_16-02-33.388.png.webp"

	p = strings.Replace(p, conf.Host, "", -1)
	_list := strings.Split(p, "/")

	t.Log(strings.Join(_list[:len(_list)-1], "/"))
	t.Log(_list[len(_list)-1])
}
