package qqwry

import (
	"bytes"
	"encoding/binary"
	"io/ioutil"
	"net"

	"os"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

const (
	INDEX_LEN       = 7
	REDIRECT_MODE_1 = 0x01
	REDIRECT_MODE_2 = 0x02
)

type QQwry struct {
	filepath string
	file     *os.File
}

type QQwryResult struct {
	IP      string
	Country string
	City    string
}

func NewQQwry(filepath string) (qqwry *QQwry) {
	qqwry = &QQwry{filepath: filepath}
	return
}

func GbkToUtf8(s []byte) string {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return ""
	}
	return string(d)
}

func (this *QQwry) ReadMode(offset uint32) byte {
	this.file.Seek(int64(offset), 0)
	mode := make([]byte, 1)
	this.file.Read(mode)
	return mode[0]
}

func (this *QQwry) ReadArea(offset uint32) []byte {
	mode := this.ReadMode(offset)
	if mode == REDIRECT_MODE_1 || mode == REDIRECT_MODE_2 {
		areaOffset := this.ReadUInt24()
		if areaOffset == 0 {
			return []byte("")
		} else {
			return this.ReadString(areaOffset)
		}
	}
	return this.ReadString(offset)
}

func (this *QQwry) ReadString(offset uint32) []byte {
	this.file.Seek(int64(offset), 0)
	data := make([]byte, 0, 30)
	buf := make([]byte, 1)
	for {
		this.file.Read(buf)
		if buf[0] == 0 {
			break
		}
		data = append(data, buf[0])
	}
	return data
}

func (this *QQwry) SearchIndex(ip uint32) uint32 {
	header := make([]byte, 8)
	this.file.Seek(0, 0)
	this.file.Read(header)

	start := binary.LittleEndian.Uint32(header[:4])
	end := binary.LittleEndian.Uint32(header[4:])

	for {
		mid := this.GetMiddleOffset(start, end)
		this.file.Seek(int64(mid), 0)
		buf := make([]byte, INDEX_LEN)
		this.file.Read(buf)
		_ip := binary.LittleEndian.Uint32(buf[:4])

		if end-start == INDEX_LEN {
			offset := byte3ToUInt32(buf[4:])
			this.file.Read(buf)
			if ip < binary.LittleEndian.Uint32(buf[:4]) {
				return offset
			} else {
				return 0
			}
		}

		if _ip > ip {
			end = mid
		} else if _ip < ip {
			start = mid
		} else if _ip == ip {
			return byte3ToUInt32(buf[4:])
		}
	}
}

func (this *QQwry) ReadUInt24() uint32 {
	buf := make([]byte, 3)
	this.file.Read(buf)
	return byte3ToUInt32(buf)
}

func (this *QQwry) GetMiddleOffset(start uint32, end uint32) uint32 {
	records := ((end - start) / INDEX_LEN) >> 1
	return start + records*INDEX_LEN
}

func byte3ToUInt32(data []byte) uint32 {
	i := uint32(data[0]) & 0xff
	i |= (uint32(data[1]) << 8) & 0xff00
	i |= (uint32(data[2]) << 16) & 0xff0000
	return i
}

func (this *QQwry) Find(ip string) (result QQwryResult, err error) {
	if this.filepath == "" {
		return
	}

	file, err := os.OpenFile(this.filepath, os.O_RDONLY, 0400)
	defer file.Close()
	if err != nil {
		return
	}
	this.file = file
	ipv4 := net.ParseIP(ip).To4()
	ipv4long := binary.BigEndian.Uint32(ipv4)
	offset := this.SearchIndex(ipv4long)
	if offset <= 0 {
		return
	}

	var country []byte
	var area []byte

	mode := this.ReadMode(offset + 4)
	if mode == REDIRECT_MODE_1 {
		countryOffset := this.ReadUInt24()
		mode = this.ReadMode(countryOffset)
		if mode == REDIRECT_MODE_2 {
			c := this.ReadUInt24()
			country = this.ReadString(c)
			countryOffset += 4
		} else {
			country = this.ReadString(countryOffset)
			countryOffset += uint32(len(country) + 1)
		}
		area = this.ReadArea(countryOffset)
	} else if mode == REDIRECT_MODE_2 {
		countryOffset := this.ReadUInt24()
		country = this.ReadString(countryOffset)
		area = this.ReadArea(offset + 8)
	} else {
		country = this.ReadString(offset + 4)
		area = this.ReadArea(offset + uint32(5+len(country)))
	}

	result = QQwryResult{
		IP:      ip,
		Country: GbkToUtf8(country),
		City:    GbkToUtf8(area),
	}
	return result, nil
}
