// Durak game as a JS Worker

importScripts('DurakCommon.js', 'util.js');

// Helper

function rankOnBoard(card) {
    return board.reduce((present, pair) => present || rank(pair[0]) == rank(card) || rank(pair[1]) == rank(card), false);
}

function validPlayer(id) {
    return id >= 0 && id < players.length;
}

function allPass() {
    return players.reduce((ap, player) => (player != defender) ? 
		ap && lastActionWas(player, 'pass') : ap, true);
}

function getDirection() {
    return (Id(attacker) < Id(defender)) ? 1 : -1;
}

function boardAttacksAllSameRank() {
    let first = board[0];
    board.map(pair => {
        if (rank(pair[0]) != rank(first)) return false;
    });
    return true;
}

function getNotCoveredBoardId(target) {
    for (let i=0; i<board.length; i++) {
        let pair = board[i];
        if (pair[0] == target && (!pair[1] || pair[1] == 0)) return i;
    }
    return -1;
}

function getFirstAttacker() {
    const s = suit(trump);
    for (let i=0, j=s*9; i<9; i++, j++) {
		players.map(player => {
			if (player.hand.includes(j)) return player; 
        });
    }
    return players[0];
}

function shuffle(deckLength) {
    const deck = [...Array(36).keys()];
    for (let i=0; i<36; i++) {
        let j = Math.floor(Math.random()*36);
        [deck[i], deck[j]] = [deck[j], deck[i]];
    }
    return deck.slice(0,deckLength);
}

function deal() {
    const dir = getDirection();
    const dealOrder = [];
    for (let i=0, j=players.length+Id(attacker); i<players.length; i++, j-=dir) {
        dealOrder.push(players[j%players.length]);
    }
    dealOrder.map(player => {
        while (deck.length > 0 && player.hand.length < 6) {
            player.hand.push(deck.splice(deck.length-1, 1)[0]);
        }
    });
}

function endHand() {
    // End game
    if (gameOver()) {
        postMessage(new Event('gameover', new Results));
        return;
    }
    // End hand
    let inc, dir = getDirection();
	if (lastActionWas(defender, 'pickup')) {
		inc = 2*dir;
		//console.log(board);
        board.map(pair => defender.hand.push(...pair));
    } else {
		inc = dir;
        board.map(pair => discard.push(...pair));
    }
    deal();
	const aid = Id(attacker);
	// Weird because of how modulus of negative numbers works
    attacker = players[(players.length+aid+inc)%players.length];
    defender = players[(players.length+aid+inc+dir)%players.length];
    turn++;
	board = [];
    players.map(player => player.actions.push({verb:'newhand'}));
    postMessage(new Event('newhand', new State));
    return;
}

function commit(action) {
    action.id = actions.length;
    actions.push(action);
    players[action.player].actions.push(action);
    postMessage(new Reply('success', action));
}

// End game

function getNumWinners() {
    return players.reduce((nWinners, player) => (win(player)) ? nWinners+1 : nWinners, 0); 
}

function draw() {
    return getNumWinners() == players.length;
}

function win(player) {
    return player.hand.length == 0 && deck.length == 0;
}

function loss(player) {
    return getNumWinners() == players.length-1 && !win(player);
}

function gameOver() {
    return getNumWinners() >= players.length-1;
}

function getResults() {
    return players.map(player => win(player) ? 'win' : 'loss');
}

// State

let players;
let deck;
let board;
let discard;
let attacker;
let defender;
let trump;
let turn;
let actions;

