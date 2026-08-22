package models

import (
	"encoding/json"
	"testing"
)

func TestMapGenSettingsUnmarshal(t *testing.T) {
	data := []byte(`{
		"width": 0,
		"height": 0,
		"starting_area": 1,
		"peaceful_mode": false,
		"autoplace_controls": {
			"coal": {"frequency": 1, "size": 1, "richness": 1}
		},
		"cliff_settings": {
			"name": "cliff",
			"cliff_elevation_0": 10,
			"cliff_elevation_interval": 40,
			"richness": 1
		},
		"property_expression_names": {
			"control:moisture:frequency": "1"
		},
		"starting_points": [
			{"x": 0, "y": 0}
		],
		"seed": null
	}`)

	var cfg MapGenSettings
	if err := json.Unmarshal(data, &cfg); err != nil {
		t.Fatalf("failed to unmarshal MapGenSettings: %v", err)
	}

	if cfg.CliffSettings.Name != "cliff" {
		t.Errorf("expected cliff settings name 'cliff', got: %s", cfg.CliffSettings.Name)
	}
}

func TestMapSettingsUnmarshal(t *testing.T) {
	data := []byte(`{
		"difficulty_settings": {
			"technology_price_multiplier": 1,
			"spoil_time_modifier": 1
		},
		"pollution": {
			"enabled": true,
			"diffusion_ratio": 0.02,
			"min_to_diffuse": 15,
			"ageing": 1,
			"expected_max_per_chunk": 150,
			"min_to_show_per_chunk": 50,
			"min_pollution_to_damage_trees": 60,
			"pollution_with_max_forest_damage": 150,
			"pollution_per_tree_damage": 50,
			"pollution_restored_per_tree_damage": 10,
			"max_pollution_to_restore_trees": 20,
			"enemy_attack_pollution_consumption_modifier": 1
		},
		"enemy_evolution": {
			"enabled": true,
			"time_factor": 0.000004,
			"destroy_factor": 0.002,
			"pollution_factor": 0.0000009
		},
		"enemy_expansion": {
			"enabled": true,
			"max_expansion_distance": 7,
			"friendly_base_influence_radius": 2,
			"enemy_building_influence_radius": 2,
			"building_coefficient": 0.1,
			"other_base_coefficient": 2.0,
			"neighbouring_chunk_coefficient": 0.5,
			"neighbouring_base_chunk_coefficient": 0.4,
			"max_colliding_tiles_coefficient": 0.9,
			"settler_group_min_size": 5,
			"settler_group_max_size": 20,
			"min_expansion_cooldown": 14400,
			"max_expansion_cooldown": 216000
		},
		"unit_group": {
			"min_group_gathering_time": 3600,
			"max_group_gathering_time": 36000,
			"max_wait_time_for_late_members": 7200,
			"max_group_radius": 30.0,
			"min_group_radius": 5.0,
			"max_member_speedup_when_behind": 1.4,
			"max_member_slowdown_when_ahead": 0.6,
			"max_group_slowdown_factor": 0.3,
			"max_group_member_fallback_factor": 3,
			"member_disown_distance": 10,
			"tick_tolerance_when_member_arrives": 60,
			"max_gathering_unit_groups": 30,
			"max_unit_group_size": 200
		},
		"steering": {
			"default": {
				"radius": 1.2,
				"separation_force": 0.005,
				"separation_factor": 1.2,
				"force_unit_fuzzy_goto_behavior": false
			},
			"moving": {
				"radius": 3,
				"separation_force": 0.01,
				"separation_factor": 3,
				"force_unit_fuzzy_goto_behavior": false
			}
		},
		"path_finder": {
			"fwd2bwd_ratio": 5,
			"goal_pressure_ratio": 2,
			"max_steps_worked_per_tick": 1000,
			"max_work_done_per_tick": 8000,
			"use_path_cache": true,
			"short_cache_size": 5,
			"long_cache_size": 25,
			"short_cache_min_cacheable_distance": 10,
			"short_cache_min_algo_steps_to_cache": 50,
			"long_cache_min_cacheable_distance": 30,
			"cache_max_connect_to_cache_steps_multiplier": 100,
			"cache_accept_path_start_distance_ratio": 0.2,
			"cache_accept_path_end_distance_ratio": 0.15,
			"negative_cache_accept_path_start_distance_ratio": 0.3,
			"negative_cache_accept_path_end_distance_ratio": 0.3,
			"cache_path_start_distance_rating_multiplier": 10,
			"cache_path_end_distance_rating_multiplier": 20,
			"stale_enemy_with_same_destination_collision_penalty": 30,
			"ignore_moving_enemy_collision_distance": 5,
			"enemy_with_different_destination_collision_penalty": 30,
			"general_entity_collision_penalty": 10,
			"general_entity_subsequent_collision_penalty": 3,
			"extended_collision_penalty": 3,
			"max_clients_to_accept_any_new_request": 10,
			"max_clients_to_accept_short_new_request": 100,
			"direct_distance_to_consider_short_request": 100,
			"short_request_max_steps": 1000,
			"short_request_ratio": 0.5,
			"min_steps_to_check_path_find_termination": 2000,
			"start_to_goal_cost_multiplier_to_terminate_path_find": 2000.0,
			"overload_levels": [0, 100, 500],
			"overload_multipliers": [2, 3, 4],
			"negative_path_cache_delay_interval": 20
		},
		"asteroids": {
			"spawning_rate": 1,
			"max_ray_portals_expanded_per_tick": 100
		},
		"max_failed_behavior_count": 3
	}`)

	var cfg MapSettings
	if err := json.Unmarshal(data, &cfg); err != nil {
		t.Fatalf("failed to unmarshal MapSettings: %v", err)
	}

	if !cfg.Pollution.Enabled {
		t.Errorf("expected pollution to be enabled")
	}
	if cfg.MaxFailedBehaviorCount != 3 {
		t.Errorf("expected MaxFailedBehaviorCount = 3, got: %d", cfg.MaxFailedBehaviorCount)
	}
}

