package local_storage

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
)

type LinksFormRequest struct {
	Id         uint64
	Urls       []string
	UrlsStatus map[string]bool
}

type TempLinksFormRequest struct {
	mu       sync.Mutex
	requests map[uint64]LinksFormRequest
	tempId   uint64
}

func NewLocalStorage() *TempLinksFormRequest {
	return &TempLinksFormRequest{
		requests: make(map[uint64]LinksFormRequest),
		tempId:   1,
	}
}

func (TL *TempLinksFormRequest) AddUrl(urls []string) uint64 {
	TL.mu.Lock()
	defer TL.mu.Unlock()

	id := TL.tempId
	TL.tempId++

	TL.requests[id] = LinksFormRequest{
		Id:         id,
		Urls:       urls,
		UrlsStatus: make(map[string]bool),
	}
	return id
}

func (TL *TempLinksFormRequest) GetUrls(id uint64) (LinksFormRequest, bool) {
	TL.mu.Lock()
	defer TL.mu.Unlock()

	urls, ok := TL.requests[id]
	return urls, ok
}

func (TL *TempLinksFormRequest) UpdateStatus(status bool, id uint64, url string) {
	TL.mu.Lock()
	defer TL.mu.Unlock()

	if urls, ok := TL.requests[id]; ok {
		urls.UrlsStatus[url] = status
		TL.requests[id] = urls
	}
}

func (TL *TempLinksFormRequest) UpdateAllStatuses(id uint64, statuses map[string]bool) {
	TL.mu.Lock()
	defer TL.mu.Unlock()

	if urls, ok := TL.requests[id]; ok {
		for u, s := range statuses {
			urls.UrlsStatus[u] = s
		}
		TL.requests[id] = urls
	}
}

type persistedState struct {
	TempId   uint64
	Requests map[uint64]LinksFormRequest
}

func (TL *TempLinksFormRequest) SaveToDisk(path string) error {
	TL.mu.Lock()
	defer TL.mu.Unlock()

	data := persistedState{
		TempId:   TL.tempId,
		Requests: TL.requests,
	}
	buf, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, buf, 0666)
}

func (TL *TempLinksFormRequest) LoadFromDisk(path string) error {
	TL.mu.Lock()
	defer TL.mu.Unlock()

	b, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}
	var data persistedState
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}
	if data.Requests == nil {
		data.Requests = make(map[uint64]LinksFormRequest)
	}
	TL.requests = data.Requests
	if data.TempId == 0 {
		data.TempId = 1
	}
	TL.tempId = data.TempId
	return nil
}
