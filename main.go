package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
)

const debug = false

type EmojiData struct {
	SupportSize []struct{ Width, Height int } `json:"supportSize"`
	Name        string                        `json:"name"`
	Mark        string                        `json:"mark"`
	Imgs        []struct{ ID, Name string }   `json:"imgs"`
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("输入链接/ID/本地文件: ")
		if !scanner.Scan() {
			break
		}
		input := strings.Trim(strings.TrimSpace(scanner.Text()), `"'`)
		if input == "" {
			break
		}

		if isLocalFile(input) {
			processLocal(input)
		} else if id := extractID(input); id != "" {
			if strings.HasPrefix(input, "http") {
				fmt.Println("提取ID:", id)
			}
			processRemote(id)
		} else {
			fmt.Println("输入无效")
			continue
		}
		fmt.Println("完成\n")
	}
}

func isLocalFile(path string) bool {
	ext := strings.ToLower(path)
	return strings.HasSuffix(ext, ".json") || strings.HasSuffix(ext, ".jtmp")
}

func extractID(input string) string {
	if !strings.HasPrefix(input, "http") {
		return input
	}
	if match := regexp.MustCompile(`[?&]id=(\d+)`).FindStringSubmatch(input); len(match) > 1 {
		return match[1]
	}
	return ""
}

func processRemote(id string) {
	url := fmt.Sprintf("https://gxh.vip.qq.com/club/item/parcel/%s/%s.json", id[len(id)-1:], id)
	if debug {
		fmt.Println("JSON:", url)
	}

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("请求失败:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var data EmojiData
	if json.Unmarshal(body, &data) != nil {
		fmt.Println("解析失败")
		return
	}

	dir := fmt.Sprintf("[ID%s] %s", id, data.Name)
	os.Mkdir(dir, os.ModePerm)
	if debug {
		os.WriteFile(fmt.Sprintf("%s/%s.json", dir, id), body, 0644)
	}
	download(data, dir)
}

func processLocal(path string) {
	body, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("读取失败:", err)
		return
	}

	var data EmojiData
	if json.Unmarshal(body, &data) != nil {
		fmt.Println("解析失败")
		return
	}

	// 从文件名提取ID
	id := regexp.MustCompile(`(\d+)`).FindString(path)
	dir := data.Name
	if id != "" {
		dir = fmt.Sprintf("[ID%s] %s", id, data.Name)
	}
	os.Mkdir(dir, os.ModePerm)
	download(data, dir)
}

func download(data EmojiData, dir string) {
	h := data.SupportSize[0].Height
	fmt.Println("名称:", data.Name)
	fmt.Printf("总数: %d，尺寸: %dx%d\n", len(data.Imgs), data.SupportSize[0].Width, h)
	if data.Mark != "" {
		fmt.Println("备注:", data.Mark)
	}

	for i, img := range data.Imgs {
		url := fmt.Sprintf("https://gxh.vip.qq.com/club/item/parcel/item/%s/%s/raw%d.gif", img.ID[:2], img.ID, h)
		path := fmt.Sprintf("%s/%d_%s.gif", dir, i+1, img.Name)
		if debug {
			fmt.Println("下载:", url)
		}
		if resp, err := http.Get(url); err == nil {
			if file, err := os.Create(path); err == nil {
				io.Copy(file, resp.Body)
				file.Close()
				fmt.Println("完成:", path)
			}
			resp.Body.Close()
		}
	}
}
