package xpath

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/x-armory/go-unmarshal/base"
	"strings"
	"testing"
	"time"
)

type MyType struct {
	Source              string     `xm:"xpath://*[@id='content_right']/div/div[2]/a[1]"`
	Id                  int        `xm:"xpath://*[@id='content_right']/div/table/tbody/tr[r[1:]]/td/span"`
	Title               string     `xm:"xpath://*[@id='content_right']/div/table/tbody/tr[r[1:]]/td/a"`
	Year                int        `xm:"xpath://*[@id='form']/input[@name='datetime']/@value pattern='^\\d+'"`
	Month               int        `xm:"xpath://*[@id='form']/input[@name='datetime']/@value pattern='^\\d+-(\\d+)'"`
	Day                 int        `xm:"xpath://*[@id='form']/input[@name='datetime']/@value pattern='^\\d+-\\d+-(\\d+)'"`
	Datetime            time.Time  `xm:"xpath://*[@id='form']/input[@name='datetime']/@value format='2006-01-02 15:04:05' timezone='UTC'"`
	Date                *time.Time `xm:"xpath://*[@id='form']/input[@name='date']/@value pattern='\\d+-\\d+-\\d+' format='2006-01-02' timezone='Asia/Shanghai'"`
	DefaultValueTime    time.Time  `xm:"xpath://*[@id='form']/input[@name='notExist']/@value pattern='\\d+-\\d+-\\d+' format='2006-01-02' timezone='Asia/Shanghai'"`
	DefaultValueTimePtr *time.Time `xm:"xpath://*[@id='form']/input[@name='notExist']/@value pattern='\\d+-\\d+-\\d+' format='2006-01-02' timezone='Asia/Shanghai'"`
	DefaultValueString  string     `xm:"xpath://*[@id='form']/input[@name='notExist']/@value"`
	DefaultValueInt     int        `xm:"xpath://*[@id='form']/input[@name='notExist']/@value"`
	DefaultValueFloat   float64    `xm:"xpath://*[@id='form']/input[@name='notExist']/@value"`
}

func TestUnmarshaler_Unmarshal(t *testing.T) {
	xpathUnmarshaler := Unmarshaler{
		DataLoader: base.DataLoader{
			ItemFilters: []base.ItemFilter{
				func(item interface{}, vars *base.Vars) (flow base.FlowControl, deep int) {
					// check item type
					data, ok := item.(*MyType)
					if !ok {
						return base.Forward, 0
					}
					// validate item
					if data.Id == 0 && data.Title == "" {
						return base.Break, 0
					}
					// process item
					fmt.Printf("%+v\n", data)
					return base.Forward, 0
				},
			},
		},
	}

	var data []*MyType
	err := xpathUnmarshaler.Unmarshal(strings.NewReader(HTML), &data)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(data))
}

