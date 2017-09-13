package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type Coubs struct {
	Coubs      []Coub `json:"coubs"`
	TotalPages int    `json:"total_pages"`
	Page       int    `json:"page"`
	PerPage    int    `json:"per_page"`
	Tag        string `json:"_"`
	CoubID     int    `json:"_"`
}

type Coub struct {
	ID        int    `json:"id"`
	Permalink string `json:"permalink"`
}

// var Tags = map[string]string{"сиськи": "boobs", "животные": "animals-pets", "смешные": "funny"}

const CoubUrl = "https://coub.com/view/"

func GetApiUrlByTag(tag string) string {
	if tag == "boobs" {
		return "https://coub.com/api/v2/timeline/tag/boobs?order_by=newest_popular&page="
	}
	if tag == "funny" {
		return "https://coub.com/api/v2/timeline/hot/funny/weekly?page="
	}
	if tag == "animals-pets" {
		return "https://coub.com/api/v2/timeline/hot/animals-pets/weekly?page="
	}
	return "https://coub.com/api/v2/timeline/tag/" + tag + "?order_by=newest_popular&page="
}

func (c *Coubs) GetCoubURL(id int, page int, tag string) string {

	if id <= len(c.Coubs)-1 && page == c.Page && tag == c.Tag {
		c.CoubID = id
		return CoubUrl + c.Coubs[id].Permalink
	}

	url := GetApiUrlByTag(tag) + strconv.Itoa(page)
	log.Println(url)
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, c)
	if err != nil {
		log.Println(err)
	}

	if id > len(c.Coubs)-1 {
		resp.Body.Close()
		id = 0
		page = page + 1
		if page > c.TotalPages {
			return "На данную тему больше видео нету 😢"
		}
		return c.GetCoubURL(id, page, tag)
	}
	c.CoubID = id
	c.Page = page
	c.Tag = tag
	return CoubUrl + c.Coubs[id].Permalink
}

func (c *Coubs) GetNext() string {
	return c.GetCoubURL(c.CoubID+1, c.Page, c.Tag)
}

func (c *Coubs) GetBoobs() string {
	log.Println("GetBoobs")
	return c.GetCoubURL(0, 1, "boobs")
}

func (c *Coubs) GetAnimal() string {
	log.Println("GetAnimal")
	return c.GetCoubURL(0, 1, "animals-pets")
}

func (c *Coubs) GetFunny() string {
	log.Println("GetFunny")
	return c.GetCoubURL(0, 1, "funny")
}
