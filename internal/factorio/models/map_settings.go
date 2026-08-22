package models

type DifficultySettings struct {
	TechnologyPriceMultiplier float64 `json:"technology_price_multiplier"`
	SpoilTimeModifier         float64 `json:"spoil_time_modifier"`
}

type PollutionSettings struct {
	Enabled                                   bool    `json:"enabled"`
	DiffusionRatio                            float64 `json:"diffusion_ratio"`
	MinToDiffuse                              float64 `json:"min_to_diffuse"`
	Ageing                                    float64 `json:"ageing"`
	ExpectedMaxPerChunk                       float64 `json:"expected_max_per_chunk"`
	MinToShowPerChunk                         float64 `json:"min_to_show_per_chunk"`
	MinPollutionToDamageTrees                 float64 `json:"min_pollution_to_damage_trees"`
	PollutionWithMaxForestDamage              float64 `json:"pollution_with_max_forest_damage"`
	PollutionPerTreeDamage                    float64 `json:"pollution_per_tree_damage"`
	PollutionRestoredPerTreeDamage            float64 `json:"pollution_restored_per_tree_damage"`
	MaxPollutionToRestoreTrees                float64 `json:"max_pollution_to_restore_trees"`
	EnemyAttackPollutionConsumptionModifier float64 `json:"enemy_attack_pollution_consumption_modifier"`
}

type EnemyEvolutionSettings struct {
	Enabled         bool    `json:"enabled"`
	TimeFactor      float64 `json:"time_factor"`
	DestroyFactor   float64 `json:"destroy_factor"`
	PollutionFactor float64 `json:"pollution_factor"`
}

type EnemyExpansionSettings struct {
	Enabled                          bool    `json:"enabled"`
	MaxExpansionDistance             float64 `json:"max_expansion_distance"`
	FriendlyBaseInfluenceRadius      float64 `json:"friendly_base_influence_radius"`
	EnemyBuildingInfluenceRadius     float64 `json:"enemy_building_influence_radius"`
	BuildingCoefficient              float64 `json:"building_coefficient"`
	OtherBaseCoefficient             float64 `json:"other_base_coefficient"`
	NeighbouringChunkCoefficient     float64 `json:"neighbouring_chunk_coefficient"`
	NeighbouringBaseChunkCoefficient float64 `json:"neighbouring_base_chunk_coefficient"`
	MaxCollidingTilesCoefficient     float64 `json:"max_colliding_tiles_coefficient"`
	SettlerGroupMinSize              int     `json:"settler_group_min_size"`
	SettlerGroupMaxSize              int     `json:"settler_group_max_size"`
	MinExpansionCooldown             int     `json:"min_expansion_cooldown"`
	MaxExpansionCooldown             int     `json:"max_expansion_cooldown"`
}

type UnitGroupSettings struct {
	MinGroupGatheringTime          int     `json:"min_group_gathering_time"`
	MaxGroupGatheringTime          int     `json:"max_group_gathering_time"`
	MaxWaitTimeForLateMembers       int     `json:"max_wait_time_for_late_members"`
	MaxGroupRadius                 float64 `json:"max_group_radius"`
	MinGroupRadius                 float64 `json:"min_group_radius"`
	MaxMemberSpeedupWhenBehind     float64 `json:"max_member_speedup_when_behind"`
	MaxMemberSlowdownWhenAhead     float64 `json:"max_member_slowdown_when_ahead"`
	MaxGroupSlowdownFactor         float64 `json:"max_group_slowdown_factor"`
	MaxGroupMemberFallbackFactor   float64 `json:"max_group_member_fallback_factor"`
	MemberDisownDistance           float64 `json:"member_disown_distance"`
	TickToleranceWhenMemberArrives int     `json:"tick_tolerance_when_member_arrives"`
	MaxGatheringUnitGroups         int     `json:"max_gathering_unit_groups"`
	MaxUnitGroupSize               int     `json:"max_unit_group_size"`
}

type SteeringSetting struct {
	Radius                      float64 `json:"radius"`
	SeparationForce             float64 `json:"separation_force"`
	SeparationFactor            float64 `json:"separation_factor"`
	ForceUnitFuzzyGotoBehavior bool    `json:"force_unit_fuzzy_goto_behavior"`
}

type SteeringSettings struct {
	Default SteeringSetting `json:"default"`
	Moving  SteeringSetting `json:"moving"`
}

