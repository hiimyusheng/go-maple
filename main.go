package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/schollz/progressbar/v3"
)

var (
	URL       = ""
	USERNAME  = ""
	PASSWORD  = ""
	FILE_PATH = "/MapleStory/file_tree.json"
)

// AuthResponse 解析登入回應的 JSON
type AuthResponse struct {
	Success bool `json:"success"`
}

type FileTree struct {
	New    map[string]string `json:"new"`
	Update map[string]string `json:"update"`
}

func init() {
	// 檢查必要的變數是否被正確注入
	if URL == "" || USERNAME == "" || PASSWORD == "" {
		fmt.Println("錯誤: 未能正確注入必要的環境變數")
		os.Exit(1)
	}
}

func main() {

	fmt.Println(`
                                                  
                 @@@@@@@@@@                       
              @@@@%#######%@@@@@                  
             @@%**++==+++**####%@                 
           @%%%%%%%#*++++++****#%@@               
         @%#+++*+++*###*++++****#%%@              
         @%*+####*++++***+++******#%@@            
        @@%*+##**#*++******+++******#%@@          
     @@@#%%*+#***#*++*****+++++++*****#%@         
   @@#*+=+#**#####*++**+++++==++++++****%@@       
  @%*+=--+#****+*#*++++++========+++++**#%@@      
  @%*+===-----====================++++++**#@      
  @@*++========================++++++++**#%@      
  @@#*++++===============+++++++++++****#%@       
   @@#**+++++++++++++++++++++++*******##%         
     @%%###************************###%@          
      @@#*#%%%%%%%%%%%%%%%%%%%%%%%%##%%           
     @%*+==++****++++++++++++*****+==+#@          
     @%+====*****+--------===========+#@          
     @%+--+*#****+=-::::::----=======+#@@%%       
     @#=:=*#*+++*++=:.:::::::---=====+#@%###      
      @%==*#*+++*++=:..::::::::---==+#%%#######%% 
      @%*++*#*+++++=:...::::::::--=+#%%########## 
         %%%%#*+***+=-:::::----*##%%%%#########   
           %%@@%%@@%%%%%%%%%%%%%%%%#########      
             ###########################          
                                                  `)

	fmt.Println("╔╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╗")
	fmt.Println("╠╬╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╬╣")
	fmt.Println("╠╣ __  __             _           _                         ╠╣")
	fmt.Println("╠╣|  \\/  |           | |         | |                        ╠╣")
	fmt.Println("╠╣| \\  / | __ _ _ __ | | ___  ___| |_ ___  _ __ _   _       ╠╣")
	fmt.Println("╠╣| |\\/| |/ _` | '_ \\| |/ _ \\/ __| __/ _ \\| '__| | | |      ╠╣")
	fmt.Println("╠╣| |  | | (_| | |_) | |  __/\\__ \\ || (_) | |  | |_| |      ╠╣")
	fmt.Println("╠╣|_|  |_|\\__,_| .__/|_|\\___||___/\\__\\___/|_|   \\__, |      ╠╣")
	fmt.Println("╠╣             | |                               __/ |      ╠╣")
	fmt.Println("╠╣ _____       |_|            _                 |___/       ╠╣")
	fmt.Println("╠╣|  __ \\                    | |               | |          ╠╣")
	fmt.Println("╠╣| |  | | _____      ___ __ | | ___   __ _  __| | ___ _ __ ╠╣")
	fmt.Println("╠╣| |  | |/ _ \\ \\ /\\ / / '_ \\| |/ _ \\ / _` |/ _` |/ _ \\ '__|╠╣")
	fmt.Println("╠╣| |__| | (_) \\ V  V /| | | | | (_) | (_| | (_| |  __/ |   ╠╣")
	fmt.Println("╠╣|_____/ \\___/ \\_/\\_/ |_| |_|_|\\___/ \\__,_|\\__,_|\\___|_|   ╠╣")
	fmt.Println("╠╬╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╬╣")
	fmt.Println("╚╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╩╝")

	// 步驟 1: 登入並獲取 Cookie
	cookie, err := login()
	if err != nil {
		fmt.Println("登入失敗:", err)
		fmt.Println("輸入任意鍵以結束程式...")
		fmt.Scanln()
		return
	}
	// fmt.Println("登入成功，開始下載檔案...")

	// 步驟 2: 下載檔案
	if err := downloadFile(cookie, FILE_PATH); err != nil {
		fmt.Println("下載失敗:", err)
		fmt.Println("輸入任意鍵以結束程式...")
		fmt.Scanln()
		return
	}
	// 步驟 3: 讀取並解析 JSON 檔案
	file, err := os.Open("updater/" + filepath.Base(FILE_PATH))
	if err != nil {
		fmt.Println("無法開啟檔案:", err)
		fmt.Println("輸入任意鍵以結束程式...")
		fmt.Scanln()
		return
	}
	defer file.Close()

	var data FileTree
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		fmt.Println("解析 JSON 失敗:", err)
		fmt.Println("輸入任意鍵以結束程式...")
		fmt.Scanln()
		return
	}

	// 從終端接收輸入
	var input string
	var version string
	for key := range data.New {
		version = key
		break
	}
	fmt.Println("目前最新版本為 V" + version + "，全新下載請輸入new，更新請輸入update")
	fmt.Print("請輸入: ")
	fmt.Scanln(&input)

	if (input != "new") && (input != "update") {
		fmt.Println("輸入錯誤")
		fmt.Println("輸入任意鍵以結束程式...")
		fmt.Scanln()
		return
	}

	if input == "new" {
		fmt.Println("開始下載最新版本 V" + version)
		if err := downloadFileFromOfficial(data.New[version]); err != nil {
			fmt.Println("下載失敗:", err)
			fmt.Println("輸入任意鍵以結束程式...")
			fmt.Scanln()
			return
		}
		fmt.Println("V" + version + " 安裝檔下載完成")
	}

	if input == "update" {
		fmt.Println("請輸入目前安裝的版本，將開始下載更新檔")
		fmt.Print("請輸入: ")
		fmt.Scanln(&input)
		for key, item := range data.Update {
			if key <= version && key > input {
				if item != "" {
					fmt.Println("開始下載更新版本 V" + key)
					downloadPath := "/MapleStory/" + key + "/" + item
					if err := downloadFile(cookie, downloadPath); err != nil {
						fmt.Println("下載失敗:", err)
						fmt.Println("輸入任意鍵以結束程式...")
						fmt.Scanln()
						return
					}
					fmt.Println("V" + key + " 更新檔下載完成")
				}
			}
		}
	}
	fmt.Println(`
--------------------------------------------------
🚀 Powered by hiimyusheng
📧 scott@cisyy.cc
🌍 GitHub: github.com/hiimyusheng
--------------------------------------------------`)
	fmt.Println("輸入任意鍵結束程式...")
	fmt.Scanln()
	return
}

