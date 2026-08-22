package models

type ServerVisibilitySettings struct {
	Public bool `json:"public"`
	Lan    bool `json:"lan"`
}

type ServerSettings struct {
	Name                                 string                   `json:"name"`
	Description                          string                   `json:"description"`
	Tags                                 []string                 `json:"tags"`
	MaxPlayers                           int                      `json:"max_players"`
	Visibility                           ServerVisibilitySettings `json:"visibility"`
	Username                             string                   `json:"username"`
	Password                             string                   `json:"password"`
	Token                                string                   `json:"token"`
	GamePassword                         string                   `json:"game_password"`
	RequireUserVerification              bool                     `json:"require_user_verification"`
	MaxUploadInKilobytesPerSecond        int                      `json:"max_upload_in_kilobytes_per_second"`
	MaxUploadSlots                       int                      `json:"max_upload_slots"`
	MinimumLatencyInTicks                int                      `json:"minimum_latency_in_ticks"`
	MaxHeartbeatsPerSecond               int                      `json:"max_heartbeats_per_second"`
	IgnorePlayerLimitForReturningPlayers bool                     `json:"ignore_player_limit_for_returning_players"`
	AllowCommands                        string                   `json:"allow_commands"`
	AutosaveInterval                     int                      `json:"autosave_interval"`
	AutosaveSlots                        int                      `json:"autosave_slots"`
	AFKAutokickInterval                  int                      `json:"afk_autokick_interval"`
	AutoPause                            bool                     `json:"auto_pause"`
	AutoPauseWhenPlayersConnect          bool                     `json:"auto_pause_when_players_connect"`
	OnlyAdminsCanPauseTheGame            bool                     `json:"only_admins_can_pause_the_game"`
	AutosaveOnlyOnServer                 bool                     `json:"autosave_only_on_server"`
	NonBlockingSaving                    bool                     `json:"non_blocking_saving"`
	MinimumSegmentSize                   int                      `json:"minimum_segment_size"`
	MinimumSegmentSizePeerCount          int                      `json:"minimum_segment_size_peer_count"`
	MaximumSegmentSize                   int                      `json:"maximum_segment_size"`
	MaximumSegmentSizePeerCount          int                      `json:"maximum_segment_size_peer_count"`
}
