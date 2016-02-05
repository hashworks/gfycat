package gfycat

import (
	"math/rand"
	"testing"
	"time"
)

const (
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var src = rand.NewSource(time.Now().UnixNano())

func generateRandomString(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

func TestGifUploadByURL(t *testing.T) {
	item, err := UploadGIFByURL("https://zippy.gfycat.com/DeficientFluffyBluefish.gif", "")
	if err != nil {
		t.Error(err.Error())
	} else {
		if item.WEBMUrl == "" {
			t.Fail()
		}
	}
}

func TestGIFUploadByFile(t *testing.T) {
	item, err := UploadGIFByFile("gfycat_test.gif", generateRandomString(10))
	if err != nil {
		t.Error(err.Error())
	} else {
		if item.WEBMUrl == "" {
			t.Fail()
		}
	}
}

func TestGetGfyInfoByRandomString(t *testing.T) {
	response, err := GetGfyInfoByGfyName("DeficientFluffyBluefish")
	if err != nil {
		if err.Error() == "json: cannot unmarshal string into Go value of type int" {
			t.Error("This fails because of weird types in the Gfycat JSON response. I informed them about this.")
		} else {
			t.Error(err.Error())
		}
	} else {
		if response.GfyItem.WEBMUrl == "" {
			t.Fail()
		}
	}
}

func TestCheckGIFByURL(t *testing.T) {
	item, err := CheckGIFByURL("http://i.imgur.com/lKi99vn.gif")
	if err != nil {
		if err.Error() == "json: cannot unmarshal string into Go value of type int" {
			t.Error("This fails because of weird types in the Gfycat JSON response. I informed them about this.")
		} else {
			t.Error(err.Error())
		}
	} else {
		if !item.URLKnown {
			t.Fail()
		}
	}
}

func TestGetOEmbedDataByGfyURL(t *testing.T) {
	oEmbed, err := GetOEmbedDataByGfyURL("https://gfycat.com/DeficientFluffyBluefish", 512, 512)
	if err != nil {
		t.Error(err.Error())
	} else {
		if oEmbed.HTML == "" {
			t.Fail()
		}
	}
}
