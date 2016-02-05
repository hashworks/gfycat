package gfycat

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type gfyError struct {
	GIFName string
	Error   string
	Task    string
	Time    int
}

type gfyResponse struct {
	GfyItem gfyItem
}

type gfyItem struct {
	GfyId              string
	GfyName            string
	GfyNumber          string
	UserName           string
	Width              int
	Height             int
	FrameRate          int
	NumFrames          int
	MP4Url             string
	WEBMUrl            string
	GIFUrl             string
	GIFSize            int
	Mp4Size            int
	WEBMSize           int
	CreateDate         string
	Views              int
	Title              string
	ExtraLemmas        []string
	MD5                string
	Tags               []string
	NSFW               string
	SAR                string
	URL                string
	Source             string
	Dynamo             string
	Subreddit          string
	RedditId           string
	RedditIdText       string
	UploadGifName      string
	Likes              int
	Dislikes           int
	Published          int
	Description        string
	ExtraLemmaText     string
	UrlMD5             string
	Task               string
	GfySize            int
	CopyrightClaimaint string
}

type gfyCheckURLResponse struct {
	URLKnown  bool
	GfyName   string
	GfyURL    string
	WEBMUrl   string
	MP4Url    string
	GIFUrl    string
	FrameRate int
}

type gfyOEmbed struct {
	Version      string
	Type         string
	ProviderName string `json:"provider_name"`
	Width        int
	Height       int
	Title        string
	HTML         string
}

func getGfyError(jsonData []byte) error {
	if strings.Index(string(jsonData), `"error":`) != -1 {
		var gfyError gfyError
		err := json.Unmarshal(jsonData, &gfyError)
		if err == nil {
			return errors.New(gfyError.Error)
		}
		return nil
	}
	return nil
}

// UploadGIFByURL uploads a GIF by URL and returns information about the generated item.
// You can provide an optional random string for the GfyName.
func UploadGIFByURL(gifURL string, ownRandomString string) (gfyItem, error) {
	var requestURL string
	if ownRandomString != "" {
		requestURL = "https://upload.gfycat.com/transcodeRelease/" + ownRandomString + "?fetchUrl=" + url.QueryEscape(gifURL)
	} else {
		requestURL = "https://upload.gfycat.com/transcode?fetchUrl=" + url.QueryEscape(gifURL)
	}
	response, err := http.DefaultClient.Get(requestURL)
	if err != nil {
		return gfyItem{}, err
	}
	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return gfyItem{}, err
	}
	err = getGfyError(bytes)
	if err != nil {
		return gfyItem{}, err
	}
	var item gfyItem
	err = json.Unmarshal(bytes, &item)
	if err != nil {
		return gfyItem{}, err
	}
	return item, nil
}

