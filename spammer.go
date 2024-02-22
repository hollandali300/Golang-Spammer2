package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
)

type Config struct {
	Token      string `json:"token"`
	GuildID    string `json:"guild_id"`
	VanityCode string `json:"vanity_code"`
}

func main() {
	fmt.Println("Discord URL Spammer is starting...")

	config, err := readConfig("config.json")
	if err != nil {
		fmt.Println("Error reading config file:", err)
		return
	}

	file, err := os.Open("tokens.txt")
	if err != nil {
		fmt.Println("Error opening tokens file:", err)
		return
	}
	defer file.Close()

	var tokens []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		tokens = append(tokens, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading tokens file:", err)
		return
	}

	var wg sync.WaitGroup

	for _, t := range tokens {
		wg.Add(1)
		go func(token string) {
			defer wg.Done()
			for {
				getVanityURL(token, config)
			}
		}(t)
	}

	select {}
}

func readConfig(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	err = json.NewDecoder(file).Decode(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func getVanityURL(token string, config *Config) (string, error) {
	url := fmt.Sprintf("https://discord.com/api/v7/guilds/%s/vanity-url", config.GuildID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("error creating HTTP request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending HTTP request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var result struct {
			Code string `json:"code"`
		}
		err = json.NewDecoder(resp.Body).Decode(&result)
		if err != nil {
			return "", fmt.Errorf("error decoding response body: %v", err)
		}
		return result.Code, nil
	}

	return "", fmt.Errorf("failed to get Vanity URL: %d", resp.StatusCode)
}
