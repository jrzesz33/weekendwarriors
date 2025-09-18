package ui

// HTML templates for different views

const createGameHTML = `
<div class="container">
	<div class="nav">
		<div class="nav-content">
			<div class="nav-brand">Golf Gamez</div>
			<div id="connection-status" class="connection-status disconnected">
				<span class="status-dot"></span>
				<span>Offline</span>
			</div>
		</div>
	</div>

	<div class="py-6">
		<div class="text-center mb-8">
			<h1 class="text-4xl font-bold text-primary mb-4">Create New Game</h1>
			<p class="text-lg text-gray-600">Start tracking your golf round with side bets</p>
		</div>

		<div class="card max-w-lg mx-auto">
			<form id="create-game-form">
				<div class="form-group">
					<label class="label" for="course-select">Golf Course</label>
					<select id="course-select" class="select" required>
						<option value="diamond-run">Diamond Run Golf Course</option>
					</select>
				</div>

				<div class="course-info mb-6">
					<div class="course-name">Diamond Run Golf Course</div>
					<div class="course-stats">
						<div class="course-stat">
							<span class="stat-value">18</span>
							<span class="stat-label">Holes</span>
						</div>
						<div class="course-stat">
							<span class="stat-value">72</span>
							<span class="stat-label">Par</span>
						</div>
						<div class="course-stat">
							<span class="stat-value">6,200</span>
							<span class="stat-label">Yards</span>
						</div>
					</div>
				</div>

				<div class="form-group">
					<label class="label">Game Options</label>
					<div class="space-y-3">
						<div class="flex items-center">
							<input type="checkbox" id="handicap-enabled" class="mr-3">
							<label for="handicap-enabled" class="text-sm">Enable Handicap Calculations</label>
						</div>
					</div>
				</div>

				<div class="form-group">
					<label class="label">Side Bets</label>
					<div class="space-y-3">
						<div class="flex items-center">
							<input type="checkbox" id="best-nine-enabled" class="mr-3">
							<label for="best-nine-enabled" class="text-sm">Best Nine</label>
						</div>
						<div class="text-xs text-gray-500 ml-6 mb-3">
							Calculate best 9 holes vs par with handicap
						</div>

						<div class="flex items-center">
							<input type="checkbox" id="putt-poker-enabled" class="mr-3">
							<label for="putt-poker-enabled" class="text-sm">Putt Putt Poker</label>
						</div>
						<div class="text-xs text-gray-500 ml-6">
							Earn cards based on putting performance
						</div>
					</div>
				</div>

				<button id="createGmBtn" type="button" onclick="window.golfGamez.createGame()" class="btn btn-primary btn-lg btn-full">
					Create Game
				</button>
			</form>
		</div>

		<div class="text-center mt-8">
			<a href="#" class="text-primary text-sm">Join existing game</a>
		</div>
	</div>
</div>
`

