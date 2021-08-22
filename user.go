package mediawiki

import "time"

const userURL = "/w/api.php"

// User mediawiki user representation.
type User struct {
	UserID           int           `json:"userid,omitempty"`
	Name             string        `json:"name"`
	EditCount        int           `json:"editcount,omitempty"`
	Registration     time.Time     `json:"registration,omitempty"`
	Groups           []string      `json:"groups,omitempty"`
	GroupMemberships []interface{} `json:"groupmemberships,omitempty"`
	Emailable        bool          `json:"emailable,omitempty"`
	Missing          bool          `json:"missing,omitempty"`
}

type userResponse struct {
	Batchcomplete bool `json:"batchcomplete"`
	Query         struct {
		Users []User `json:"users"`
	} `json:"query"`
}