// UploadGIFByFile uploads a GIF by a file path and returns information about the generated item.
// You must provide a random string for the GfyName.
func UploadGIFByFile(gifFilePath, ownRandomString string) (gfyItem, error) {
	if ownRandomString == "" {
		return gfyItem{}, errors.New("You must provide your own random string.")
	}
	if gifFilePath == "" {
		return gfyItem{}, errors.New("You must provide a GIF file path.")
	}
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	w.WriteField("Content-Type", "image/gif")
	w.WriteField("acl", "private")
	w.WriteField("AWSAccessKeyId", "AKIAIT4VU4B7G2LQYKZQ")
	w.WriteField("policy", "eyAiZXhwaXJhdGlvbiI6ICIyMDIwLTEyLTAxVDEyOjAwOjAwLjAwMFoiLAogICAgICAgICAgICAiY29uZGl0aW9ucyI6IFsKICAgICAgICAgICAgeyJidWNrZXQiOiAiZ2lmYWZmZSJ9LAogICAgICAgICAgICBbInN0YXJ0cy13aXRoIiwgIiRrZXkiLCAiIl0sCiAgICAgICAgICAgIHsiYWNsIjogInByaXZhdGUifSwKCSAgICB7InN1Y2Nlc3NfYWN0aW9uX3N0YXR1cyI6ICIyMDAifSwKICAgICAgICAgICAgWyJzdGFydHMtd2l0aCIsICIkQ29udGVudC1UeXBlIiwgIiJdLAogICAgICAgICAgICBbImNvbnRlbnQtbGVuZ3RoLXJhbmdlIiwgMCwgNTI0Mjg4MDAwXQogICAgICAgICAgICBdCiAgICAgICAgICB9")
	w.WriteField("success_action_status", "200")
	w.WriteField("signature", "mk9t/U/wRN4/uU01mXfeTe2Kcoc=")
	w.WriteField("key", ownRandomString)
	f, err := os.Open(gifFilePath)
	if err != nil {
		return gfyItem{}, err
	}
	fw, err := w.CreateFormFile("file", filepath.Base(gifFilePath))
	if err != nil {
		return gfyItem{}, err
	}
	if _, err = io.Copy(fw, f); err != nil {
		return gfyItem{}, err
	}
	w.Close()

	req, err := http.NewRequest("POST", "https://gifaffe.s3.amazonaws.com/", &b)
	if err != nil {
		return gfyItem{}, err
	}
	req.Header.Add("Content-Type", w.FormDataContentType())
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return gfyItem{}, err
	}
	if response.ContentLength != 0 {
		bytes, err := ioutil.ReadAll(response.Body)
		if err == nil {
			err = errors.New(string(bytes))
		}
		return gfyItem{}, err
	}
	response, err = http.DefaultClient.Get("https://upload.gfycat.com/transcode/" + url.QueryEscape(ownRandomString))
	bytes, err := ioutil.ReadAll(response.Body)
	err = getGfyError(bytes)
	if err != nil {
		return gfyItem{}, err
	}
	var item gfyItem
	err = json.Unmarshal(bytes, &item)
	if err != nil {
		return gfyItem{}, err
	}
	return item, nil
}

// GetGfyInfoByRandomString returns information about the provided GfyName.
func GetGfyInfoByGfyName(gfyName string) (gfyResponse, error) {
	response, err := http.DefaultClient.Get("https://gfycat.com/cajax/get/" + url.QueryEscape(gfyName))
	if err != nil {
		return gfyResponse{}, err
	}
	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return gfyResponse{}, err
	}
	err = getGfyError(bytes)
	if err != nil {
		return gfyResponse{}, err
	}
	var rsp gfyResponse
	err = json.Unmarshal(bytes, &rsp)
	if err != nil {
		return gfyResponse{}, err
	}
	return rsp, nil
}

// CheckGIFByURL checks if a provided GIF URL was already posted to gfycat.com.
func CheckGIFByURL(gifURL string) (gfyCheckURLResponse, error) {
	response, err := http.DefaultClient.Get("https://gfycat.com/cajax/checkUrl/" + url.QueryEscape(gifURL))
	if err != nil {
		return gfyCheckURLResponse{}, err
	}
	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return gfyCheckURLResponse{}, err
	}
	err = getGfyError(bytes)
	if err != nil {
		return gfyCheckURLResponse{}, err
	}
	var rsp gfyCheckURLResponse
	err = json.Unmarshal(bytes, &rsp)
	if err != nil {
		return gfyCheckURLResponse{}, err
	}
	return rsp, nil
}

// GetOEmbedDataByGfyURL returns oembed information about the provided gfycat.com URL.
func GetOEmbedDataByGfyURL(gfyURL string, maxWidth, maxHeight int) (gfyOEmbed, error) {
	response, err := http.DefaultClient.Get("https://api.gfycat.com/v1/oembed?url=" + url.QueryEscape(gfyURL) + "&maxwidth=" + strconv.Itoa(maxWidth) + "&maxheight=" + strconv.Itoa(maxHeight))
	if err != nil {
		return gfyOEmbed{}, err
	}
	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return gfyOEmbed{}, err
	}
	err = getGfyError(bytes)
	if err != nil {
		return gfyOEmbed{}, err
	}
	var rsp gfyOEmbed
	err = json.Unmarshal(bytes, &rsp)
	if err != nil {
		return gfyOEmbed{}, err
	}
	return rsp, nil
}
