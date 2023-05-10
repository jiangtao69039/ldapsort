package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"path"
	"regexp"
	"sort"
	"strings"
)
func init()  {

}

func main() {

	//err := fortressFileSort10("/home/tttz/Desktop/rs/exp.ldif","/home/tttz/Desktop/rs/result_sort.ldif")
	//str, _ := os.Getwd()
	err :=fortressFileSort_decodebase64(path.Join("/home/tttz/Desktop/0009.txt"),path.Join("/home/jiangtao/Desktop/0601ldap","data_sort_decode.ldif"))
	fmt.Println(err)

}

type Fortress_Block struct {
	content               []string //这一块的文本内容
	dn_count				int
	dn_line					string
    id                      int
}

func NewFortressBlock() *Fortress_Block {
	return &Fortress_Block{
		content: make([]string,0),
		dn_count: 0,
		dn_line:"",

	}
}

func fortressFileSort_decodebase64(src string, desc string) error {
	fi, err := os.Open(src)
	if err != nil {
		return err
	}
	defer fi.Close()

	reader := bufio.NewReader(fi)

	CRLF_Reg, _ := regexp.Compile("^\\s+$")
	dn_reg, _ := regexp.Compile("^dn::?\\s(.*)\\s+$")
	continue_line, _ := regexp.Compile("^\\s+(.+)")

	var blpck_ptr *Fortress_Block = nil
	blocks := make([]*Fortress_Block, 0)
	curr_dn_line := ""
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			if blpck_ptr != nil {
				blocks = append(blocks, blpck_ptr)
			}
			break
		}
		//dn块开始
		if CRLF_Reg.MatchString(line) {
			if blpck_ptr == nil {
				blpck_ptr = NewFortressBlock()
			} else {
				blocks = append(blocks, blpck_ptr)
				blpck_ptr = NewFortressBlock()
			}
		} else if dn_reg.MatchString(line) {
			//dn行匹配到
			curr_dn_line += line
			blpck_ptr.content = append(blpck_ptr.content, line)
			blpck_ptr.dn_line = curr_dn_line
		} else if continue_line.MatchString(line) && curr_dn_line != "" {
			//dn行第二..行
			curr_dn_line += line
			blpck_ptr.content = append(blpck_ptr.content, line)
			blpck_ptr.dn_line = curr_dn_line
		} else {
			if blpck_ptr == nil {
				blpck_ptr = NewFortressBlock()
			}
			curr_dn_line = ""
			//可能是行内 行外数据,最后用dn判断
			blpck_ptr.content = append(blpck_ptr.content, line)
		}
	}
	i := 0
	for _, value := range blocks {
		value.id = i
		i++
	}
	needSortBlocks := make([]*Fortress_Block, 0)
	base64reg, _ := regexp.Compile("^dn::?\\s*")
	blankreg, _ := regexp.Compile("\\s+")

	for _, value := range blocks {
		if value.dn_line != "" {
			value.dn_line = strings.ReplaceAll(value.dn_line, "\n", "")
			value.dn_line = blankreg.ReplaceAllString(value.dn_line, "")
			if !strings.Contains(value.dn_line, ",") {
				base64Line := base64reg.ReplaceAllString(value.dn_line, "")
				b, err := base64.StdEncoding.DecodeString(base64Line)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				value.dn_line = string(b)
			}
			value.dn_count = len(strings.Split(value.dn_line, ","))
			needSortBlocks = append(needSortBlocks, value)
		}
	}

	sort.Slice(needSortBlocks, func(i, j int) bool {
		return needSortBlocks[i].dn_count < needSortBlocks[j].dn_count
	})



	//TODO文件是否存在
	_ = os.Remove(desc)
	_, _ = os.Create(desc)
	fin, err := os.OpenFile(desc, os.O_TRUNC|os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer fin.Close()
	base64line_reg, _ := regexp.Compile("(.*)::?")
	for _, value := range needSortBlocks {
		if value.dn_line == "" {
			continue
		} else {
			temp := value
			skipDn := false
			for _, v := range temp.content {
				fmt.Print(v)
				if(!skipDn){
					//还没有输出dn
					_, _ = fin.WriteString(v)
					if(strings.HasPrefix(v,"dn:")){
						skipDn=true
					}
					continue
				}
				if(strings.Contains(v,"::")){
					//说明是单行base64,目前没发现除DN外有2行的base64
					findString := base64line_reg.FindString(v)
					findString = strings.ReplaceAll(findString,"::",": ")
					v = base64line_reg.ReplaceAllString(v,"")
					v= blankreg.ReplaceAllString(v,"")
					b, _ := base64.StdEncoding.DecodeString(v)
					result := findString + string(b)+"\n"
					_, _ = fin.WriteString(result)
				}else {
					_, _ = fin.WriteString(v)
				}

			}
			//needSortBlocks = needSortBlocks[1:]
		}
		_, _ = fin.WriteString("62-CLOUDQUERY-SPLIT\n")

	}

	return nil
}

