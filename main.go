package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/hekmon/transmissionrpc/v2"
	"gopkg.in/yaml.v3"
)

// Config 定义配置结构体
type Config struct {
	Pattern  string `yaml:"pattern"`
	Tracker  string `yaml:"tracker"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

// loadConfig 从配置文件加载配置
func loadConfig(filename string) (*Config, error) {
	config := &Config{
		Host: "localhost",
		Port: 9091,
	}
	if filename == "" {
		return config, nil
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return config, nil
		}
		return nil, fmt.Errorf("读取配置文件错误：%v", err)
	}

	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("解析配置文件错误：%v", err)
	}

	return config, nil
}

func main() {
	// 定义命令行参数
	pattern := flag.String("pattern", "", "正则表达式，用于匹配要替换的tracker地址")
	newTracker := flag.String("tracker", "", "新的tracker地址")
	host := flag.String("host", "localhost", "Transmission RPC 主机地址")
	port := flag.Int("port", 9091, "Transmission RPC 端口")
	username := flag.String("username", "", "Transmission RPC 用户名")
	password := flag.String("password", "", "Transmission RPC 密码")
	configFile := flag.String("config", "", "配置文件路径")
	flag.Parse()

	// 从配置文件加载配置
	config, err := loadConfig(*configFile)
	if err != nil {
		log.Fatalf("加载配置文件失败：%v\n", err)
	}

	// 命令行参数优先级高于配置文件
	if *pattern != "" {
		config.Pattern = *pattern
	}
	if *newTracker != "" {
		config.Tracker = *newTracker
	}
	if *host != "localhost" {
		config.Host = *host
	}
	if *port != 9091 {
		config.Port = *port
	}
	if *username != "" {
		config.Username = *username
	}
	if *password != "" {
		config.Password = *password
	}

	// 验证必需参数
	if config.Pattern == "" || config.Tracker == "" {
		fmt.Println("请提供所有必需的参数：pattern 和 tracker（通过命令行参数或配置文件）")
		flag.Usage()
		os.Exit(1)
	}

	// 编译正则表达式
	reg, err := regexp.Compile(config.Pattern)
	if err != nil {
		log.Fatalf("正则表达式编译错误：%v\n", err)
	}

	// 创建 Transmission RPC 客户端
	conf := transmissionrpc.AdvancedConfig{
		HTTPS: true,
		Port:  uint16(config.Port),
	}

	client, err := transmissionrpc.New(config.Host, config.Username, config.Password, &conf)
	if err != nil {
		log.Fatalf("连接到 Transmission 失败：%v\n", err)
	}

	// 获取所有种子
	torrents, err := client.TorrentGet(context.Background(), []string{"id", "name", "trackers"}, nil)
	if err != nil {
		log.Fatalf("获取种子列表失败：%v\n", err)
	}

	successCount := 0
	failCount := 0

	// 处理每个种子
	for _, torrent := range torrents {
		modified := false
		for _, tracker := range torrent.Trackers {
			if reg.MatchString(tracker.Announce) {
				// 找到匹配的tracker，准备更新
				fmt.Printf("[%s] 找到匹配的tracker：%s\n", *torrent.Name, tracker.Announce)

				// 更新tracker
				err := client.TorrentSet(context.Background(), transmissionrpc.TorrentSetPayload{
					IDs:           []int64{*torrent.ID},
					TrackerRemove: []int64{tracker.ID},
					TrackerAdd:    []string{config.Tracker},
				})

				if err != nil {
					fmt.Printf("[%s] 更新tracker失败：%v\n", *torrent.Name, err)
					failCount++
				} else {
					fmt.Printf("[%s] 已替换为新的tracker：%s\n", *torrent.Name, *newTracker)
					successCount++
				}
				modified = true
				break
			}
		}

		if !modified {
			fmt.Printf("[%s] 未找到匹配的tracker地址\n", *torrent.Name)
			failCount++
		}
	}

	// 输出处理结果统计
	fmt.Printf("\n处理完成！成功：%d，失败：%d\n", successCount, failCount)
}
