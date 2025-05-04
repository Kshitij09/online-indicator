package main

import (
	"fmt"
	"net/http"
)

func Ping(accountId, sessionId, baseUrl string) error {
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/ping/%s", baseUrl, accountId), nil)
	if err != nil {
		return err
	}
	req.Header.Set("X-Session-Token", sessionId)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