const gameSetupHTML = `
<div class="container">
	<div class="nav">
		<div class="nav-content">
			<div class="nav-brand">Golf Gamez</div>
			<div id="connection-status" class="connection-status connecting">
				<span class="status-dot"></span>
				<span>Connecting...</span>
			</div>
		</div>
	</div>

	<div class="py-6">
		<div class="text-center mb-6">
			<h1 class="text-3xl font-bold text-primary mb-2">Game Setup</h1>
			<div class="game-status setup">
				<span class="status-dot"></span>
				<span>Setup</span>
			</div>
		</div>

		<!-- Share Links Section -->
		<div class="share-link-container mb-6">
			<h3 class="share-link-title">Share Your Game</h3>

			<div class="mb-4">
				<label class="label text-sm">Player Share Link (Edit Access)</label>
				<div class="share-link-input">
					<input type="text" id="share-link-input" class="input" readonly>
					<button id="copy-share-btn" class="btn btn-secondary copy-button">
						<span>Copy</span>
					</button>
				</div>
			</div>

			<div>
				<label class="label text-sm">Spectator Link (View Only)</label>
				<div class="share-link-input">
					<input type="text" id="spectator-link-input" class="input" readonly>
					<button id="copy-spectator-btn" class="btn btn-secondary copy-button">
						<span>Copy</span>
					</button>
				</div>
			</div>
		</div>

		<!-- Add Player Section -->
		<div class="card mb-6">
			<div class="card-header">
				<h3 class="card-title">Add Players</h3>
				<p class="card-subtitle">Add up to 4 players to your game</p>
			</div>

			<form id="add-player-form">
				<div class="form-group">
					<label class="label" for="player-name">Player Name</label>
					<input type="text" id="player-name" class="input" placeholder="Enter player name" maxlength="100" required>
				</div>

				<div class="form-group">
					<label class="label" for="player-gender">Gender (for handicap guidance)</label>
					<select id="player-gender" class="select">
						<option value="">Select gender</option>
						<option value="male">Male</option>
						<option value="female">Female</option>
						<option value="other">Other</option>
					</select>
				</div>

				<div class="handicap-guide mb-4">
					<div class="handicap-guide-title">Handicap Guidance</div>
					<div class="handicap-suggestions" id="handicap-guide">
						Choose your gender above for handicap suggestions
					</div>
				</div>

				<div class="form-group">
					<label class="label" for="player-handicap">Handicap</label>
					<input type="number" id="player-handicap" class="input" min="0" max="54" step="0.1" placeholder="0.0" required>
				</div>

				<button type="submit" class="btn btn-accent btn-full">Add Player</button>
			</form>
		</div>

		<!-- Current Players -->
		<div class="card mb-6">
			<div class="card-header">
				<h3 class="card-title">Current Players</h3>
			</div>
			<div id="player-list">
				<p class="text-gray-500 text-center py-8">No players added yet</p>
			</div>
		</div>

		<!-- Start Game Button -->
		<div class="text-center">
			<button id="start-game-btn" onclick="window.golfGamez.startGame()"
					class="btn btn-primary btn-lg" disabled>
				Start Game
			</button>
			<p class="text-sm text-gray-500 mt-2">Add at least 1 player to start</p>
		</div>
	</div>
</div>

<style>
.share-link-container {
	animation: fadeIn 0.5s ease-in-out;
}

@keyframes fadeIn {
	from { opacity: 0; transform: translateY(20px); }
	to { opacity: 1; transform: translateY(0); }
}

.player-card {
	animation: slideIn 0.3s ease-out;
}

@keyframes slideIn {
	from { opacity: 0; transform: translateX(-20px); }
	to { opacity: 1; transform: translateX(0); }
}
</style>
`

