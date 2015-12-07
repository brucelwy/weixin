package weixin

import (
	"bytes"
	"errors"
	"fmt"
	"text/template"
)

// 用户基本信息
const (
	UserInfoUpdateRemarkURL     = "https://api.weixin.qq.com/cgi-bin/user/info/updateremark?access_token=%s"
	UserInfoGetUserInfoURL      = "https://api.weixin.qq.com/cgi-bin/user/info?access_token=%s&openid=%s&lang=%s"
	UserInfoBatchGetUserInfoURL = "https://api.weixin.qq.com/cgi-bin/user/info/batchget?access_token=%s"
)

// UpdateUserRemark 设置用户备注名
func UpdateUserRemark(openId, remark string) (err error) {
	url := fmt.Sprintf(UserInfoUpdateRemarkURL, AccessToken())
	body := fmt.Sprintf(`{"openid":"%s","remark":"%s"}`, openId, remark)
	return Post(url, []byte(body))
}

// Lang 国家地区语言版本
type Lang string

// 微信支持的语言
const (
	LangZHCN = "zh_CN" // 简体
	LangZHTW = "zh_TW" // 繁体
	LangEN   = "en"    // 英语
)

// UserInfo 用户基本信息
type UserInfo struct {
	WeixinError
	Subscribe     int    `json:"subscribe"`
	OpenId        string `json:"openid"`
	NickName      string `json:"nickname"`
	Sex           int    `json:"sex"`
	Language      Lang   `json:"language"`
	City          string `json:"city"`
	Province      string `json:"province"`
	Country       string `json:"country"`
	HeadImgURL    string `json:"headimgurl"`
	SubscribeTime int    `json:"subscribe_time"`
	UnionId       string `json:"unionid"`
	Remark        string `json:"remark"`
	GroupId       int    `json:"groupid"`
}

// GetUserInfo 获取用户基本信息（包括UnionID机制）
func GetUserInfo(openId string, lang ...Lang) (info *UserInfo, err error) {
	if len(lang) == 0 {
		lang = []Lang{LangZHCN}
	}
	url := fmt.Sprintf(UserInfoGetUserInfoURL, AccessToken(), openId, lang[0])

	info = &UserInfo{}
	err = GetUnmarshal(url, info)
	if err != nil {
		return nil, err
	}

	return info, nil
}

// BatchGetUserInfo 获取用户基本信息（包括UnionID机制）
func BatchGetUserInfo(openIds []string, lang ...Lang) (infos []UserInfo, err error) {
	if len(openIds) == 0 {
		return nil, errors.New("openIds is blank")
	}

	if len(lang) == 0 {
		lang = []Lang{LangZHCN}
	}
	url := fmt.Sprintf(UserInfoBatchGetUserInfoURL, AccessToken())

	tlp := `{"user_list":[{{range $index, $elmt := .}}{{if $index}},{{end}}{"openid":"{{.}}","lang":"` + string(lang[0]) + `"}{{end}}]}`
	t := template.New("user_info_list_request")
	t, _ = t.Parse(tlp)
	var buf bytes.Buffer
	t.Execute(&buf, openIds)

	wapper := &struct {
		WeixinError
		UserInfoList []UserInfo `json:"user_info_list"`
	}{}
	err = PostUnmarshal(url, buf.Bytes(), wapper)
	if err != nil {
		return nil, err
	}

	return wapper.UserInfoList, nil
}
