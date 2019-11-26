package baidu_pc_service_test

import (
	"fmt"
	"html_parse_api/services/baidu_pc_service"
	"testing"
)

func TestDecodeBaiduEncURL(t *testing.T) {
	fmt.Println(baidu_pc_service.DecodeBaiduEncURL("https://m.baidu.com/from=844b/bd_page_type=1/ssid=0/uid=0/pu=usm%404%2Csz%401320_2001%2Cta%40iphone_1_11.0_25_11.0/baiduid=EB9A2C07DAED944FBF46F430DA5A1064/w=0_10_/t=iphone/l=1/tc?ref=www_iphone&lid=8180081909666280377&order=2&fm=alop&isAtom=1&is_baidu=0&tj=vmp_zxent_atom_2_0_10_l1&clk_info=%7B%22tplname%22%3A%22vmp_zxent_atom%22%7D&wd=&eqid=71857d282b41f400100000005db67c82&w_qd=IlPT2AEptyoA_ykxxu5czBqvJlxSom9nw5MXmQDRtucsMOUuUkFaxZnffOe&bdver=2&tcplug=1&sec=42533&di=b113792dd4017667&bdenc=1&nsrc=FnwSf2lhoyucCAidpUHIxHW31uZkM99rR4GSQmUoOagP8lyEuUuszvEUxkOcaoVMqypnznnNjXwx%2BgARILRC7C0epdk16opLL054spyk8L4Vr3Z4JcR%2BTbJ5rTXTzD%2BCn3iRdKgKWF6O3N8vteRV5g%3D%3D"))
}
