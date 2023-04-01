<? set_include_path('/home3/calmprepared/public_html/'); ?>
<!DOCTYPE html>
<html>
<head>
    <title>Durak - Hunimal</title>
    <!-- Generic Hunimal -->
    <? include('generic.php'); ?>
    <!-- Hunimal-specific -->
	<link rel="stylesheet" href="style/Durak.css">
    <script src="script/util.js"></script>
    <script src="script/DurakCommon.js"></script>
	<script src="script/Durak.js"></script>
</head>
<body>
    <h1>Computer Durak</h1>
    <? include('menu.php'); ?>
    <div id="container">
        <? include('navbar.php'); ?>
        <div id="main">
			<div style='width: 920px'>
				<div>Your opponent: <span id="oppoName"></span></div>
				<canvas style='float: left;' id='durak-canvas' width=700px height=500px></canvas>
				<div style='width: 200px; float: right;'>
					<input type='checkbox' id='viewHandsCB'>
					View Opponent's Hand<br>
					<input type='checkbox' id='viewGridCB'>
					View Grid<br>
					<input type='checkbox' id='viewDeckCB'>
					View Deck<br>
					<input type='checkbox' id='controlOpponentsCB'>
					Control Opponents<br>
					<button onclick='newGame(2,36);' id='resetButton'>New Game</button><br>
					<textarea readonly id='log' rows='10' cols='30'></textarea>
				</div>
			</div>
			<div style='clear: both;'></div>
			<p>Card svg images were retrieved from <a href='https://tekeye.uk/playing_cards/svg-playing-cards'>this site</a> and are in the public domain.</p>
        </div>
    </div>
    <? include('footer.php'); ?>
</body>
</html>