func TestServerSettingsUnmarshal(t *testing.T) {
	data := []byte(`{
		"name": "Factorio Server",
		"description": "Factorio Docker Server",
		"tags": ["game", "tags"],
		"max_players": 0,
		"visibility": {
			"public": true,
			"lan": true
		},
		"username": "",
		"password": "",
		"token": "",
		"game_password": "",
		"require_user_verification": true,
		"max_upload_in_kilobytes_per_second": 0,
		"max_upload_slots": 5,
		"minimum_latency_in_ticks": 0,
		"max_heartbeats_per_second": 60,
		"ignore_player_limit_for_returning_players": false,
		"allow_commands": "admins-only",
		"autosave_interval": 10,
		"autosave_slots": 5,
		"afk_autokick_interval": 0,
		"auto_pause": true,
		"auto_pause_when_players_connect": false,
		"only_admins_can_pause_the_game": true,
		"autosave_only_on_server": true,
		"non_blocking_saving": false,
		"minimum_segment_size": 25,
		"minimum_segment_size_peer_count": 20,
		"maximum_segment_size": 100,
		"maximum_segment_size_peer_count": 10
	}`)

	var cfg ServerSettings
	if err := json.Unmarshal(data, &cfg); err != nil {
		t.Fatalf("failed to unmarshal ServerSettings: %v", err)
	}

	if cfg.Name != "Factorio Server" {
		t.Errorf("expected Name = 'Factorio Server', got: %s", cfg.Name)
	}
	if !cfg.Visibility.Public {
		t.Errorf("expected Visibility.Public = true")
	}
}

func TestPlayerListsUnmarshal(t *testing.T) {
	data := []byte(`["player1", "player2"]`)

	var admins AdminList
	if err := json.Unmarshal(data, &admins); err != nil {
		t.Fatalf("failed to unmarshal AdminList: %v", err)
	}

	if len(admins) != 2 || admins[0] != "player1" {
		t.Errorf("unexpected AdminList content: %v", admins)
	}
}
