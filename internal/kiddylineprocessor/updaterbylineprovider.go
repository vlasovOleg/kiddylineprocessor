package kiddylineprocessor

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"strconv"
	"time"
)

func (kp *Kiddylineprocessor) updaterByLineProvider(sport string) (float32, error) {
	response := struct {
		Lines struct {
			// Coefficient string `json:"BASEBALL,FOOTBALL,SOCCER"`
			CoefficientBASEBALL string `json:"BASEBALL"`
			CoefficientFOOTBALL string `json:"FOOTBALL"`
			CoefficientSOCCER   string `json:"SOCCER"`
		} `json:"lines"`
	}{}

	resp, err := kp.httpClient.Get(kp.config.LinesProviderAddress + sport)
	if err != nil {
		return 0, err
	}

	if resp.StatusCode != 200 {
		responseText, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return 0, err
		}
		return 0, errors.New(string(responseText))
	}

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return 0, err
	}

	if response.Lines.CoefficientBASEBALL != "" {
		c, err := strconv.ParseFloat(response.Lines.CoefficientBASEBALL, 32)
		return float32(c), err
	}
	if response.Lines.CoefficientFOOTBALL != "" {
		c, err := strconv.ParseFloat(response.Lines.CoefficientFOOTBALL, 32)
		return float32(c), err
	}
	if response.Lines.CoefficientSOCCER != "" {
		c, err := strconv.ParseFloat(response.Lines.CoefficientSOCCER, 32)
		return float32(c), err
	}
	return 0, nil
}

func (kp *Kiddylineprocessor) updaterByLineProviderBaseball() {
	for {
		time.Sleep(kp.config.UpdateByProviderBaseball * time.Second)
		c, err := kp.updaterByLineProvider("baseball")
		if err != nil {
			kp.loger.Error("Kiddylineprocessor : updaterByProviderBaseball : updaterByProvider : ", err.Error())
			continue
		}
		kp.store.BaseballRepository().UpdateCoefficient(c)
		if err != nil {
			kp.loger.Error("Kiddylineprocessor : updaterByProviderBaseball : UpdateCoefficient : ", err.Error())
			continue
		}
		kp.loger.Debug("Kiddylineprocessor : updaterByProviderBaseball : new coefficient : ", c)
	}
}

func (kp *Kiddylineprocessor) updaterByLineProviderFootball() {
	for {
		time.Sleep(kp.config.UpdateByProviderFootball * time.Second)
		c, err := kp.updaterByLineProvider("football")
		if err != nil {
			kp.loger.Error("Kiddylineprocessor : updaterByProviderFootball : updaterByProvider : ", err.Error())
			continue
		}
		err = kp.store.FootballRepository().UpdateCoefficient(c)
		if err != nil {
			kp.loger.Error("Kiddylineprocessor : updaterByProviderFootball : UpdateCoefficient : ", err.Error())
			continue
		}
		kp.loger.Debug("Kiddylineprocessor : updaterByProviderFootball : new coefficient : ", c)
	}
}

func (kp *Kiddylineprocessor) updaterByLineProviderSoccer() {
	for {
		time.Sleep(kp.config.UpdateByProviderSoccer * time.Second)
		c, err := kp.updaterByLineProvider("soccer")
		if err != nil {
			kp.loger.Error("Kiddylineprocessor : updaterByProviderSoccer : updaterByProvider : ", err.Error())
			continue
		}
		err = kp.store.SoccerRepository().UpdateCoefficient(c)
		if err != nil {
			kp.loger.Error("Kiddylineprocessor : updaterByProviderSoccer : UpdateCoefficient : ", err.Error())
			continue
		}
		kp.loger.Debug("Kiddylineprocessor : updaterByProviderSoccer : new coefficient : ", c)
	}
}