type PathFinderSettings struct {
	Fwd2BwdRatio                                    float64   `json:"fwd2bwd_ratio"`
	GoalPressureRatio                               float64   `json:"goal_pressure_ratio"`
	MaxStepsWorkedPerTick                           int       `json:"max_steps_worked_per_tick"`
	MaxWorkDonePerTick                              int       `json:"max_work_done_per_tick"`
	UsePathCache                                    bool      `json:"use_path_cache"`
	ShortCacheSize                                  int       `json:"short_cache_size"`
	LongCacheSize                                   int       `json:"long_cache_size"`
	ShortCacheMinCacheableDistance                  float64   `json:"short_cache_min_cacheable_distance"`
	ShortCacheMinAlgoStepsToCache                   int       `json:"short_cache_min_algo_steps_to_cache"`
	LongCacheMinCacheableDistance                   float64   `json:"long_cache_min_cacheable_distance"`
	CacheMaxConnectToCacheStepsMultiplier           float64   `json:"cache_max_connect_to_cache_steps_multiplier"`
	CacheAcceptPathStartDistanceRatio               float64   `json:"cache_accept_path_start_distance_ratio"`
	CacheAcceptPathEndDistanceRatio                 float64   `json:"cache_accept_path_end_distance_ratio"`
	NegativeCacheAcceptPathStartDistanceRatio       float64   `json:"negative_cache_accept_path_start_distance_ratio"`
	NegativeCacheAcceptPathEndDistanceRatio         float64   `json:"negative_cache_accept_path_end_distance_ratio"`
	CachePathStartDistanceRatingMultiplier         float64   `json:"cache_path_start_distance_rating_multiplier"`
	CachePathEndDistanceRatingMultiplier           float64   `json:"cache_path_end_distance_rating_multiplier"`
	StaleEnemyWithSameDestinationCollisionPenalty   float64   `json:"stale_enemy_with_same_destination_collision_penalty"`
	IgnoreMovingEnemyCollisionDistance              float64   `json:"ignore_moving_enemy_collision_distance"`
	EnemyWithDifferentDestinationCollisionPenalty   float64   `json:"enemy_with_different_destination_collision_penalty"`
	GeneralEntityCollisionPenalty                   float64   `json:"general_entity_collision_penalty"`
	GeneralEntitySubsequentCollisionPenalty          float64   `json:"general_entity_subsequent_collision_penalty"`
	ExtendedCollisionPenalty                        float64   `json:"extended_collision_penalty"`
	MaxClientsToAcceptAnyNewRequest                 int       `json:"max_clients_to_accept_any_new_request"`
	MaxClientsToAcceptShortNewRequest                int       `json:"max_clients_to_accept_short_new_request"`
	DirectDistanceToConsiderShortRequest            float64   `json:"direct_distance_to_consider_short_request"`
	ShortRequestMaxSteps                            int       `json:"short_request_max_steps"`
	ShortRequestRatio                               float64   `json:"short_request_ratio"`
	MinStepsToCheckPathFindTermination             int       `json:"min_steps_to_check_path_find_termination"`
	StartToGoalCostMultiplierToTerminatePathFind    float64   `json:"start_to_goal_cost_multiplier_to_terminate_path_find"`
	OverloadLevels                                  []int     `json:"overload_levels"`
	OverloadMultipliers                             []float64 `json:"overload_multipliers"`
	NegativePathCacheDelayInterval                  int       `json:"negative_path_cache_delay_interval"`
}

type AsteroidSettings struct {
	SpawningRate                 float64 `json:"spawning_rate"`
	MaxRayPortalsExpandedPerTick int     `json:"max_ray_portals_expanded_per_tick"`
}

type MapSettings struct {
	DifficultySettings     DifficultySettings     `json:"difficulty_settings"`
	Pollution              PollutionSettings      `json:"pollution"`
	EnemyEvolution         EnemyEvolutionSettings `json:"enemy_evolution"`
	EnemyExpansion         EnemyExpansionSettings `json:"enemy_expansion"`
	UnitGroup              UnitGroupSettings      `json:"unit_group"`
	Steering               SteeringSettings       `json:"steering"`
	PathFinder             PathFinderSettings     `json:"path_finder"`
	Asteroids              AsteroidSettings       `json:"asteroids"`
	MaxFailedBehaviorCount int                    `json:"max_failed_behavior_count"`
}