func fortressFileSort10(src string, desc string) error {
	fi, err := os.Open(src)
	if err != nil {
		return err
	}
	defer fi.Close()

	reader := bufio.NewReader(fi)

	CRLF_Reg, _ := regexp.Compile("^\\s+$")
	dn_reg, _ := regexp.Compile("^dn::?\\s(.*)\\s+$")
	continue_line, _ := regexp.Compile("^\\s+(.+)")

	var blpck_ptr *Fortress_Block = nil
	blocks := make([]*Fortress_Block, 0)
	curr_dn_line := ""
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			if blpck_ptr != nil {
				blocks = append(blocks, blpck_ptr)
			}
			break
		}
		//dn块开始
		if CRLF_Reg.MatchString(line) {
			if blpck_ptr == nil {
				blpck_ptr = NewFortressBlock()
			} else {
				blocks = append(blocks, blpck_ptr)
				blpck_ptr = NewFortressBlock()
			}
		} else if dn_reg.MatchString(line) {
			//dn行匹配到
			curr_dn_line += line
			blpck_ptr.content = append(blpck_ptr.content, line)
			blpck_ptr.dn_line = curr_dn_line
		} else if continue_line.MatchString(line) && curr_dn_line != "" {
			//dn行第二..行
			curr_dn_line += line
			blpck_ptr.content = append(blpck_ptr.content, line)
			blpck_ptr.dn_line = curr_dn_line
		} else {
			if blpck_ptr == nil {
				blpck_ptr = NewFortressBlock()
			}
			curr_dn_line = ""
			//可能是行内 行外数据,最后用dn判断
			blpck_ptr.content = append(blpck_ptr.content, line)
		}
	}
	i := 0
	for _, value := range blocks {
		value.id = i
		i++
	}
	needSortBlocks := make([]*Fortress_Block, 0)
	base64reg, _ := regexp.Compile("^dn::?\\s*")
	blankreg, _ := regexp.Compile("\\s+")

	for _, value := range blocks {
		if value.dn_line != "" {

			value.dn_line = strings.ReplaceAll(value.dn_line, "\n", "")
			value.dn_line = blankreg.ReplaceAllString(value.dn_line, "")
			if !strings.Contains(value.dn_line, ",") {
				base64Line := base64reg.ReplaceAllString(value.dn_line, "")
				b, err := base64.StdEncoding.DecodeString(base64Line)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				value.dn_line = string(b)
			}else {
				if strings.HasPrefix(value.dn_line,"dn:") {
					value.dn_line=value.dn_line[3:]
				}
			}
			value.dn_count = len(strings.Split(value.dn_line, ","))
			needSortBlocks = append(needSortBlocks, value)
		}
	}

	//

	//
	sort.Slice(needSortBlocks, func(i, j int) bool {
		return needSortBlocks[i].dn_count < needSortBlocks[j].dn_count
	})

	//
	dn_map := make(map[string]string)
	for _,value := range needSortBlocks{
		dn_map[value.dn_line]=value.dn_line
	}

	notExistBlocks := make([]*Fortress_Block, 0)
	result := make([]*Fortress_Block, 0)

	for _,value := range needSortBlocks{
		if value.dn_count <=2 {
			result= append(result, value)
			continue;
		}
		parentDN := value.dn_line[strings.Index(value.dn_line,",")+1:]
		if _, ok := dn_map[parentDN]; ! ok {
			// 不存在
			notExistBlocks=append(notExistBlocks, value)
		}else {
			result= append(result, value)
		}

	}

	needSortBlocks = result

	//TODO文件是否存在
	_ = os.Remove(desc)
	_, _ = os.Create(desc)
	fin, err := os.OpenFile(desc, os.O_TRUNC|os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer fin.Close()
	for _, value := range needSortBlocks {

			for _, v := range value.content {
				fmt.Print(v)
				_, _ = fin.WriteString(v)
			}

		_, _ = fin.WriteString("\n")

	}


	//-----------------------------
	_ = os.Remove("/home/jiangtao/Desktop/rs/notexist.ldif")
	_, _ = os.Create("/home/jiangtao/Desktop/rs/notexist.ldif")
	fin, err = os.OpenFile("/home/jiangtao/Desktop/rs/notexist.ldif", os.O_TRUNC|os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer fin.Close()
	for _, value := range notExistBlocks {

			for _, v := range value.content {
				fmt.Print(v)
				_, _ = fin.WriteString(v)
			}
		_, _ = fin.WriteString("\n")

	}

	return nil
}

