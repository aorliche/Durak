// A random Durak AI
// Attempts to play all legal and illegal moves in a random order
// 1 second delay between getting state change notice and flurry of moves

importScripts('DurakCommon.js', 'util.js');

// State
let myId = 1;
let state;
let active;
let passTimeout;
let pickUpTimeout;
let actionQ = [];
let eventQ = [];
		
setInterval(intervalFunc, 1000);

function intervalFunc() {
	while (eventQ.length > 0) {
		const e = eventQ.pop();
		if (e.action && e.action.player == myId) {
			dequeActionQ(e.action);
		}
		if (e.type == 'Reply') {
			if (e.code != 'success') {
				console.log('unkown Reply code');
				console.log(e);
			} 
		} else if (e.type == 'Error') {
			// ignore, dequed above
		} else if (e.type == 'Event') {
			if (e.code == 'activate') {
				active = e.details;
			} else if (e.code == 'getactionq') {
				console.log(actionQ);
			} else if (e.code == 'newhand') {
				state = e.details;
			} else if (e.code == 'statechange') {
				console.log('got statechange');
				state = e.details;
			} else {
				console.log(`RandomPlayer${myId} doesn't understand Event code ${e.code}`);
			}
		} else {
			console.log(`RandomPlayer${myId} doesn't understand Message type ${e.type}`);
		}
	}
	if (!active) return;
	if (actionQ.length == 0) {
		if (eventQ.length == 0) {
			considerPickUp();
		} 
		if (actionQ.length == 0) {
			randomMove();
		}
	}
}

onmessage = function(e) {
	e = e.data;
	eventQ.push(e);
}

function considerPickUp() {
	if (myId == state.defender && getNumNotCovered(state.board) > 0) {
		actionQ.push(new Action(state.defender,'pickup',null,null,null,state.defender));
		console.log('picking up');
	}
}

function randomMove() {
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
