package gophy

type ImageData struct {
	URL    string `json:"url"`
	Width  string `json:"width"`
	Height string `json:"height"`
	Size   string `json:"size"`
	Frames string `json:"frames"`
}

type Gif struct {
	Type               string `json:"type"`
	Id                 string `json:"id"`
	URL                string `json:"url"`
	Tags               string `json:"tags"`
	BitlyGifURL        string `json:"bitly_gif_url"`
	BitlyFullscreenURL string `json:"bitly_fullscreen_url"`
	BitlyTiledURL      string `json:"bitly_tiled_url"`
	Images             struct {
		Original               ImageData `json:"original"`
		FixedHeight            ImageData `json:"fixed_height"`
		FixedHeightStill       ImageData `json:"fixed_height_still"`
		FixedHeightDownsampled ImageData `json:"fixed_height_downsampled"`
		FixedWidth             ImageData `json:"fixed_width"`
		FixedwidthStill        ImageData `json:"fixed_width_still"`
		FixedwidthDownsampled  ImageData `json:"fixed_width_downsampled"`
	} `json:"images"`
}

type paginatedResults struct {
	Data       []*Gif `json:"data"`
	Pagination struct {
		TotalCount int `json:"total_count"`
	} `json:"pagination"`
}

type singleResult struct {
	Data *Gif `json:"data"`
}
