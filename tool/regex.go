/*
@Desc : 正则工具
@Version : 1.0.0
@Time : 2020/9/3 11:23
@Author : hammercui
@File : regex
@Company: Sdbean
*/
package tool

import "regexp"

func VerifyEmailFormat(email string) bool {
	//pattern := `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*` //匹配电子邮箱
	pattern := `^[0-9a-z][_.0-9a-z-]{0,31}@([0-9a-z][0-9a-z-]{0,30}[0-9a-z]\.){1,4}[a-z]{2,4}$`

	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}
//mobile verify
func VerifyMobileFormat(mobileNum string) bool {
	//国家
	regular := "^86"
	reg := regexp.MustCompile(regular)
	ret := reg.MatchString(mobileNum)
	if(ret){
		return ret
	}
	//手机号
	regular = "1[34578][0-9]{9}"
	reg = regexp.MustCompile(regular)
	return reg.MatchString(mobileNum)
}