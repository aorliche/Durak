// Input Messages

class Action {
    constructor (player, verb, id, card, target, defender, query) {
        this.type = 'Action';
        this.player = player;
        this.verb = verb;
        this.id = id;
        this.card = card;
        this.target = target;
        this.defender = defender;
		this.query = query || false;
		this.turn = state.turn;
    }
}

function actionsEqual(a1, a2) {
	return a1.type == a2.type && a1.player == a2.player && a1.verb == a2.verb
		&& a1.card == a2.card && a1.target == a2.target && a1.defender == a2.defender
		&& a1.query == a2.query && a1.turn == a2.turn;
}

class NewGame {
    constructor (nPlayer, deckLength) {
        this.type = 'NewGame';
        this.nPlayer = nPlayer;
        this.deckLength = deckLength;
    }
}

// State

class Results {
    constructor() {
        this.results = getResults();
    }
}

class State {
    constructor () {
        this.players = players;
        this.discard = discard;
        this.deck = deck;
		this.board = board;
        this.actions = actions;
        this.attacker = Id(attacker);
        this.defender = Id(defender);
        this.trump = trump;
        this.turn = turn;
    }
}

// Output Messages

class Error {
    constructor(code, details, action) {
        this.type = 'Error';
        this.code = code;
        this.details = details;
		this.action = action;
    }
}

function strErr(e) {
    return `Error: ${e.code}: ${e.details}`;
}

class Reply {
    constructor(code, action, source) {
        this.type = 'Reply';
        this.code = code;
        this.action = action;
        this.state = new State;
		this.source = source;
    }
}

class Event {
    constructor(code, details) {
        this.type = 'Event';
        this.code = code;
        this.details = details;
    }
}

class Player {
    constructor () {
        this.hand = [];
        this.actions = [];
    }
}

// Explanation

const suits = ['hearts', 'diamonds', 'clubs', 'spades'];
const ranks = ['6', '7', '8', '9', '10', 'jack', 'queen', 'king', 'ace'];

// Helper

function rank(card) {
    return card ? card % 9 : card;
}

function suit(card) {
    return card ? Math.floor(card/9) : card;
}

function isTrump(card, trump) {
    return suit(card) == suit(trump);
}

function beats(card, target, trump) {
    if (isTrump(card, trump)) return isTrump(target, trump) ? rank(card) > rank(target) : true;
    if (isTrump(target, trump)) return false;
    if (suit(card) == suit(target)) return rank(card) > rank(target);
    return false;
}

function getName(card) {
    return `${ranks[rank(card)]} of ${suits[suit(card)]}`;
}

function Id(player) {
    return players.indexOf(player);
}

// Used in UI

function getNumNotCovered(board) {
    return board.reduce((notCovered, pair) => (pair[1] || pair[1] === 0) ? notCovered : notCovered+1, 0);
}

function lastActionWas(player, verb) {
	return player.actions.length > 0 && last(player.actions).verb == verb;
}
