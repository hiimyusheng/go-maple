package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

// Synology NAS 設定
const (
	NAS_URL   = "https://nas.cisyy.cc"
	USERNAME  = "mapler"
	PASSWORD  = "*8syjhmUY&b4"
	FILE_PATH = "/MapleStory/file_tree.json" // 要下載的檔案路徑
	VERSION   = "268"
)

// AuthResponse 解析登入回應的 JSON
type AuthResponse struct {
	Success bool `json:"success"`
}

type FileTree map[string]struct {
	New    string `json:"new"`
	Update string `json:"update"`
}

func main() {
	// 步驟 1: 登入並獲取 Cookie
	cookie, err := login()
	if err != nil {
		fmt.Println("登入失敗:", err)
		return
	}
	// fmt.Println("登入成功，開始下載檔案...")

	// 步驟 2: 下載檔案
	if err := downloadFile(cookie, FILE_PATH); err != nil {
		fmt.Println("下載失敗:", err)
		return
	}
	// 步驟 3: 讀取並解析 JSON 檔案
	file, err := os.Open("updater/" + filepath.Base(FILE_PATH))
	if err != nil {
		fmt.Println("無法開啟檔案:", err)
		return
	}
	defer file.Close()

	var data []FileTree
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		fmt.Println("解析 JSON 失敗:", err)
		return
	}

	// 從終端接收輸入
	var input string
	fmt.Println("目前預載版本為V" + VERSION + "，全新下載請輸入new，更新請輸入update")
	fmt.Print("請輸入: ")
	fmt.Scanln(&input)

	if (input != "new") && (input != "update") {
		fmt.Println("輸入錯誤")
		return
	}

	if input == "new" {
		var maxVersion string
		for _, item := range data {
			for key := range item {
				if key > maxVersion {
					maxVersion = key
				}
			}
		}
		for _, value := range data {
			if value[maxVersion].New != "" {
				fmt.Println("開始下載最新版本V" + maxVersion)
				// fmt.Println(value[maxVersion].New)
				downloadPath := "/MapleStory/" + maxVersion + "/" + value[maxVersion].New
				if err := downloadFile(cookie, downloadPath); err != nil {
					fmt.Println("下載失敗:", err)
					return
				}
				fmt.Println(value[maxVersion].New + " 下載完成")
			}
		}
		return
	}
	if input == "update" {
		for _, item := range data {
			for version, value := range item {
				if version > VERSION && value.Update != "" {
					fmt.Println("開始下載更新版本V" + version)
					// fmt.Println(value.Update)
					downloadPath := "/MapleStory/" + version + "/" + item[version].Update
					if err := downloadFile(cookie, downloadPath); err != nil {
						fmt.Println("下載失敗:", err)
						return
					}
					fmt.Println(item[version].Update + " 下載完成")
				}
			}
		}
		return
	}
}

// login 進行登入並獲取 Cookie
func login() (string, error) {
	loginURL := fmt.Sprintf("%s/webapi/auth.cgi?api=SYNO.API.Auth&version=3&method=login&account=%s&passwd=%s&session=FileStation&format=cookie", NAS_URL, url.QueryEscape(USERNAME), url.QueryEscape(PASSWORD))

	// fmt.Println("loginURL:", loginURL)
	resp, err := http.Get(loginURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var authResp AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return "", err
	}

	if !authResp.Success {
		return "", fmt.Errorf("登入失敗")
	}

	// 取得 Set-Cookie 標頭
	cookies := resp.Header.Values("Set-Cookie")
	if len(cookies) == 0 {
		return "", fmt.Errorf("未收到 Cookie")
	}

	return cookies[0], nil
}

// downloadFile 下載檔案
func downloadFile(cookie, file_path string) error {
	downloadURL := fmt.Sprintf("%s/webapi/entry.cgi?api=SYNO.FileStation.Download&version=2&method=download&path=[\"%s\"]&mode=download", NAS_URL, url.QueryEscape(file_path))

	// fmt.Println("downloadURL:", downloadURL)
	// 建立 HTTP 請求
	req, err := http.NewRequest("GET", downloadURL, nil)
	if err != nil {
		return err
	}

	// 設定 Cookie
	req.Header.Set("Cookie", cookie)

	// 發送請求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 檢查回應狀態
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("下載失敗，HTTP 狀態碼: %d", resp.StatusCode)
	}
	filename := filepath.Base(file_path)
	// 儲存到本地檔案
	if err := os.MkdirAll("updater", os.ModePerm); err != nil {
		return err
	}
	file, err := os.Create("updater/" + filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	return err
}
