package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// 生成密码哈希
	passwords := []string{"admin123", "user123"}

	for _, pwd := range passwords {
		hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
		if err != nil {
			fmt.Printf("生成 %s 的哈希失败: %v\n", pwd, err)
			continue
		}
		fmt.Printf("密码: %s\n哈希: %s\n\n", pwd, string(hash))
	}
}
