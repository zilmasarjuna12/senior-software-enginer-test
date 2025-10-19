package wazuh

import (
	"fmt"
)

func (w *Wazuh) GetRules(queryString string) ([]byte, error) {
	err := w.authenticate()
	if err != nil {
		return nil, err
	}

	resp, err := w.Client.R().SetQueryString(queryString).Get("/rules")
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("GetRules failed: %s", resp.String())
	}

	return resp.Body(), nil
}

func (w *Wazuh) GetRulesFiles(queryString string) ([]byte, error) {
	err := w.authenticate()
	if err != nil {
		return nil, err
	}

	resp, err := w.Client.R().SetQueryString(queryString).Get("/rules/files")
	if err != nil {
		return nil, err
	}

	return resp.Body(), nil
}