onmessage = function(e) {
	e = e.data;
	//console.log(e);
    if (e.type == 'NewGame') {
        if (e.nPlayer < 2 || e.nPlayer > 6) {
            postMessage(new Error('malformed', 'bad number of players'));
            return;
        }
        if (e.deckLength < 0 || e.deckLength > 36) {
            postMessage(new Error('malformed', 'bad deck length'));
            return;
        }
        players = [...Array(e.nPlayer)].map(player => new Player);
        deck = shuffle(e.deckLength);
        discard = [];
		board = [];
        trump = deck[0];
        deal();
        attacker = getFirstAttacker();
        defender = players[(Id(attacker)+1)%players.length];
        turn = 0;
		actions = [];
        postMessage(e); // confirmation
        postMessage(new Event('newhand', new State));
        return;
    }
	if (e.type != 'Action') {
		postMessage(new Error('malformed', 'only Action allowed besides NewGame'));
		return;
	}
    if (gameOver()) {
        postMessage(new Error('illegal', 'game is over', e));
        return;
    }
    if (!validPlayer(e.player)) {
        postMessage(new Error('malformed', 'bad player', e));
        return;
    }
    if (truthy(e.id)) {
        postMessage(new Error('malformed', 'id hasn\'t been assigned yet', e));
        return;
    }
    const p = players[e.player];

	// Debug only
	if (e.verb == 'stack') {
		// find card and remove
		for (let i=0; i<players.length; i++) {
			for (let j=0; j<players[i].hand.length; j++) {
				if (players[i].hand[j] == e.card) {
					players[i].hand.splice(j,1);
					p.hand.push(e.card);
					postMessage(new Reply('success', e));
					return;
				}
			}
		}
		for (let i=0; i<board.length; i++) {
			if (board[i][0] == e.card || board[i][1] == e.card) {
				board.splice(i,1);
				p.hand.push(e.card);
				postMessage(new Reply('success', e));
				return;
			}
		}
		for (let i=0; i<discard.length; i++) {
			if (discard[i] == e.card) {
				discard.splice(i,1);
				p.hand.push(e.card);
				postMessage(new Reply('success', e));
				return;
			}
		}
		for (let i=0; i<deck.length; i++) {
			if (deck[i] == e.card) {
				deck.splice(i,1);
				p.hand.push(e.card);
				postMessage(new Reply('success', e));
				return;
			}
		}
		console.assert(false, "Bad in stack deck in DurakWorker");
		postMessage(new Error('malformed', 'Card does not exist', e));
		return;
	}
	if (e.verb == 'discardDeck') {
		const nToRemove = e.card;
		let startIdx = deck.length-nToRemove;
		if (startIdx < 0) startIdx = 0;
		if (startIdx >= deck.length) return; // no need to send a reply
		discard.push(...deck.splice(startIdx, nToRemove));
		postMessage(new Reply('success', e));
		return;
	}

    if (p != defender) {
        switch (e.verb) {
            case 'play': {
                if (!p.hand.includes(e.card)) {
                    postMessage(new Error('illegal', 'card not in hand', e));
                    return;
                }
                if (e.defender != Id(defender)) {
                    postMessage(new Error('malformed', 'bad defender value', e));
					return;
                }
                if (e.target) {
                    postMessage(new Error('malformed', 'should not have target', e));
                    return;
                }
                if (getNumNotCovered(board) >= defender.hand.length) {
                    postMessage(new Error('illegal', 'cannot play more cards than defender can cover', e));
                    return;
                }
                if ((board.length == 0 && p == attacker) || rankOnBoard(e.card)) {
					board.push([e.card]);
                    p.hand.splice(p.hand.indexOf(e.card),1);
                    commit(e);
                    return;
                } else {
                    if (p != attacker) {
                        postMessage(new Error('illegal', 'rank is not on board and player is not attacker', e));
                        return;
                    } else if (board.length != 0) {
                        postMessage(new Error('illegal', 'rank is not on board and board is not empty', e));
                        return;
                    } else {
                        postMessage(new Error('fatal', 'unexpected condition in play', e));
                        return;
                    }
                }
            }
            case 'pass': {
                /*if (p == attacker && board.length == 0) {
                    postMessage(new Error('illegal', 'cannot pass attacker turn', e));
                    return;
                }*/
				commit(e);
                if (allPass() && (getNumNotCovered(board) == 0 || lastActionWas(defender, 'pickup'))) {
                    endHand();
                }
				return;
            }
            default: {
                postMessage(new Error('malformed', 'bad verb', e));
                return;
            }
        }
    } 
    if (p == defender) {
        if (e.defender != e.player) {
            postMessage(new Error('malformed', 'bad defender value', e));
            return;
        }
		if (lastActionWas(defender, 'pickup')) {
            postMessage(new Error('illegal', 'already picking up', e));
            return;
        }
        switch(e.verb) {
            case 'reverse': {
                if (!p.hand.includes(e.card)) {
                    postMessage(new Error('illegal', 'card not in hand', e));
                    return;
                }
                if (e.target) {
                    postMessage(new Error('malformed', 'should not have target', e));
                    return;
                }
                if (board.length == 0) {
                    postMessage(new Error('illegal', 'cannot reverse on empty board', e));
                    return;
                }
                if (getNumNotCovered(board) != board.length) {
                    postMessage(new Error('illegal', 'board has covers', e));
                    return;
                }
                if (!boardAttacksAllSameRank()) {
                    postMessage(new Error('fatal', 'board not all same rank', e));
                    return;
                }
				if (rank(board[0][0]) != rank(e.card)) {
					postMessage(new Error('illegal', 
						`${getName(e.card)} is not the same rank as ${getName(board[0][0])}`, e));
					return;
				}
				// Just checking legality
				if (e.query == true) {
    				postMessage(new Reply('success', e));
					return
				}
                board.push([e.card]);
                p.hand.splice(p.hand.indexOf(e.card),1);
				[attacker, defender] = [defender, attacker]
                commit(e);
                return;
            }
            case 'cover': {
                if (!p.hand.includes(e.card)) {
                    postMessage(new Error('illegal', 'card not in hand', e));
                    return;
                }
                if (!beats(e.card, e.target, trump)) {
                    postMessage(new Error('illegal', 
						`${getName(e.card)} doesn't beat ${getName(e.target)}`, e));
                    return;
                }
                const boardId = getNotCoveredBoardId(e.target);
                if (boardId == -1) {
                    postMessage(new Error('illegal', 'not on board and uncovered', e));
                    return;
                }
                board[boardId][1] = e.card;
                p.hand.splice(p.hand.indexOf(e.card),1);
				commit(e);
                if (allPass() && (getNumNotCovered(board) == 0 || lastActionWas(defender, 'pickup'))) {
                    endHand();
                }
				return;
            }
            case 'pickup': {
                if (truthy(e.card) || truthy(e.target)) {
                    postMessage(new Error('malformed', 'should be no card or target on pickup', e));
                    return;
                }
                commit(e);
                return;
            }
            default: {
                postMessage(new Error('malformed', 'bad verb', e));
                return;
            }
        }
    } 
}
