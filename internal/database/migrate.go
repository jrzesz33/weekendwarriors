package database

import (
	"database/sql"
	"fmt"

	"github.com/rs/zerolog/log"
)

// Migrate runs database migrations
func Migrate(db *sql.DB) error {
	// Create migrations table if it doesn't exist
	if err := createMigrationsTable(db); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	migrations := getMigrations()

	for _, migration := range migrations {
		applied, err := isMigrationApplied(db, migration.Version)
		if err != nil {
			return fmt.Errorf("failed to check migration %s: %w", migration.Version, err)
		}

		if applied {
			log.Debug().Str("version", migration.Version).Msg("Migration already applied")
			continue
		}

		log.Info().Str("version", migration.Version).Msg("Applying migration")
		if err := applyMigration(db, migration); err != nil {
			return fmt.Errorf("failed to apply migration %s: %w", migration.Version, err)
		}
	}

	return nil
}

type Migration struct {
	Version string
	Name    string
	SQL     string
}

func createMigrationsTable(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS migrations (
			version TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`
	_, err := db.Exec(query)
	return err
}

func isMigrationApplied(db *sql.DB, version string) (bool, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM migrations WHERE version = ?", version).Scan(&count)
	return count > 0, err
}

func applyMigration(db *sql.DB, migration Migration) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Execute migration SQL
	if _, err := tx.Exec(migration.SQL); err != nil {
		return err
	}

	// Record migration
	if _, err := tx.Exec("INSERT INTO migrations (version, name) VALUES (?, ?)", migration.Version, migration.Name); err != nil {
		return err
	}

	return tx.Commit()
}

func getMigrations() []Migration {
	return []Migration{
		{
			Version: "001",
			Name:    "Create initial schema",
			SQL: `
				-- Course data table
				CREATE TABLE course_data (
					id TEXT PRIMARY KEY,
					course_name TEXT NOT NULL,
					hole INTEGER NOT NULL,
					par INTEGER NOT NULL,
					handicap_ranking INTEGER NOT NULL,
					yardage INTEGER,
					description TEXT,
					UNIQUE(course_name, hole)
				);

				-- Games table
				CREATE TABLE games (
					id TEXT PRIMARY KEY,
					course TEXT NOT NULL,
					status TEXT NOT NULL CHECK (status IN ('setup', 'in_progress', 'completed', 'abandoned')),
					handicap_enabled BOOLEAN NOT NULL DEFAULT 1,
					side_bets TEXT, -- JSON array
					share_token TEXT UNIQUE NOT NULL,
					spectator_token TEXT UNIQUE NOT NULL,
					current_hole INTEGER,
					created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					started_at TIMESTAMP,
					completed_at TIMESTAMP,
					final_results TEXT -- JSON
				);

				-- Players table
				CREATE TABLE players (
					id TEXT PRIMARY KEY,
					game_id TEXT NOT NULL,
					name TEXT NOT NULL,
					handicap REAL NOT NULL,
					gender TEXT CHECK (gender IN ('male', 'female', 'other')),
					position INTEGER NOT NULL,
					created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					FOREIGN KEY (game_id) REFERENCES games(id) ON DELETE CASCADE,
					UNIQUE(game_id, name),
					UNIQUE(game_id, position)
				);

				-- Scores table
				CREATE TABLE scores (
					id TEXT PRIMARY KEY,
					player_id TEXT NOT NULL,
					game_id TEXT NOT NULL,
					hole INTEGER NOT NULL,
					strokes INTEGER NOT NULL,
					putts INTEGER NOT NULL,
					par INTEGER NOT NULL,
					handicap_stroke BOOLEAN NOT NULL DEFAULT 0,
					score_to_par INTEGER NOT NULL,
					effective_score INTEGER NOT NULL,
					created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					FOREIGN KEY (player_id) REFERENCES players(id) ON DELETE CASCADE,
					FOREIGN KEY (game_id) REFERENCES games(id) ON DELETE CASCADE,
					UNIQUE(player_id, hole)
				);

				-- Side bet calculations table
				CREATE TABLE side_bet_calculations (
					id TEXT PRIMARY KEY,
					game_id TEXT NOT NULL,
					player_id TEXT NOT NULL,
					bet_type TEXT NOT NULL CHECK (bet_type IN ('best_nine', 'putt_putt_poker')),
					calculation_data TEXT NOT NULL, -- JSON
					current_position INTEGER,
					final_position INTEGER,
					is_winner BOOLEAN DEFAULT 0,
					calculated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					FOREIGN KEY (game_id) REFERENCES games(id) ON DELETE CASCADE,
					FOREIGN KEY (player_id) REFERENCES players(id) ON DELETE CASCADE,
					UNIQUE(player_id, bet_type)
				);

				-- Putt putt poker cards table
				CREATE TABLE putt_putt_poker_cards (
					id TEXT PRIMARY KEY,
					player_id TEXT NOT NULL,
					game_id TEXT NOT NULL,
					hole INTEGER,
					action TEXT NOT NULL CHECK (action IN ('starting', 'one_putt', 'hole_in_one', 'penalty')),
					cards_change INTEGER NOT NULL,
					penalty_amount REAL,
					total_cards INTEGER NOT NULL,
					created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					FOREIGN KEY (player_id) REFERENCES players(id) ON DELETE CASCADE,
					FOREIGN KEY (game_id) REFERENCES games(id) ON DELETE CASCADE
				);

				-- Poker hands table
				CREATE TABLE poker_hands (
					id TEXT PRIMARY KEY,
					game_id TEXT NOT NULL,
					player_id TEXT NOT NULL,
					total_cards_earned INTEGER NOT NULL,
					dealt_cards TEXT NOT NULL, -- JSON array
					best_hand_cards TEXT NOT NULL, -- JSON array
					hand_type TEXT NOT NULL,
					hand_rank INTEGER NOT NULL,
					hand_description TEXT,
					position INTEGER NOT NULL,
					deal_timestamp TIMESTAMP NOT NULL,
					random_seed TEXT NOT NULL,
					FOREIGN KEY (game_id) REFERENCES games(id) ON DELETE CASCADE,
					FOREIGN KEY (player_id) REFERENCES players(id) ON DELETE CASCADE,
					UNIQUE(game_id, player_id)
				);

				-- Create indexes
				CREATE INDEX idx_games_status ON games(status);
				CREATE INDEX idx_games_share_token ON games(share_token);
				CREATE INDEX idx_games_spectator_token ON games(spectator_token);
				CREATE INDEX idx_games_created_at ON games(created_at);
				CREATE INDEX idx_players_game_id ON players(game_id);
				CREATE INDEX idx_scores_player_id ON scores(player_id);
				CREATE INDEX idx_scores_game_id ON scores(game_id);
				CREATE INDEX idx_scores_hole ON scores(hole);
				CREATE INDEX idx_scores_game_hole ON scores(game_id, hole);
				CREATE INDEX idx_sidebet_game_bet ON side_bet_calculations(game_id, bet_type);
				CREATE INDEX idx_sidebet_player_bet ON side_bet_calculations(player_id, bet_type);
				CREATE INDEX idx_poker_cards_player ON putt_putt_poker_cards(player_id);
				CREATE INDEX idx_poker_cards_game ON putt_putt_poker_cards(game_id);
				CREATE INDEX idx_poker_cards_hole ON putt_putt_poker_cards(hole);
				CREATE INDEX idx_poker_hands_game ON poker_hands(game_id);
				CREATE INDEX idx_poker_hands_rank ON poker_hands(hand_rank);
				CREATE INDEX idx_course_data_course ON course_data(course_name);
			`,
		},
		{
			Version: "002",
			Name:    "Insert Diamond Run course data",
			SQL: `
				INSERT INTO course_data (id, course_name, hole, par, handicap_ranking) VALUES
				('hole_dr_01', 'diamond-run', 1, 4, 10),
				('hole_dr_02', 'diamond-run', 2, 3, 18),
				('hole_dr_03', 'diamond-run', 3, 5, 2),
				('hole_dr_04', 'diamond-run', 4, 4, 8),
				('hole_dr_05', 'diamond-run', 5, 3, 16),
				('hole_dr_06', 'diamond-run', 6, 4, 12),
				('hole_dr_07', 'diamond-run', 7, 4, 6),
				('hole_dr_08', 'diamond-run', 8, 3, 14),
				('hole_dr_09', 'diamond-run', 9, 5, 4),
				('hole_dr_10', 'diamond-run', 10, 4, 9),
				('hole_dr_11', 'diamond-run', 11, 3, 17),
				('hole_dr_12', 'diamond-run', 12, 5, 1),
				('hole_dr_13', 'diamond-run', 13, 3, 15),
				('hole_dr_14', 'diamond-run', 14, 4, 11),
				('hole_dr_15', 'diamond-run', 15, 4, 7),
				('hole_dr_16', 'diamond-run', 16, 4, 13),
				('hole_dr_17', 'diamond-run', 17, 5, 3),
				('hole_dr_18', 'diamond-run', 18, 4, 5);
			`,
		},
	}
}