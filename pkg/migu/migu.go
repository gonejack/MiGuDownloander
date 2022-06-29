package migu

import (
	"encoding/json"
	"fmt"
	"gopkg.in/h2non/gentleman.v2"
	"log"
	"migugui/pkg/downloader"
	"net/url"
	"strings"
	"time"

	"github.com/spf13/cast"
	"github.com/wailsapp/wails"
)

var music []Music
var Type string
var suffix string
var gen = gentleman.New()

type MiGu struct {
	runtime *wails.Runtime
	Total   int
}

func get(keyword string, page int) {
	log.Printf("[*] Search %s Page [%d]", keyword, page)
	keywordencode := url.QueryEscape(keyword)
	var u = `http://pd.musicapp.migu.cn/MIGUM2.0/v1.0/content/search_all.do?&ua=Android_migu&version=5.0.1&text=%s&pageNo=%d&pageSize=10&searchSwitch={"song":1,"album":0,"singer":0,"tagSong":0,"mvSong":0,"songlist":0,"bestShow":1}`
	u = fmt.Sprintf(u, keywordencode, page)
	var migu = &MiGuResult{}
	rsp, err := gen.Get().URL(u).Do()
	if err == nil {
		err = rsp.JSON(migu)
	}
	if err != nil {
		log.Println("Get Music List Error : ", err)
		return
	}
	extract(migu)
	total := cast.ToInt(migu.SongResultData.TotalCount)
	log.Printf("Total : [%d] Current : [%d]", total, len(music))
	if len(music) < total {
		get(keyword, page+1)
	}
}
func extract(result *MiGuResult) {
	for _, m := range result.SongResultData.Result {
		u := `http://218.205.239.34/MIGUM2.0/v1.0/content/sub/listenSong.do?toneFlag=%s&netType=00&copyrightId=0&&contentId=%s&channel=0`
		username := ""
		for _, singer := range m.Singers {
			username += singer.Name + " "
		}
		albums := ""
		for _, album := range m.Albums {
			albums += album.Name + " "
		}
		u = fmt.Sprintf(u, Type, m.ContentID)
		music = append(music, Music{
			ID:     m.ID,
			Name:   m.Name,
			Album:  albums,
			Singer: username,
			URL:    u,
		})
	}
}

func (m *MiGu) WailsInit(runtime *wails.Runtime) error {
	m.runtime = runtime
	return nil
}
func (m *MiGu) Search(keyword string) []Music {
	music = []Music{}
	if strings.ToUpper(Type) == "HQ" {
		Type = "HQ&formatType=HQ&resourceType=2"
		suffix = "mp3"
	} else {
		Type = "SQ&formatType=SQ&resourceType=E"
		suffix = "flac"
	}
	get(keyword, 1)
	m.Total = len(music)
	return music
}
func (m *MiGu) GetResult() []Music {
	return music
}
func (m *MiGu) BatchDownload(request string) bool {
	path := m.runtime.Dialog.SelectDirectory()
	if path == "" {
		return false
	}
	var music = []Music{}
	if err := json.Unmarshal([]byte(request), &music); err != nil {
		log.Println("[!] Get Task Request Error : ", err)
		return false
	}
	m.Total = len(music)
	downloader.Count = 0
	go downloader.Start(path, 20)
	for _, v := range music {
		downloader.Push(downloader.Task{
			Name: v.Name + "_" + v.ID + "." + suffix,
			URL:  v.URL,
		})
	}
	time.Sleep(3 * time.Second)
	downloader.Lock.Wait()
	return true
}
func (m *MiGu) GetProgress() float32 {
	log.Println("Progress =>> ", downloader.Count, m.Total)
	return float32(downloader.Count) / float32(m.Total)
}
