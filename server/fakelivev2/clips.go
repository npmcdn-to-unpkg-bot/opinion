package fakelive2

type videoType int
const (
	local videoType =   iota
	vod
	live
	embedclip
	autoencode
)

type Clips struct {
	Ad_first_clip               int64   `json:"ad_first_clip,string"`
	Ad_insertion_offset         int64   `json:"ad_insertion_offset,string"`
	Ad_insertion_pattern        int64   `json:"ad_insertion_pattern,string"`
	Ad_link                     string  `json:"ad_link"`
	Ad_type                     int64   `json:"ad_type,string"`
	Allow_comments              int64   `json:"allow_comments,string"`
	Aspect                      float64 `json:"aspect,string"`
	Clicks                      int64   `json:"clicks,string"`
	Date                        int64   `json:"date,string"`
	Date_lastmod                int64   `json:"date_lastmod,string"`
	Description                 string  `json:"description"`
	Description_seo             string  `json:"description_seo"`
	Dislikes                    int64   `json:"dislikes,string"`
	Downloadable                int64   `json:"downloadable,string"`
	Downloadable_condition      int64   `json:"downloadable_condition,string"`
	Downloadable_xfiles         string  `json:"downloadable_xfiles"`
	Duration                    int64   `json:"duration,string"`
	Id                          int64   `json:"id,string"`
	Id_import                   string  `json:"id_import"`
	Id_user                     int64   `json:"id_user,string"`
	Img_icon                    string  `json:"img_icon"`
	Img_poster                  string  `json:"img_poster"`
	Img_social                  string  `json:"img_social"`
	Img_thumbnail               string  `json:"img_thumbnail"`
	Interactivity_randomization int64   `json:"interactivity_randomization,string"`
	Interactivity_spacing       int64   `json:"interactivity_spacing,string"`
	Interactivity_start_delay   int64   `json:"interactivity_start_delay,string"`
	Interactivity_timing        int64   `json:"interactivity_timing,string"`
	Is_3d                       int64   `json:"is_3d,string"`
	Is_ad                       int64   `json:"is_ad,string"`
	Is_featured                 int64   `json:"is_featured,string"`
	Is_indexable                int64   `json:"is_indexable,string"`
	Is_searchable               int64   `json:"is_searchable,string"`
	Is_skippable                int64   `json:"is_skippable,string"`
	Is_skippable_after          int64   `json:"is_skippable_after,string"`
	Is_visitable                int64   `json:"is_visitable,string"`
	Likes                       int64   `json:"likes,string"`
	Privacy                     int64   `json:"privacy,string"`
	Privacy_access_level        int64   `json:"privacy_access_level,string"`
	Socialize                   int64   `json:"socialize,string"`
	Sprite_img                  string  `json:"sprite_img"`
	Sprite_vtt                  string  `json:"sprite_vtt"`
	Status                      int64   `json:"status,string"`
	Status_moderation           int64   `json:"status_moderation,string"`
	Store_on_sale               int64   `json:"store_on_sale,string"`
	Store_play_trailer          int64   `json:"store_play_trailer,string"`
	Tags                        string  `json:"tags"`
	Title                       string  `json:"title"`
	Title_url                   string  `json:"title_url"`
	Type                        int64   `json:"type,string"`
	Views                       int64   `json:"views,string"`
	Views_complete              int64   `json:"views_complete,string"`
	Views_embed                 int64   `json:"views_embed,string"`
	Views_page                  int64   `json:"views_page,string"`
}

type Clip_files struct {
	Embed_flash            string `json:"embed_flash"`
	Embed_html5            string `json:"embed_html5"`
	Encoding_source        string `json:"encoding_source"`
	Id                     int64  `json:"id,string"`
	Id_clip                int64  `json:"id_clip,string"`
	Id_quality             int64  `json:"id_quality,string"`
	Live_flash             string `json:"live_flash"`
	Live_html5_dash        string `json:"live_html5_dash"`
	Live_ios               string `json:"live_ios"`
	Live_ms                string `json:"live_ms"`
	Live_rtsp              string `json:"live_rtsp"`
	Vod_flash              string `json:"vod_flash"`
	Vod_flash_trailer      string `json:"vod_flash_trailer"`
	Vod_html5_h264         string `json:"vod_html5_h264"`
	Vod_html5_h264_trailer string `json:"vod_html5_h264_trailer"`
	Vod_html5_webm         string `json:"vod_html5_webm"`
	Vod_html5_webm_trailer string `json:"vod_html5_webm_trailer"`
}
