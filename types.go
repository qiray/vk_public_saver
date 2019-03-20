package main

//Read https://vk.com/dev/attachments_w for more info about attachments
/*
6. Граффити (type = graffiti)
8. Заметка (type = note)
9. Контент приложения (type = app)
11. Вики-страница (type = page)
12. Альбом с фотографиями (type=album)
13. Список фотографий (type=photos_list)
14. Товар (type = market)
15. Подборка товаров (type = market_album)
16. Карточки (type = pretty_cards)
*/

//Some structs are generated by https://mholt.github.io/json-to-go/

//Attachment - common struct for vk attachment (photo, video, poll etc.)
type Attachment struct {
	Type  string `json:"type"`
	Photo struct {
		ID        int    `json:"id"`
		AlbumID   int    `json:"album_id"`
		OwnerID   int    `json:"owner_id"`
		UserID    int    `json:"user_id"`
		Text      string `json:"text"`
		Date      int    `json:"date"`
		AccessKey string `json:"access_key"`
	} `json:"photo"`
	PostedPhoto struct {
		ID       int    `json:"id"`
		OwnerID  int    `json:"owner_id"`
		Photo130 string `json:"photo_130"`
		Photo604 string `json:"photo_604"`
	} `json:"posted_photo"`
	Video struct {
		ID          int    `json:"id"`
		OwnerID     int    `json:"owner_id"`
		Title       string `json:"title"`
		Duration    int    `json:"duration"`
		Description string `json:"description"`
		Date        int    `json:"date"`
		Comments    int    `json:"comments"`
		Views       int    `json:"views"`
		Width       int    `json:"width"`
		Height      int    `json:"height"`
		AccessKey   string `json:"access_key"`
	} `json:"video"`
	Audio struct {
		ID         int    `json:"id"`
		OwnerID    int    `json:"owner_id"`
		Artist     string `json:"artist"`
		Title      string `json:"title"`
		URL        string `json:"url"`
		Duration   int    `json:"duration"`
		Date       int    `json:"date"`
		AlbumID    int    `json:"album_id"`
		IsHq       bool   `json:"is_hq"`
		TrackCode  string `json:"track_code"`
		IsExplicit bool   `json:"is_explicit"`
	} `json:"audio"`
	Doc struct {
		ID      int    `json:"id"`
		OwnerID int    `json:"owner_id"`
		Size    int    `json:"size"`
		Title   string `json:"title"`
		Date    int    `json:"date"`
		Type    int    `json:"type"`
		Ext     string `json:"ext"`
		URL     string `json:"url"`
	} `json:"doc"`
	Link struct {
		URL         string `json:"url"`
		Title       string `json:"title"`
		Caption     string `json:"caption"`
		Description string `json:"description"`
	} `json:"link"`
	Poll struct {
		ID       int    `json:"id"`
		OwnerID  int    `json:"owner_id"`
		Question string `json:"question"`
		Votes    int    `json:"votes"`
	} `json:"poll"`
}

//Post - type for vk post structure
type Post struct {
	ID          int          `json:"id"`
	FromID      int          `json:"from_id"`
	OwnerID     int          `json:"owner_id"`
	SignerID    int          `json:"signer_id"`
	Date        int          `json:"date"`
	MarkedAsAds int          `json:"marked_as_ads"`
	PostType    string       `json:"post_type"`
	Text        string       `json:"text"`
	IsPinned    int          `json:"is_pinned"`
	Attachments []Attachment `json:"attachments"`
	Comments    struct {
		Count int `json:"count"`
	} `json:"comments"`
	Likes struct {
		Count int `json:"count"`
	} `json:"likes"`
	Reposts struct {
		Count int `json:"count"`
	} `json:"reposts"`
	Views struct {
		Count int `json:"count"`
	} `json:"views"`
}

//PostsResponse - response struct
type PostsResponse struct {
	Response []struct {
		Count int    `json:"count"`
		Items []Post `json:"items"`
	} `json:"response"`
}

//AppSettings - struct for application settings
type AppSettings struct {
	AppID      string `json:"app_id"`
	APIVersion string `json:"api_version"`
	userdata   map[string]string
	token      string
}

//Version - struct for version info
type Version struct {
	Major int
	Minor int
	Build int
}
