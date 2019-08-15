package crates

type attributes struct {
	Repository string `json:"repository"`
	Branch     string `json:"branch"`
}

type badges struct {
	BadgeType  string     `json:"badge_type"`
	Attributes attributes `json:"attributes"`
}

type links struct {
	VersionDownloads    string `json:"version_downloads"`
	Versions            int    `json:"versions"`
	Owners              string `json:"owners"`
	OwnerTeam           string `json:"owner_team"`
	OwnerUser           string `json:"owner_user"`
	ReverseDependencies string `json:"reverse_dependencies"`
}

// Crate is basic crate metadata
type Crate struct {
	Id              string   `json:"id"`
	Name            string   `json:"name"`
	UpdatedAt       string   `json:"updated_at"`
	Versions        []int    `json:"versions"`
	Keywords        []string `json:"keywords"`
	Categories      []string `json:"categories"`
	Badges          []badges `json:"badges"`
	CreatedAt       string   `json:"created_at"`
	Downloads       int      `json:"downloads"`
	RecentDownloads int      `json:"recent_downloads"`
	MaxVersion      string   `json:"max_version"`
	Description     string   `json:"description"`
	Homepage        string   `json:"homepage"`
	Documentation   string   `json:"documentation"`
	Repository      string   `json:"repository"`
	Links           links    `json:"links"`
	ExactMatch      bool     `json:"exact_match"`
}

type features struct {
}

type dlinks struct {
	Dependencies     string `json:"dependencies"`
	VersionDownloads string `json:"version_downloads"`
	Authors          string `json:"authors"`
}

type published_by struct {
	Id     int    `json:"id"`
	Login  string `json:"login"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
	Url    string `json:"url"`
}

// Versions are specific releases
type Versions struct {
	Id          int          `json:"id"`
	Crate       string       `json:"crate"`
	Num         string       `json:"num"`
	DlPath      string       `json:"dl_path"`
	ReadmePath  string       `json:"readme_path"`
	UpdatedAt   string       `json:"updated_at"`
	CreatedAt   string       `json:"created_at"`
	Downloads   int          `json:"downloads"`
	Features    features     `json:"features"`
	Yanked      bool         `json:"yanked"`
	License     string       `json:"license"`
	Links       dlinks       `json:"links"`
	CrateSize   int          `json:"crate_size"`
	PublishedBy published_by `json:"published_by"`
}

type keywords struct {
	Id        string `json:"id"`
	Keyword   string `json:"keyword"`
	CreatedAt string `json:"created_at"`
	CratesCnt int    `json:"crates_cnt"`
}

type categories struct {
	Id          string `json:"id"`
	Category    string `json:"category"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	CratesCnt   int    `json:"crates_cnt"`
}

// CrateData is the root data structure of a crates.io/crates JSON
type CrateData struct {
	Crate      Crate        `json:"crate"`
	Versions   []Versions   `json:"versions"`
	Keywords   []keywords   `json:"keywords"`
	Categories []categories `json:"categories"`
}

type Dependencies struct {
	Id              int    `json:"id"`
	VersionId       int    `json:"version_id"`
	CrateId         string `json:"crate_id"`
	Req             string `json:"req"`
	Optional        bool   `json:"optional"`
	DefaultFeatures bool   `json:"default_features"`
	Target          string `json:"target"`
	Kind            string `json:"kind"`
	Downloads       int    `json:"downloads"`
}

type DepRoot struct {
	Dependencies []Dependencies `json:"dependencies"`
}
