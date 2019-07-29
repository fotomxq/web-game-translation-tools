package module

import (
	"archive/zip"
	"crypto/sha1"
	"encoding/hex"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"sort"
	"strings"
	"time"
)

//方法集合
type FileType struct {
}

var (
	//分隔符
	Sep string = string(os.PathSeparator)
	//方法集
	File FileType
)

//创建临时文件到指定目录
// 文件名称以文件内容的SHA1为主
//param dirSrc string 存放目录 eg : abc/
//param content []byte 文件内容
//param FileType string 文件类型
//return string,string,string,error 存储路径，相对路径，文件名称，错误
func (File *FileType) CreateDownloadLS(dirSrc string, content []byte, FileType string) (string, string, string, error) {
	//计算SHA1
	hasher := sha1.New()
	_, err := hasher.Write(content)
	if err != nil {
		return "", "", "", err
	}
	sha := hasher.Sum(nil)
	shaStr := hex.EncodeToString(sha)
	//创建文件路径
	fileName := shaStr + "." + FileType
	nowTime := time.Now().Format("2006010215")
	dsrc := "ls" + Sep + nowTime
	src := dirSrc + Sep + dsrc
	err = File.CreateFolder(src)
	src += Sep + fileName
	//删除60分钟之前的数据
	sinceTime := time.Now().Add(-time.Minute * 30).Format("2006010215")
	sinceTimeDSrc := dirSrc + Sep + sinceTime
	fileList, err := File.GetFileList(dirSrc, []string{}, false)
	if err == nil {
		for _, v := range fileList {
			if v == nowTime || v == sinceTime {
				continue
			}
			err = File.DeleteF(sinceTimeDSrc + Sep + "v")
			if err != nil {
				return src, dsrc, fileName, File.WriteFile(src, content)
			}
		}
	}
	return src, dsrc, fileName, File.WriteFile(src, content)
}

//从互联网下载文件
func (File *FileType) DownloadByURL(url string, params url.Values, destDir string, name string) error {
	//从URL下载数据结构
	resp,err := HttpGet.GetData(url,params,"", false)
	if err != nil{
		return err
	}
	if File.IsFolder(destDir) == false{
		if err := File.CreateFolder(destDir); err != nil{
			return err
		}
	}
	//将数据保存到指定的文件
	err = File.WriteFile(destDir + Sep + name, resp)
	if err != nil{
		return err
	}
	//反馈成功
	return nil
}

//移动文件或文件夹
//param src string 文件路径
//param dest string 新路径
//return error
func (File *FileType) MoveF(src string, dest string) error {
	return os.Rename(src, dest)
}

//删除文件或文件夹
//param src string 文件路径
//return error
func (File *FileType) DeleteF(src string) error {
	return os.RemoveAll(src)
}

//判断文件或文件夹是否存在
//param src string 文件路径
//return bool 是否存在
func (File *FileType) IsExist(src string) bool {
	_, err := os.Stat(src)
	return err == nil && os.IsNotExist(err) == false
}

//读取文件
//param src string 文件路径
//return []byte,error 文件数据,错误
func (File *FileType) LoadFile(src string) ([]byte, error) {
	fd, err := os.Open(src)
	if err != nil {
		return nil, err
	}
	defer fd.Close()
	c, err := ioutil.ReadAll(fd)
	if err != nil {
		return nil, err
	}
	return c, nil
}

//写入文件
//param src string 文件路径
//param content []byte 写入内容
//return error
func (File *FileType) WriteFile(src string, content []byte) error {
	return ioutil.WriteFile(src, content, 0666)
}

