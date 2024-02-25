package util

import (
	"bytes"
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"github.com/denisbrodbeck/machineid"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"vocab-builder/pkg/server/conf"
)

func getLicenseServerAddress(LicenseServerAddressPublicationLink string) string {
	response, err := http.Get(LicenseServerAddressPublicationLink)
	if err != nil {
		fmt.Println("HTTP request failed:", err)
		os.Exit(1)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		fmt.Println("HTTP request returned status code:", response.StatusCode)
		os.Exit(1)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Failed to read response body:", err)
		os.Exit(1)
	}
	res := strings.Trim(string(body), "\n")
	return res
}

func getActivationCode(protectedMachineCode string, activationCode string) string {
	url := getLicenseServerAddress(conf.Cfg.LicenseServerAddressPublicationLink)

	data := map[string]string{
		"protected_machine_code": protectedMachineCode,
		"activation_code":        activationCode,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
	}

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println(err)
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
	}

	type Response struct {
		Code    int         `json:"code"`
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
	}
	var resp Response

	err = json.Unmarshal(body, &resp)
	if err != nil {
		log.Fatal(err)
	}
	if resp.Code != 0 {
		log.Fatalf("激活失败: %s\n", resp.Message)
	}
	return resp.Data.(string)
}

func generateSecretKey(s string) string {
	secretKey := ""
	for i := 0; i < 10; i++ {
		hash := sha512.New()
		hash.Write([]byte(s))
		hashedBytes := hash.Sum(nil)
		s = fmt.Sprintf("%x", hashedBytes)
		secretKey = s
	}
	return secretKey
}

func Activate() bool {
	const SECRET = "vocab-builder-machine-secret-key-do-not-change-it"
	protectedMachineID, err := machineid.ProtectedID(SECRET)
	if err != nil {
		log.Fatalf("获取机器码失败: %v", err)
	}

	secretKey := getActivationCode(protectedMachineID, conf.Cfg.Machine.ActivationCode)

	if secretKey != generateSecretKey(protectedMachineID) {
		log.Fatalf("尚未激活，您的机器码是: %s\n", protectedMachineID)
		return false
	}
	return true
}
