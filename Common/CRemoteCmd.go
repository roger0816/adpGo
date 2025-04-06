package Common

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/crypto/ssh"
)

func runSh(ip string, scriptPath string) {
	// 動態取得使用者家目錄
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("取得家目錄失敗:", err)
		return
	}

	// 組合出私鑰的路徑，例如 ~/.ssh/id_rsa
	keyPath := filepath.Join(homeDir, ".ssh", "id_rsa")

	// 使用 os.ReadFile 讀取私鑰檔案
	key, err := os.ReadFile(keyPath)
	if err != nil {
		fmt.Println("讀取私鑰失敗:", err)
		return
	}

	// 解析私鑰
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		fmt.Println("解析私鑰失敗:", err)
		return
	}

	// 設定 SSH 配置，使用 SSH 金鑰認證
	config := &ssh.ClientConfig{
		User: "your-username", // 替換成 VPS 上的使用者名稱
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // 生產環境中請使用正確的 host key 驗證
		Timeout:         5 * time.Second,
	}

	// 指定 VPS 的 IP 和 SSH 端口（預設為 22），例如 "192.168.1.100:22"
	client, err := ssh.Dial("tcp", ip, config)
	if err != nil {
		fmt.Println("SSH 連線失敗:", err)
		return
	}
	defer client.Close()

	// 建立 SSH session
	session, err := client.NewSession()
	if err != nil {
		fmt.Println("創建 session 失敗:", err)
		return
	}
	defer session.Close()

	// 使用傳入的 scriptPath 組成執行命令，例如："/bin/sh /path/to/your_script.sh"
	cmd := fmt.Sprintf("/bin/sh %s", scriptPath)
	err = session.Run(cmd)
	if err != nil {
		fmt.Println("執行命令失敗:", err)
		return
	}

	fmt.Println("腳本執行成功")
}