const gameViewHTML = `
<div class="container">
	<div class="nav">
		<div class="nav-content">
			<div class="nav-brand">Golf Gamez</div>
			<div id="connection-status" class="connection-status connected">
				<span class="status-dot"></span>
				<span>Connected</span>
			</div>
		</div>
	</div>

	<div class="py-4">
		<!-- Game Header -->
		<div class="text-center mb-6">
			<h1 class="text-2xl font-bold text-primary mb-2">Diamond Run</h1>
			<div class="game-status in-progress">
				<span class="status-dot"></span>
				<span>In Progress</span>
			</div>
		</div>

		<!-- Navigation Tabs -->
		<div class="flex mb-6 bg-white rounded-lg p-1 shadow-md">
			<button class="tab-btn active flex-1 py-2 px-4 rounded-md font-medium transition-colors"
					data-tab="scorecard">Scorecard</button>
			<button class="tab-btn flex-1 py-2 px-4 rounded-md font-medium transition-colors"
					data-tab="score-entry">Enter Scores</button>
			<button class="tab-btn flex-1 py-2 px-4 rounded-md font-medium transition-colors"
					data-tab="leaderboard">Leaderboard</button>
			<button class="tab-btn flex-1 py-2 px-4 rounded-md font-medium transition-colors"
					data-tab="side-bets">Side Bets</button>
		</div>

		<!-- Scorecard Tab -->
		<div id="scorecard-tab" class="tab-content active">
			<div id="scorecard-container" class="scorecard">
				<!-- Scorecard content will be dynamically generated -->
				<div class="scorecard-header">
					<h3 class="text-lg font-semibold">Live Scorecard</h3>
				</div>
				<div class="p-4 text-center text-gray-500">
					<p>Scorecard will appear here as scores are entered</p>
				</div>
			</div>
		</div>

		<!-- Score Entry Tab -->
		<div id="score-entry-tab" class="tab-content">
			<div class="card">
				<div class="card-header">
					<h3 class="card-title">Record Score</h3>
					<p class="card-subtitle">Enter stroke and putt counts for each hole</p>
				</div>

				<form id="score-entry-form">
					<div class="form-group">
						<label class="label" for="score-player-select">Player</label>
						<select id="score-player-select" class="select" required>
							<option value="">Select player</option>
							<!-- Player options will be populated dynamically -->
						</select>
					</div>

					<div class="form-group">
						<label class="label" for="score-hole">Hole</label>
						<div class="number-stepper">
							<button type="button" class="minus-btn">-</button>
							<input type="number" id="score-hole" min="1" max="18" value="1" required>
							<button type="button" class="plus-btn">+</button>
						</div>
					</div>

					<div class="score-entry-grid">
						<div class="score-input-group">
							<label class="score-input-label" for="score-strokes">Strokes</label>
							<div class="number-stepper">
								<button type="button" class="minus-btn">-</button>
								<input type="number" id="score-strokes" class="score-input" min="1" max="15" value="4" required>
								<button type="button" class="plus-btn">+</button>
							</div>
						</div>

						<div class="score-input-group">
							<label class="score-input-label" for="score-putts">Putts</label>
							<div class="number-stepper">
								<button type="button" class="minus-btn">-</button>
								<input type="number" id="score-putts" class="score-input" min="0" max="10" value="2" required>
								<button type="button" class="plus-btn">+</button>
							</div>
						</div>
					</div>

					<div class="hole-info mb-4">
						<div class="hole-number">Hole 1</div>
						<div class="par-info">
							<span>Par 4</span>
							<span class="mx-2">•</span>
							<span>350 yards</span>
						</div>
					</div>

					<button type="submit" class="btn btn-primary btn-lg btn-full">
						Record Score
					</button>
				</form>
			</div>
		</div>

		<!-- Leaderboard Tab -->
		<div id="leaderboard-tab" class="tab-content">
			<div id="leaderboard-container" class="leaderboard">
				<div class="leaderboard-header">
					<h3 class="leaderboard-title">Live Leaderboard</h3>
				</div>
				<div class="p-4 text-center text-gray-500">
					<p>Leaderboard will update as scores are entered</p>
				</div>
			</div>
		</div>

		<!-- Side Bets Tab -->
		<div id="side-bets-tab" class="tab-content">
			<!-- Best Nine Side Bet -->
			<div class="side-bet-card mb-4">
				<div class="side-bet-header">
					<h3 class="side-bet-title">Best Nine</h3>
					<span class="side-bet-status">Active</span>
				</div>
				<p class="text-sm text-gray-600 mb-4">
					Best 9 holes vs par with handicap adjustment
				</p>
				<div id="best-nine-standings">
					<p class="text-center text-gray-500 py-4">
						Standings will appear after more scores are entered
					</p>
				</div>
			</div>

			<!-- Putt Putt Poker Side Bet -->
			<div class="side-bet-card">
				<div class="side-bet-header">
					<h3 class="side-bet-title">Putt Putt Poker</h3>
					<span class="side-bet-status">Active</span>
				</div>
				<p class="text-sm text-gray-600 mb-4">
					Earn cards based on putting performance
				</p>
				<div id="putt-poker-status">
					<p class="text-center text-gray-500 py-4">
						Card counts will update as putts are recorded
					</p>
				</div>
			</div>
		</div>
	</div>
</div>

<style>
.tab-btn {
	color: var(--gray-600);
	font-size: 0.875rem;
}

.tab-btn.active {
	background-color: var(--primary-green);
	color: var(--white);
}

.tab-content {
	display: none;
}

.tab-content.active {
	display: block;
	animation: fadeIn 0.3s ease-in-out;
}

@keyframes fadeIn {
	from { opacity: 0; }
	to { opacity: 1; }
}

.score-entry-grid {
	gap: 1rem;
}

@media (max-width: 640px) {
	.tab-btn {
		font-size: 0.75rem;
		padding: 0.5rem 0.25rem;
	}
}
</style>

<script>
// Tab switching functionality
document.querySelectorAll('.tab-btn').forEach(btn => {
	btn.addEventListener('click', () => {
		const tabName = btn.dataset.tab;

		// Remove active class from all tabs and content
		document.querySelectorAll('.tab-btn').forEach(b => b.classList.remove('active'));
		document.querySelectorAll('.tab-content').forEach(c => c.classList.remove('active'));

		// Add active class to clicked tab and corresponding content
		btn.classList.add('active');
		document.getElementById(tabName + '-tab').classList.add('active');
	});
});
</script>
`

