package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
)

const debug = true

type EmojiData struct {
	SupportSize []struct {
		Width, Height int
	} `json:"supportSize"`
	Name string `json:"name"`
	Imgs []struct {
		ID, Name string
	} `json:"imgs"`
}

func main() {
	var input string
	fmt.Print("请输入 QQ 表情链接或 ID: ")
	fmt.Scanln(&input)

	// 从URL提取ID或直接使用输入的ID
	id := extractID(input)
	if id == "" {
		fmt.Println("无法获取 ID，退出。")
		waitExit()
		return
	}

	if strings.HasPrefix(input, "http") {
		fmt.Println("获取到表情 ID:", id)
	}
	downloadEmoji(id)
	fmt.Println("所有表情下载完成！")
	waitExit()
}

func waitExit() {
	fmt.Print("按回车键退出...")
	fmt.Scanln()
}

// extractID 从URL或直接返回ID
func extractID(input string) string {
	if !strings.HasPrefix(input, "http") {
		return input
	}

	// 直接从URL正则提取id参数
	re := regexp.MustCompile(`[?&]id=(\d+)`)
	if match := re.FindStringSubmatch(input); len(match) > 1 {
		return match[1]
	}
	return ""
}

// downloadEmoji 下载表情包
func downloadEmoji(id string) {
	jsonURL := fmt.Sprintf("https://gxh.vip.qq.com/club/item/parcel/%s/%s_android.json", id[len(id)-1:], id)
	if debug {
		fmt.Println("JSON 链接:", jsonURL)
	}

	resp, err := http.Get(jsonURL)
	if err != nil {
		fmt.Println("请求失败:", err)
		return
	}
	defer resp.Body.Close()

	var data EmojiData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		fmt.Println("解析失败:", err)
		return
	}

	saveDir := fmt.Sprintf("[ID%s] %s", id, data.Name)
	os.Mkdir(saveDir, os.ModePerm)

	height := data.SupportSize[0].Height
	fmt.Printf("表情总数: %d，尺寸: %dx%d\n", len(data.Imgs), data.SupportSize[0].Width, height)

	for idx, img := range data.Imgs {
		imgURL := fmt.Sprintf("https://gxh.vip.qq.com/club/item/parcel/item/%s/%s/raw%d.gif", img.ID[:2], img.ID, height)
		savePath := fmt.Sprintf("%s/%d_%s.gif", saveDir, idx+1, img.Name)

		if debug {
			fmt.Println("下载:", imgURL)
		}

		if resp, err := http.Get(imgURL); err == nil {
			if file, err := os.Create(savePath); err == nil {
				io.Copy(file, resp.Body)
				file.Close()
				fmt.Println("完成:", savePath)
			}
			resp.Body.Close()
		}
	}
}