var HTML = `
<html><head>
    <meta http-equiv="X-UA-Compatible" content="IE=edge,chrome=1">
    <meta http-equiv="content-type" content="text/html;charset=utf-8">
    <meta content="always" name="referrer">
    <script src="https://ss1.bdstatic.com/5eN1bjq8AAUYm2zgoY3K/r/www/nocache/imgdata/seErrorRec.js"></script>
    <title>页面不存在_百度搜索</title>
</head>

<body link="#0000cc">
    <div id="wrapper" class="wrapper_l">
        <div id="head">
            <div class="head_wrapper">
                <div class="s_form">
                    <div class="s_form_wrapper">
                        <a href="//www.baidu.com/" id="result_logo"><img src="//www.baidu.com/img/baidu_jgylogo3.gif" alt="到百度首页" title="到百度首页"></a>
                        <form id="form" name="f" action="//www.baidu.com/s" class="fm">
                            <input type="hidden" name="ie" value="utf-8">
                            <input type="hidden" name="f" value="8">
                            <input type="hidden" name="rsv_bp" value="1">
                            <input type="hidden" name="ch" value="">
                            <input type="hidden" name="tn" value="baiduerr">
                            <input type="hidden" name="bar" value="">
                            <input type="hidden" name="date" value="其他字符2019-07-03其他字符">
                            <input type="hidden" name="datetime" value="2019-07-02 12:01:02">
                            <span class="bg s_ipt_wr">
                                <input id="kw" name="wd" class="s_ipt" value="" maxlength="255" autocomplete="off" autofocus="">
                            </span><span class="bg s_btn_wr">
                                <input type="submit" id="su" value="百度一下" class="bg s_btn">
                            </span>
                    </form>
                </div>
            </div>
        </div>
    </div>
    <div class="s_tab" id="s_tab">
        <b>网页</b>
        <a href="http://tieba.baidu.com/f?kw=&amp;fr=wwwt" wdfield="kw">贴吧</a>
        <a href="http://zhidao.baidu.com/q?ct=17&amp;pn=0&amp;tn=ikaslist&amp;rn=10&amp;word=&amp;fr=wwwt" wdfield="word">知道</a>
        <a href="http://music.baidu.com/search?fr=ps&amp;ie=utf-8&amp;key=" wdfield="key">音乐</a>
        <a href="http://image.baidu.com/i?tn=baiduimage&amp;ps=1&amp;ct=201326592&amp;lm=-1&amp;cl=2&amp;nc=1&amp;ie=utf-8&amp;word=" wdfield="word">图片</a>
        <a href="http://v.baidu.com/v?ct=301989888&amp;rn=20&amp;pn=0&amp;db=0&amp;s=25&amp;ie=utf-8&amp;word=" wdfield="word">视频</a>
        <a href="http://map.baidu.com/m?word=&amp;fr=ps01000" wdfield="word">地图</a>
        <a href="http://wenku.baidu.com/search?word=&amp;lm=0&amp;od=0&amp;ie=utf-8" wdfield="word">文库</a>
        <a href="//www.baidu.com/more/">更多»</a>
    </div>
    <div id="wrapper_wrapper">
        <div id="content_right"><div class="cr-content"><div class="opr-toplist-title">今日搜索热点</div><table class="opr-toplist-table"><thead><tr><th>排名</th></tr></thead><tbody><tr><td><span class="opr-index-hot1 opr-index-item">1</span><a target="_blank" href="//www.baidu.com/s?word=%E9%9B%B7%E4%BD%B3%E9%9F%B3%E5%A6%BB%E5%AD%90%E5%9B%9E%E5%BA%94&amp;sa=re_dl_seError_1&amp;tn=SE_fengyunbangS_fs9rizg7&amp;rsv_dl=fyb_n_erro" class="opr-item-text">雷佳音妻子回应</a></td></tr><tr><td><span class="opr-index-hot2 opr-index-item">2</span><a target="_blank" href="//www.baidu.com/s?word=%E6%B3%95%E6%A4%8D%E7%89%A9%E4%BA%BA%E5%AE%89%E4%B9%90%E6%AD%BB%E6%A1%88&amp;sa=re_dl_seError_2&amp;tn=SE_fengyunbangS_fs9rizg7&amp;rsv_dl=fyb_n_erro" class="opr-item-text">法植物人安乐死案</a></td></tr><tr><td><span class="opr-index-hot3 opr-index-item">3</span><a target="_blank" href="//www.baidu.com/s?word=%E8%AE%B8%E6%98%95%E9%87%8D%E8%BF%94%E4%B8%96%E7%95%8C%E7%AC%AC%E4%B8%80&amp;sa=re_dl_seError_3&amp;tn=SE_fengyunbangS_fs9rizg7&amp;rsv_dl=fyb_n_erro" class="opr-item-text">许昕重返世界第一</a></td></tr><tr><td><span class="opr-index-item">4</span><a target="_blank" href="//www.baidu.com/s?word=%E5%9E%83%E5%9C%BE%E5%88%86%E7%B1%BB%E5%88%86%E5%87%BA%E9%A6%96%E9%A5%B0&amp;sa=re_dl_seError_4&amp;tn=SE_fengyunbangS_fs9rizg7&amp;rsv_dl=fyb_n_erro" class="opr-item-text">垃圾分类分出首饰</a></td></tr><tr><td><span class="opr-index-item">5</span><a target="_blank" href="//www.baidu.com/s?word=%E5%8F%96%E6%B6%88%E4%B8%80%E5%8D%A1%E9%80%9A%E5%BC%80%E5%8D%A1%E8%B4%B9&amp;sa=re_dl_seError_5&amp;tn=SE_fengyunbangS_fs9rizg7&amp;rsv_dl=fyb_n_erro" class="opr-item-text">取消一卡通开卡费</a></td></tr><tr><td><span class="opr-index-item">6</span><a target="_blank" href="//www.baidu.com/s?word=%E8%8C%83%E5%86%B0%E5%86%B0%E6%9D%8E%E6%99%A8%E8%81%9A%E9%A4%90&amp;sa=re_dl_seError_6&amp;tn=SE_fengyunbangS_fs9rizg7&amp;rsv_dl=fyb_n_erro" class="opr-item-text">范冰冰李晨聚餐</a></td></tr><tr><td><span class="opr-index-item">7</span><a target="_blank" href="//www.baidu.com/s?word=%E6%80%80%E7%89%B9%E5%A1%9E%E5%BE%B7%20%E5%BC%80%E6%8B%93%E8%80%85&amp;sa=re_dl_seError_7&amp;tn=SE_fengyunbangS_fs9rizg7&amp;rsv_dl=fyb_n_erro" class="opr-item-text">怀特塞德 开拓者</a></td></tr><tr><td><span class="opr-index-item">8</span><a target="_blank" href="//www.baidu.com/s?word=%E5%8F%B0%E9%A3%8E%E5%B0%86%E7%99%BB%E9%99%86%E6%B5%B7%E5%8D%97&amp;sa=re_dl_seError_8&amp;tn=SE_fengyunbangS_fs9rizg7&amp;rsv_dl=fyb_n_erro" class="opr-item-text">台风将登陆海南</a></td></tr><tr><td><span class="opr-index-item">9</span><a target="_blank" href="//www.baidu.com/s?word=%E8%8D%B7%E5%85%B0%E5%BC%9F%E5%A5%A5%E6%96%AF%E5%8D%A1%E8%AF%84%E5%A7%94&amp;sa=re_dl_seError_9&amp;tn=SE_fengyunbangS_fs9rizg7&amp;rsv_dl=fyb_n_erro" class="opr-item-text">荷兰弟奥斯卡评委</a></td></tr><tr><td><span class="opr-index-item">10</span><a target="_blank" href="//www.baidu.com/s?word=%E6%97%A5%E6%9C%AC%E9%87%8D%E5%90%AF%E5%95%86%E4%B8%9A%E6%8D%95%E9%B2%B8&amp;sa=re_dl_seError_10&amp;tn=SE_fengyunbangS_fs9rizg7&amp;rsv_dl=fyb_n_erro" class="opr-item-text">日本重启商业捕鲸</a></td></tr></tbody></table><div class="opr-toplist-info"><span>来源：</span><a target="_blank" href="http://www.baidu.com/link?url=sLR63PtaB7kc3YkTtzDy1k3mbTm1DXDMu-nLcijZx8DmWgOff4lBxqmY-LGDyHqw">百度风云榜</a><span>&nbsp;-&nbsp;</span><a target="_blank" href="http://www.baidu.com/link?url=01vNBVXR2eaJxETl9PI3hcrvKCcwaJIKk5FkpO7l5YI_Q_pC24ogIBoZxI0LAu5oYpAdhRH42nzxAqfui0YnHK">实时热点</a></div></div></div><div id="content_left">
            <div class="nors">
                <div class="norsSuggest">
                    <h3 class="norsTitle">很抱歉，您要访问的页面不存在！</h3>
                    <p class="norsTitle2">温馨提示：</p>
                    <ol>
                        <li>请检查您访问的网址是否正确</li>
                        <li>如果您不能确认访问的网址，请浏览<a href="//www.baidu.com/more/index.html">百度更多</a>页面查看更多网址。</li>
                        <li>回到顶部重新发起搜索</li>
                        <li>如有任何意见或建议，请及时<a href="http://qingting.baidu.com/index">反馈给我们</a>。</li>
                    </ol>
                </div>
            </div>
        </div>
    </div>
    <div id="foot">
        <span id="help" style="float:left;padding-left:121px">
            <a href="http://help.baidu.com/question" target="_blank">帮助</a>
            <a href="http://www.baidu.com/search/jubao.html" target="_blank">举报</a>
            <a href="http://jianyi.baidu.com" target="_blank">给百度提建议</a>
        </span>
    </div>
</div></body></html>
`
