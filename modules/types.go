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

type AuddEntrp struct {
	ExecutionTime string `json:"execution_time"`
	Result        []struct {
		Offset string `json:"offset"`
		Songs  []struct {
			Album       string `json:"album"`
			Artist      string `json:"artist"`
			Label       string `json:"label"`
			ReleaseDate string `json:"release_date"`
			SongLink    string `json:"song_link"`
			Timecode    string `json:"timecode"`
			Title       string `json:"title"`
		} `json:"songs"`
	} `json:"result"`
	Status string `json:"status"`
}
type YoutubeVideo struct {
	ID        string `json:"id"`
	URL       string `json:"url"`
	Title     string `json:"title"`
	Thumbnail string `json:"thumbnail"`
	Channel   string `json:"channel"`
}

type IPData struct {
	IP, Hostname, City, Region, Country, Loc, Org, Postal, Timezone string
	Error                                                           struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

type Pins struct {
	ResourceResponse struct {
		Data struct {
			Results []struct {
				Images struct {
					Orig struct {
						URL string `json:"url,omitempty"`
					} `json:"orig,omitempty"`
				} `json:"images,omitempty"`
			} `json:"results,omitempty"`
		} `json:"data,omitempty"`
	} `json:"resource_response,omitempty"`
}

type MusicSearch struct {
	Tracks struct {
		Hits []struct {
			Track struct {
				Key      string `json:"key"`
				Title    string `json:"title"`
				Subtitle string `json:"subtitle"`
				Images   struct {
					Background string `json:"background"`
					Coverart   string `json:"coverart"`
					Coverarthq string `json:"coverarthq"`
					Joecolor   string `json:"joecolor"`
				} `json:"images"`
			} `json:"track"`
		} `json:"hits"`
	} `json:"tracks"`
}

type Lyrics struct {
	Sections []struct {
		Type string   `json:"type"`
		Text []string `json:"text,omitempty"`
	}
}

type OxfordDict struct {
	Results []struct {
		LexicalEntries []struct {
			Entries []struct {
				Pronunciations []struct {
					AudioFile string `json:"audioFile"`
				} `json:"pronunciations"`
				Senses []struct {
					Definitions []string `json:"definitions"`
					Examples    []struct {
						Text string `json:"text"`
					} `json:"examples"`
					ID              string `json:"id"`
					SemanticClasses []struct {
						ID   string `json:"id"`
						Text string `json:"text"`
					} `json:"semanticClasses"`
					ShortDefinitions []string `json:"shortDefinitions"`
					Synonyms         []struct {
						Language string `json:"language"`
						Text     string `json:"text"`
					} `json:"synonyms"`
				} `json:"senses"`
			} `json:"entries"`
		} `json:"lexicalEntries"`
	} `json:"results"`
}

type YT []struct {
	ID            string `json:"id"`
	Title         string `json:"title"`
	PublishedTime string `json:"publishedTime"`
	Duration      string `json:"duration"`
	ViewCount     struct {
		Text  string `json:"text"`
		Short string `json:"short"`
	} `json:"viewCount"`
	Thumbnails []struct {
		URL    string `json:"url"`
		Width  int    `json:"width"`
		Height int    `json:"height"`
	} `json:"thumbnails"`
	DescriptionSnippet []struct {
		Text string `json:"text"`
		Bold bool   `json:"bold,omitempty"`
	} `json:"descriptionSnippet"`
	Channel struct {
		Name string `json:"name"`
	} `json:"channel"`
	Link string `json:"link"`
}

type YTVideo struct {
	ID            string
	Title         string
	PublishedTime string
	Duration      string
	ViewCount     string
	Thumbnail     string
	Description   string
	Channel       string
	Link          string
}
