package watermeter

import (
	"encoding/json"
	"io"
	"net/http"
)

type Watermeter struct {
	addr     string
	Incoming chan *Telegram
}

type Telegram struct {
	Info *Info
	Data *Data
}

type Info struct {
	ProductName     string `json:"product_name"`
	ProductType     string `json:"product_type"`
	Serial          string `json:"serial"`
	FirmwareVersion string `json:"firmware_version"`
	ApiVersion      string `json:"api_version"`
}

type Data struct {
	WifiSsid           string  `json:"wifi_ssid"`
	WifiStrength       int     `json:"wifi_strength"`
	TotalLiterM3       float64 `json:"total_liter_m3"`
	ActiveLiterLpm     int     `json:"active_liter_lpm"`
	TotalLiterOffsetM3 int     `json:"total_liter_offset_m3"`
}

func New(addr string) (*Watermeter, error) {
	return &Watermeter{addr: addr}, nil
}

func (w *Watermeter) Start() {
	go w.run()
}

func (w *Watermeter) run() {
	info, err := w.readInfo()
	if err != nil {
		return
	}
	defer close(w.Incoming)

	for {
		data, err := w.readData()
		if err != nil {
			return
		}
		w.Incoming <- &Telegram{info, data}
	}
}

func (w *Watermeter) readInfo() (*Info, error) {
	resp, err := http.Get(w.addr + "/api")
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	_ = resp.Body.Close()
	if err != nil {
		return nil, err
	}
	var data Info
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (w *Watermeter) readData() (*Data, error) {
	resp, err := http.Get(w.addr + "/api/v1/data")
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	_ = resp.Body.Close()
	if err != nil {
		return nil, err
	}
	var data Data
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}