//追加写入文件
//param src string 文件路径
//param content []byte 写入内容
//return error
func (File *FileType) WriteFileAppend(src string, content []byte) error {
	f, err := os.OpenFile(src, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(content)
	if err != nil {
		return err
	}
	return nil
}

//复制文件
//param src string 文件路径
//param dest string 新路径
//return error
func (File *FileType) CopyFile(src string, dest string) error {
	srcF, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcF.Close()
	destF, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destF.Close()
	_, err = io.Copy(destF, srcF)
	if err != nil {
		return err
	}
	return nil
}

//判断是否为文件
//param src string 文件路径
//return bool 是否为文件
func (File *FileType) IsFile(src string) bool {
	info, err := os.Stat(src)
	return err == nil && !info.IsDir()
}

//创建多级文件夹
//param src string 新文件夹路径
//return error
func (File *FileType) CreateFolder(src string) error {
	return os.MkdirAll(src, os.ModePerm)
}

//判断是否为文件夹
//param src string 文件夹路径
//return bool 是否为文件夹
func (File *FileType) IsFolder(src string) bool {
	info, err := os.Stat(src)
	return err == nil && info.IsDir()
}

//复制文件夹
// 自递归复制文件夹
//param src string 源路径
//param dest string 目标路径
//return bool 是否成功
func (File *FileType) CopyFolder(src string, dest string) bool {
	dir, err := ioutil.ReadDir(src)
	if err != nil {
		return false
	}
	if File.CreateFolder(dest) != nil {
		return false
	}
	for _, v := range dir {
		vSrc := src + Sep + v.Name()
		vDest := dest + Sep + v.Name()
		if v.IsDir() == true {
			if File.CreateFolder(vDest) != nil {
				return false
			}
			if File.CopyFolder(vSrc, vDest) == false {
				return false
			}
		} else {
			if File.CopyFile(vSrc, vDest) != nil {
				return false
			}
		}
	}
	return true
}

//获取文件信息
//param src string 文件路径
//return os.FileInfo,error 文件信息，错误
func (File *FileType) GetFileInfo(src string) (os.FileInfo, error) {
	c, err := os.Stat(src)
	return c, err
}

//获取文件大小
//param src string 文件路径
//return int64,bool 文件大小，错误
func (File *FileType) GetFileSize(src string) (int64, error) {
	info, err := os.Stat(src)
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}

//获取文件名称分割序列
//param src string 文件路径
//return map[string]string,error 文件名称序列，错误 eg : {"name","abc.jpg","type":"jpg","only-name":"abc"}
func (File *FileType) GetFileNames(src string) (map[string]string, error) {
	info, err := os.Stat(src)
	if err != nil {
		return nil, err
	}
	res := map[string]string{
		"name":      info.Name(),
		"type":      "",
		"only-name": info.Name(),
	}
	names := strings.Split(res["name"], ".")
	if len(names) < 2 {
		return res, nil
	}
	res["type"] = names[len(names)-1]
	res["only-name"] = names[0]
	for i := range names {
		if i != 0 && i < len(names)-1 {
			res["only-name"] = res["only-name"] + "." + names[i]
		}
	}
	return res, nil
}

//获取文件列表
// 按照文件名，倒叙排列返回
//param src string 查询的文件夹路径
//param filtre []string 仅保留的文件，文件夹除外
//param isSrc bool 返回是否为文件路径
//return []string,error 文件列表,错误
func (File *FileType) GetFileList(src string, filters []string, isSrc bool) ([]string, error) {
	//初始化
	var fs []string
	//读取目录
	dir, err := ioutil.ReadDir(src)
	if err != nil {
		return nil, err
	}
	//遍历目录文件
	for _, v := range dir {
		var appendSrc string
		if isSrc == true {
			appendSrc = src + Sep + v.Name()
		} else {
			appendSrc = v.Name()
		}
		if v.IsDir() == true || len(filters) < 1 {
			fs = append(fs, appendSrc)
			continue
		}
		names := strings.Split(v.Name(), ".")
		if len(names) == 1 {
			fs = append(fs, appendSrc)
			continue
		}
		t := names[len(names)-1]
		for _, filterValue := range filters {
			if t != filterValue {
				continue
			}
			fs = append(fs, appendSrc)
		}
	}
	//对数组进行倒叙排序
	sort.Sort(sort.Reverse(sort.StringSlice(fs)))
	//返回
	return fs, nil
}

//查询文件夹下文件个数
//param src string 文件夹路径
//return int,error 文件个数,错误
func (File *FileType) GetFileListCount(src string) (int, error) {
	dir, err := ioutil.ReadDir(src)
	if err != nil {
		return 0, err
	}
	var res int
	for range dir {
		res += 1
	}
	return res, nil
}

//获取文件SHA1值
//param src string 文件路径
//return string,error SHA1值,错误
func (File *FileType) GetFileSha1(src string) (string, error) {
	content, err := File.LoadFile(src)
	if err != nil {
		return "", err
	}
	if content != nil {
		sha := sha1.New()
		_, err = sha.Write(content)
		if err != nil {
			return "", err
		}
		res := sha.Sum(nil)
		return hex.EncodeToString(res), nil
	}
	return "", nil
}

//获取并创建时间序列创建的多级文件夹
//eg : Return and create the path ,"[src]/201611/"
//eg : Return and create the path ,"[src]/201611/2016110102-03[appendFileType]"
//param src string 文件路径
//param appendFileType string 是否末尾追加文件类型，如果指定值，则返回
//return string,error 新时间周期目录，错误
func (File *FileType) GetTimeDirSrc(src string, appendFileType string) (string, error) {
	newSrc := src + Sep + time.Now().Format("200601")
	err := File.CreateFolder(newSrc)
	if err != nil {
		return "", err
	}
	newSrc = newSrc + Sep
	if appendFileType != "" {
		newSrc = newSrc + time.Now().Format("20060102-03") + appendFileType
	}
	return newSrc, nil
}

//压缩文件夹
//param src string 源文件
//param zipSrc string 目标压缩包
//return error 错误信息
func (File *FileType) ZipDir(src string, zipSrc string) error {
	//构建ZIP文件
	d, _ := os.Create(zipSrc)
	defer d.Close()
	w := zip.NewWriter(d)
	defer w.Close()
	//读取目录
	dir, err := ioutil.ReadDir(src)
	if err != nil {
		return err
	}
	//遍历文件
	for _, fileSrc := range dir {
		file, err := os.Open(src + Sep + fileSrc.Name())
		defer file.Close()
		err = File.ZipDirC(file, "", w)
		if err != nil {
			return err
		}
	}
	return nil
}

//压缩目录子操作结构
//param file *os.File 文件句柄
//param prefix string
//param zw *zip.Writer 写入ZIP句柄
//return error 错误代码
func (File *FileType) ZipDirC(file *os.File, prefix string, zw *zip.Writer) error {
	info, err := file.Stat()
	if err != nil {
		return err
	}
	if info.IsDir() {
		prefix = prefix + "/" + info.Name()
		info, err := file.Readdir(-1)
		if err != nil {
			return err
		}
		for _, vInfo := range info {
			f, err := os.Open(file.Name() + "/" + vInfo.Name())
			if err != nil {
				return err
			}
			defer f.Close()
			err = File.ZipDirC(f, prefix, zw)
			if err != nil {
				return err
			}
		}
	} else {
		header, err := zip.FileInfoHeader(info)
		header.Name = prefix + "/" + header.Name
		if err != nil {
			return err
		}
		writer, err := zw.CreateHeader(header)
		if err != nil {
			return err
		}
		_, err = io.Copy(writer, file)
		file.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

//解压文件
//param zipSrc string 目标压缩包
//param dest string 解压到... eg : /dir/
//return error 错误信息
func (File *FileType) UnZip(zipSrc string, dest string) error {
	reader, err := zip.OpenReader(zipSrc)
	if err != nil {
		return err
	}
	defer reader.Close()
	for _, file := range reader.File {
		rc, err := file.Open()
		if err != nil {
			return err
		}
		filename := dest + file.Name
		err = os.MkdirAll(File.GetDir(filename), 0755)
		if err != nil {
			rc.Close()
			return err
		}
		w, err := os.Create(filename)
		if err != nil {
			rc.Close()
			return err
		}
		_, err = io.Copy(w, rc)
		if err != nil {
			rc.Close()
			w.Close()
			return err
		}
		w.Close()
		rc.Close()
	}
	return nil
}

//获取目录路径
//param path string 地址路径
//return string 返回值
func (File *FileType) GetDir(path string) string {
	return File.SubString(path, 0, strings.LastIndex(path, "/"))
}

//截取字符串
//param str string 字符串
//param start int 开始位置
//param end int 结束位置
//return string 结果字符串
func (File *FileType) SubString(str string, start, end int) string {
	rs := []rune(str)
	length := len(rs)

	if start < 0 || start > length {
		panic("start is wrong")
	}

	if end < start || end > length {
		panic("end is wrong")
	}

	return string(rs[start:end])
}