// login 進行登入並獲取 Cookie
func login() (string, error) {
	loginURL := fmt.Sprintf("%s/webapi/auth.cgi?api=SYNO.API.Auth&version=3&method=login&account=%s&passwd=%s&session=FileStation&format=cookie", URL, url.QueryEscape(USERNAME), url.QueryEscape(PASSWORD))

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

func downloadFileFromOfficial(downloadUrl string) error {
	// 發送請求
	resp, err := http.Get(downloadUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	// 檢查回應狀態
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("下載失敗，HTTP 狀態碼: %d", resp.StatusCode)
	}

	// 取得檔案大小
	contentLength := resp.Header.Get("Content-Length")
	totalSize, err := strconv.Atoi(contentLength)
	if err != nil {
		return fmt.Errorf("無法取得檔案大小: %v", err)
	}

	filename := filepath.Base(downloadUrl)
	// 儲存到本地檔案
	if err := os.MkdirAll("updater", os.ModePerm); err != nil {
		return err
	}
	file, err := os.Create("updater/" + filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// 建立進度條
	bar := progressbar.NewOptions(totalSize,
		progressbar.OptionSetDescription("下載中..."),
		progressbar.OptionShowCount(),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(40),
		progressbar.OptionThrottle(65*time.Millisecond),
		progressbar.OptionShowIts(),
		progressbar.OptionOnCompletion(func() {
			fmt.Println("\n下載完成")
		}),
	)

	// 使用 TeeReader 來更新進度條
	_, err = io.Copy(io.MultiWriter(file, bar), resp.Body)
	return err
}

// downloadFile 下載檔案
func downloadFile(cookie, file_path string) error {
	downloadURL := fmt.Sprintf("%s/webapi/entry.cgi?api=SYNO.FileStation.Download&version=2&method=download&path=[\"%s\"]&mode=download", URL, url.QueryEscape(file_path))

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

	// 取得檔案大小
	contentLength := resp.Header.Get("Content-Length")
	totalSize, err := strconv.Atoi(contentLength)
	if err != nil {
		return fmt.Errorf("無法取得檔案大小: %v", err)
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

	// 建立進度條
	bar := progressbar.NewOptions(totalSize,
		progressbar.OptionSetDescription("下載中..."),
		progressbar.OptionShowCount(),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(40),
		progressbar.OptionThrottle(65*time.Millisecond),
		progressbar.OptionShowIts(),
		progressbar.OptionOnCompletion(func() {
			fmt.Println("\n下載完成")
		}),
	)

	// 使用 TeeReader 來更新進度條
	_, err = io.Copy(io.MultiWriter(file, bar), resp.Body)
	return err
}
