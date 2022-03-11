package modules

import "time"

type FakeID struct {
	Results []struct {
		Gender string `json:"gender"`
		Name   struct {
			Title string `json:"title"`
			First string `json:"first"`
			Last  string `json:"last"`
		} `json:"name"`
		Location struct {
			Street struct {
				Number int    `json:"number"`
				Name   string `json:"name"`
			} `json:"street"`
			City     string `json:"city"`
			State    string `json:"state"`
			Country  string `json:"country"`
			Postcode int    `json:"postcode"`
		} `json:"location"`
		Email string `json:"email"`
		Dob   struct {
			Date time.Time `json:"date"`
			Age  int       `json:"age"`
		} `json:"dob"`
		Phone string `json:"phone"`
		Cell  string `json:"cell"`
		Nat   string `json:"nat"`
	} `json:"results"`
}

type InstSearch struct {
	Users []struct {
		Position int `json:"position"`
		User     struct {
			Pk                         string      `json:"pk"`
			Username                   string      `json:"username"`
			FullName                   string      `json:"full_name"`
			IsPrivate                  bool        `json:"is_private"`
			ProfilePicURL              string      `json:"profile_pic_url"`
			ProfilePicID               string      `json:"profile_pic_id"`
			IsVerified                 bool        `json:"is_verified"`
			FollowFrictionType         int         `json:"follow_friction_type"`
			HasAnonymousProfilePicture bool        `json:"has_anonymous_profile_picture"`
			HasHighlightReels          bool        `json:"has_highlight_reels"`
			LatestReelMedia            int         `json:"latest_reel_media"`
			LiveBroadcastID            interface{} `json:"live_broadcast_id"`
			ShouldShowCategory         bool        `json:"should_show_category"`
			Seen                       int         `json:"seen"`
		} `json:"user,omitempty"`
	} `json:"users"`
}

type PintrestResp struct {
	ResourceResponse struct {
		Status string `json:"status"`
		Data   struct {
			Results []struct {
				Objects []struct {
					RecentPinImages struct {
						One92X []struct {
							URL           string `json:"url"`
							Width         int    `json:"width"`
							Height        int    `json:"height"`
							DominantColor string `json:"dominant_color"`
						} `json:"192x"`
					} `json:"recent_pin_images"`
				} `json:"objects"`
			} `json:"results"`
		} `json:"data"`
	} `json:"resource_response"`
}

type UrbanDict struct {
	List []struct {
		Definition string `json:"definition"`
		ThumbsUp   int    `json:"thumbs_up"`
		Author     string `json:"author"`
		Word       string `json:"word"`
		Example    string `json:"example"`
		ThumbsDown int    `json:"thumbs_down"`
	} `json:"list"`
}

type Bin struct {
	Number struct {
		Length int `json:"length"`
	} `json:"number"`
	Scheme  string `json:"scheme"`
	Type    string `json:"type"`
	Brand   string `json:"brand"`
	Prepaid bool   `json:"prepaid"`
	Country struct {
		Numeric  string `json:"numeric"`
		Alpha2   string `json:"alpha2"`
		Name     string `json:"name"`
		Emoji    string `json:"emoji"`
		Currency string `json:"currency"`
	} `json:"country"`
	Bank struct {
		Name  string `json:"name"`
		Phone string `json:"phone"`
	} `json:"bank"`
}

type TGraph struct {
	Ok     bool `json:"ok"`
	Result struct {
		URL string `json:"url"`
	} `json:"result"`
}

type AuddApi struct {
	Status string `json:"status"`
	Error  struct {
		ErrorCode    int    `json:"error_code"`
		ErrorMessage string `json:"error_message"`
	} `json:"error"`
	Result struct {
		Artist      string `json:"artist"`
		Title       string `json:"title"`
		Album       string `json:"album"`
		ReleaseDate string `json:"release_date"`
		Label       string `json:"label"`
		SongLink    string `json:"song_link"`
		AppleMusic  struct {
			Previews []struct {
				URL string `json:"url"`
			} `json:"previews"`
			Artwork struct {
				Width  int    `json:"width"`
				Height int    `json:"height"`
				URL    string `json:"url"`
			} `json:"artwork"`
		} `json:"apple_music"`
		Spotify struct {
			Album struct {
				Name   string `json:"name"`
				Images []struct {
					Width  int    `json:"width"`
					Height int    `json:"height"`
					URL    string `json:"url"`
				} `json:"images"`
			} `json:"album"`
			ExternalUrls struct {
				Spotify string `json:"spotify"`
			} `json:"external_urls"`
		} `json:"spotify"`
	} `json:"result"`
}

type YoutubeVideo struct {
	Title     string `json:"title,omitempty"`
	ID        string `json:"id,omitempty"`
	Thumbnail string `json:"thumbnail,omitempty"`
	Author    string `json:"author,omitempty"`
	Published string `json:"published,omitempty"`
	Duration  string `json:"duration,omitempty"`
	Views     string `json:"views,omitempty"`
}
type IPData struct {
	IP       string `json:"ip"`
	Hostname string `json:"hostname"`
	City     string `json:"city"`
	Region   string `json:"region"`
	Country  string `json:"country"`
	Loc      string `json:"loc"`
	Org      string `json:"org"`
	Postal   string `json:"postal"`
	Timezone string `json:"timezone"`
	Error    string `json:"error,omitempty"`
}

type IPData struct {
	IP       string `json:"ip"`
	Hostname string `json:"hostname"`
	City     string `json:"city"`
	Region   string `json:"region"`
	Country  string `json:"country"`
	Loc      string `json:"loc"`
	Org      string `json:"org"`
	Postal   string `json:"postal"`
	Timezone string `json:"timezone"`
}
