// A random Durak AI
// Attempts to play all legal and illegal moves in a random order
// 1 second delay between getting state change notice and flurry of moves

importScripts('DurakCommon.js', 'util.js');

// State
let name = getRandomId();
let myId = 1;
let state;
let active;
let passTimeout;
let pickUpTimeout;
let actionQ = [];

onmessage = function(e) {
	e = e.data;
	//console.log(e);
	switch (e.type) {
		case 'Event': {
			if (e.code == 'activate') {
				active = e.details;
				if (active) {
					setTimeout(doRandomMove, 1000);
					considerPickUp();
				}
			} else if (e.code == 'getactionq') {
				console.log(actionQ);
			} else if (e.code == 'newhand') {
				state = e.details;
				if (myId == state.attacker) {
					setTimeout(doRandomMove, 1000);
				}
			} else if (e.code == 'statechange') {
				state = e.details;
				setTimeout(doRandomMove, 1000);
				considerPickUp();
			} else {
				postMessage(new Error("malformed", 
					`RandomPlayer_${name} doesn't understand code ${e.code}`));
			}
			return;
		}
		case 'Reply': {
			dequeActionQ(e.action);
			if (!active) return;
			if (e.action.turn != state.turn) return; // old actions
			considerPickUp();
			if (e.code == 'success') {
				state = e.state;
				if (passTimeout) {
					clearTimeout(passTimeout);
					passTimeout = null;
				}
				if (pickUpTimeout) {
					clearTimeout(pickUpTimeout);
					pickUpTimeout = null;
				}
				if (myId != state.defender) {
					setTimeout(doRandomMove, 1000);
				}
			}
			if (!passTimeout && e.action.verb != 'pass' && actionQ.length == 0) {
				passTimeout = setTimeout(e => {
					if (actionQ.length == 0) {
						actionQ.push(new Action(myId,'pass',null,null,null,state.defender));
						postMessage(last(actionQ));
					}
				}, 2000);
			}
			return;
		}
		case 'Error': {
			if (e.action && e.action.player == myId) {
				dequeActionQ(e.action);
			}
			return;
		}
		default: {
			console.log(`RandomPlayer_${name} doesn't understand type ${e.type}`);
			return;
		}
	}
}

function considerPickUp() {
	console.log('entered considering pickup');
	if (pickUpTimeout) return;
	if (!active) return;
	if (myId == state.defender && getNumNotCovered(state.board) > 0) {
		pickUpTimeout = setTimeout(e => {
			if (myId == state.defender && getNumNotCovered(state.board) > 0) {
				actionQ.push(new Action(state.defender,'pickup',null,null,null,state.defender));
				postMessage(last(actionQ));
				console.log('pick up ok');
			}
		}, 2000);
		console.log('set pickup timeout');
	}
}

function doRandomMove() {
	if (!active) return;
	const q = [];
	if (myId == state.defender) {
		// Random chance of just picking up
		/*if (Math.random() < 0.2) {
        	q.push(new Action(state.defender,'pickup',null,null,null,state.defender));
		}*/
		state.players[myId].hand.map(card => {
			// Try every reversal
			q.push(new Action(myId,'reverse',null,
				card,null,state.defender));
			// Try every cover
			state.board.map(pair => {
				if (!truthy(pair[1])) {
					q.push(new Action(myId,'cover',null,
						card,pair[0],state.defender));
				}
			});
		});
	} else {
		// Try every play
		state.players[myId].hand.map(card => {
			q.push(new Action(myId,'play',null,
				card,null,state.defender));
		});
	}
	shuffleArray(q);
	q.map(action => {
		postMessage(action);
	});
	actionQ.push(...q);
}

function dequeActionQ(action) {
	let found = false;
	for (let i=0; i<actionQ.length; i++) {
		if (actionsEqual(actionQ[i], action)) {
			actionQ.splice(i--,1);
			found = true;
			break;
		}
	}
	if (!found) {
		console.log('Spurious action');
		console.log(action);
	}
}
