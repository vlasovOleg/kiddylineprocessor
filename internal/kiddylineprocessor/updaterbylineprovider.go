package kiddylineprocessor

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
	"time"
)

func (kp *Kiddylineprocessor) requestLine(sport string) (float32, error) {
	response := struct {
		Lines struct {
			CoefficientBASEBALL string `json:"BASEBALL"`
			CoefficientFOOTBALL string `json:"FOOTBALL"`
			CoefficientSOCCER   string `json:"SOCCER"`
		} `json:"lines"`
	}{}

	resp, err := kp.httpClient.Get(kp.config.LinesProvider.Address + sport)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
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
	return 0, errors.New("wrong data")
}

func (kp *Kiddylineprocessor) updateFromProvider(ctx context.Context, wg *sync.WaitGroup, sportName string, syncInterval time.Duration, storeFunc func(float32) error) {
	errHandler := func(msg string, err error) {
		kp.errors.Mutex.Lock()
		kp.errors.data[sportName] = err.Error()
		kp.errors.Mutex.Unlock()
		kp.loger.Error("Kiddylineprocessor : updateFromProvider : ", sportName, " : ", msg, " : ", err.Error())
	}

	update := func() {
		c, err := kp.requestLine(sportName)
		if err != nil {
			errHandler("requestLine", err)
			return
		}
		err = storeFunc(c)
		if err != nil {
			errHandler("storeFunc", err)
			return
		}
		if _, ok := kp.errors.data[sportName]; ok {
			kp.errors.Mutex.Lock()
			delete(kp.errors.data, sportName)
			kp.errors.Mutex.Unlock()
		}
		kp.loger.Debug("Kiddylineprocessor : updateFromProvider : ", sportName, " : new coefficient : ", c)
	}

	for {
		time.Sleep(syncInterval)
		select {
		case <-ctx.Done():
			kp.loger.Info("Kiddylineprocessor : updateFromProvider : ", sportName, " : stopped")
			wg.Done()
			return
		default:
			update()
		}
	}
}