const spectatorViewHTML = `
<div class="container">
	<div class="nav">
		<div class="nav-content">
			<div class="nav-brand">Golf Gamez</div>
			<div class="badge badge-secondary">Spectator Mode</div>
		</div>
	</div>

	<div class="py-6">
		<div class="text-center mb-6">
			<h1 class="text-3xl font-bold text-primary mb-2">Diamond Run</h1>
			<div class="game-status in-progress">
				<span class="status-dot"></span>
				<span>Live Game</span>
			</div>
		</div>

		<!-- Live Updates Notice -->
		<div class="alert alert-info mb-6">
			<strong>Live Updates Enabled</strong><br>
			Scores and standings update automatically as players enter their scores.
		</div>

		<!-- Navigation for Spectator -->
		<div class="flex mb-6 bg-white rounded-lg p-1 shadow-md">
			<button class="tab-btn active flex-1 py-2 px-4 rounded-md font-medium"
					data-tab="spectator-leaderboard">Leaderboard</button>
			<button class="tab-btn flex-1 py-2 px-4 rounded-md font-medium"
					data-tab="spectator-scorecard">Scorecard</button>
			<button class="tab-btn flex-1 py-2 px-4 rounded-md font-medium"
					data-tab="spectator-side-bets">Side Bets</button>
		</div>

		<!-- Spectator Leaderboard -->
		<div id="spectator-leaderboard-tab" class="tab-content active">
			<div id="spectator-leaderboard" class="leaderboard">
				<!-- Live leaderboard content -->
			</div>
		</div>

		<!-- Spectator Scorecard -->
		<div id="spectator-scorecard-tab" class="tab-content">
			<div id="spectator-scorecard" class="scorecard">
				<!-- Live scorecard content -->
			</div>
		</div>

		<!-- Spectator Side Bets -->
		<div id="spectator-side-bets-tab" class="tab-content">
			<!-- Side bet standings for spectators -->
		</div>
	</div>
</div>
`

const notFoundHTML = `
<div class="container">
	<div class="nav">
		<div class="nav-content">
			<div class="nav-brand">Golf Gamez</div>
		</div>
	</div>

	<div class="min-h-screen flex items-center justify-center">
		<div class="text-center">
			<h1 class="text-6xl font-bold text-gray-300 mb-4">404</h1>
			<h2 class="text-2xl font-bold text-gray-700 mb-4">Game Not Found</h2>
			<p class="text-gray-500 mb-8">
				The game you're looking for might have ended or the link may be invalid.
			</p>
			<a href="/" onclick="window.golfGamez.navigateTo('/')"
			   class="btn btn-primary">Create New Game</a>
		</div>
	</div>
</div>
`
