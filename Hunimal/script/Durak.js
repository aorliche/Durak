
let newGame;

// For convenience
let state;

// Debug
let stackHand;
let discardDeck;
let getActionQ;

window.addEventListener("load", e => {
	const canvas = document.querySelector("#durak-canvas");
	const viewHands = document.querySelector('#viewHandsCB');
	const viewGrid = document.querySelector('#viewGridCB');
	const viewDeck = document.querySelector('#viewDeckCB');
    const controlOpp = document.querySelector('#controlOpponentsCB');
	const log = document.querySelector('#log');

	let ctx = canvas.getContext('2d');
	ctx.textAlign = 'center';
    
	ctx.save()
	ctx.font = 'bold 30px sans';
	ctx.fillText('Loading...', canvas.width/2, canvas.height/2-20);
	ctx.restore();

	stackHand = function(playerIdx, cardNameOrNumber) {
		if (Number.isInteger(cardNameOrNumber)) {
			worker.postMessage(new Action(playerIdx,'stack',null,
				cardNameOrNumber,null,state.defender));
			return;
		}
		const parts = cardNameOrNumber.split(" ");
		const rank = ranks.indexOf(parts[0].toLowerCase());
		const suit = suits.indexOf(parts[2].toLowerCase());
		if (rank == -1 || suit == -1) {
			//console.log(`Error getting idx for ${cardNameOrNumber}`);
			return;
		}
		worker.postMessage(new Action(playerIdx,'stack',null,
			suit*9+rank,null,null));
	}

	discardDeck = function(n) {
		worker.postMessage(new Action(null,'discardDeck',null,n,null,null));
	}

	getActionQ = function() {
		oppo.postMessage(new Event('getactionq'));
	}

	// Card constructor
	class UICard {
		constructor(idx) {
			this.idx = idx;
		}

		contains(p) {
			let tp = [p[0]-this.p[0], p[1]-this.p[1]];
			let rp = rotate(tp, this.angle);
			return inBoundingBox(rp[0],rp[1],
				-this.w/2,-this.h/2,this.w,this.h);
		}

		// pushAlias(p, w, h, angle) {
			// this.aliases.push([this.p, this.w, this.h, this.angle]);
			// this.p = p;
			// this.w = w;
			// this.h = h;
			// this.angle = angle;
		// }

		// popAlias() {
			// let alias = this.aliases.pop();
			// this.p = alias[0];
			// this.w = alias[1];
			// this.h = alias[2];
			// this.angle = alias[3];
		// }
	
		draw(ctx, show) {
			let img;
			switch (show)  {
				case 'yes': img = cardImages[this.idx]; break;
				case 'no': img = cardBackImage; break;
				case 'blank': img = null; break;
			}
			drawRotated(ctx, img, this.p[0], this.p[1], 
				this.w, this.h, this.angle, 
				ctx => {
					// Drawing cards under the top one in the deck
					if (img === null) {
						ctx.save();
						ctx.fillStyle = 'white';
						ctx.fillRect(0, 0, this.w, this.h);
						ctx.restore();
					}
					// All cards get a border
                    ctx.save();
                    if (this == hovering || this == selected) {
                        ctx.lineWidth = '4';
                        ctx.strokeStyle = '#4f4';
                    } else {
                        ctx.strokeStyle = 'black';
                    }
					ctx.strokeRect(0, 0, this.w, this.h);
                    ctx.restore();
				});
		}

		// get player() {
			// for (let i=0; i<players.length; i++) {
				// let found = false;
				// players[i].hand.map(c => {
					// if (c === this) found = true;
				// });
				// if (found) return players[i];
			// }
		// }

		// get handIdx() {
			// const player = this.player;
			// for (let i=0; i<player.hand.length; i++) {
				// if (player.hand[i] === this) {
					// return i;
				// }
			// }
		// }

		// get isUnmetAttackOnBoard() {
			// for (let i=0; i<board.length; i++) {
				// if (board[i][0] === this && !board[i][1]) return true;
			// }
			// return false;
		// }
	}

	// Button constructor
	class UICheckButton {
		constructor(ctx, text, params, cbs, checked, visible) {
			this.text = text;
			this.params = params;
			this.cbs = cbs;
			this.checked = checked || false;
			this.visible = (visible === false) ? false : true; // if left undefined

			// Calculate width or height
			if (!(this.params.w && this.params.h)) {
				ctx.save();
				this.params.font = this.params.font ? this.params.font : 
					'bold 24px Sans';
				ctx.font = this.params.font;
				const fontSize = 
					parseInt(ctx.font.match(/(\d+)px/)[1]);
				const tm = ctx.measureText(this.text);
				this.params.w = this.params.w ? this.params.h : 
					tm.width + fontSize;
				this.params.h = this.params.h ? this.params.h : 
					1.4*fontSize;
				ctx.restore();
			}
		}

		contains(p) {
			const tp = this.params.p;
			const w = this.params.w;
			const h = this.params.h;
			return inBoundingBox(p[0],p[1],tp[0]-w/2,tp[1]-h/2,w,h);
		}

		hover() {
			const cb = this.cbs.hover;
			if (!cb) return;
			cb(this);
		}

		click() {
			const cb = this.cbs.click;
			if (!cb) return;
			cb(this);
		}

		draw(ctx) {
			if (!this.visible) return;
			const p = this.params.p;
			const w = this.params.w;
			const h = this.params.h;
			ctx.save();
			ctx.fillStyle = this.checked ? 'pink' : '#ccc';
			ctx.strokeStyle = 'black';
			ctx.fillRect(p[0]-w/2, p[1]-h/2, w, h);
			ctx.strokeRect(p[0]-w/2, p[1]-h/2, w, h);
			ctx.fillStyle = 'red';
			ctx.font = this.params.font;
			ctx.fillText(this.text, p[0], p[1]+0.16*h);
			ctx.restore();
		}
	}

	// UI state
	let selected;
	let hovering;
	let lastDiscard;
	let savedReverseAction;
	let savedCoverAction;

	// UI elements
	let passButton = new UICheckButton(ctx, 'Pass', 
		{p: [180, 375], font: 'bold 18px sans'},
		{
			click: b => {
				//if (players[0] === defender) return;
				b.checked = !b.checked;
				b.draw(ctx);
                worker.postMessage(new Action(state.attacker,'pass',null,null,null,state.defender));
			}, 
			hover: b => {
				//if (players[0] === defender) return;
				const checkSav = b.checked;
				b.checked = true;
				b.draw(ctx)
				b.checked = checkSav;
			},
		}, false);

	let pickUpButton = new UICheckButton(ctx, 'Pick up',
		{p: [260, 375], font: 'bold 18px sans'},
		{
			click: b => {
				//if (players[0] !== defender) return;
				b.checked = true;
				b.draw(ctx);
                worker.postMessage(new Action(state.defender,'pickup',null,null,null,state.defender));
			},
			hover: b => {
				//if (players[0] !== defender) return;
				const checkSav = b.checked;
				b.checked = true;
				b.draw(ctx)
				b.checked = checkSav;
			},
		}, false);

	let reverseButton = new UICheckButton(ctx, 'Reverse',
		{p: [275, 250], font: 'bold 18px sans'},
		{
			click: b => {
				if (!b.visible) return;
				console.assert(savedReverseAction);
                worker.postMessage(savedReverseAction);
				savedReverseAction = null;
				savedCoverAction = null;
				b.visible = false;
				repaint();
			},
			hover: b => {
				if (!b.visible) return;
				const checkSav = b.checked;
				b.checked = true;
				b.draw(ctx)
				b.checked = checkSav;
			}
		}, false, false);

	let buttons = [passButton, pickUpButton, reverseButton];

	// Load images
	const cardImages = {};
	const numImagesToLoad = 36+3;
	let numImagesLoaded = 0;

	function loadingComplete() {
		return numImagesLoaded === numImagesToLoad;
	} 
    
    function onLoadFn() {
        numImagesLoaded++;	
		if (loadingComplete()) newGame(2, 36);
    }
    
    for (let i=0; i<36; i++) {
        const s = suits[suit(i)];
        const r = ranks[rank(i)];
        cardImages[i] = new Image;
        cardImages[i].addEventListener('load', onLoadFn);
        cardImages[i].src = `cards/fronts/${s}_${r}.png`;
    }
    
    const cardBackImage = new Image;
	cardBackImage.addEventListener('load', onLoadFn);
	cardBackImage.src = 'cards/backs/astronaut.png';
    
    const swordImage = new Image;
    swordImage.addEventListener('load', onLoadFn);
    swordImage.src = 'images/sword.png';
    
    const shieldImage = new Image;
    shieldImage.addEventListener('load', onLoadFn);
    shieldImage.src = 'images/shield.png';

	viewHands.addEventListener('change', e => repaint());
	viewGrid.addEventListener('change', e => repaint());
	viewDeck.addEventListener('change', e => repaint());
	controlOpp.addEventListener('change', e => {
		oppo.postMessage(new Event('activate', !controlOpp.checked));
	});
    
    // Game
    
    const worker = new Worker('script/DurakWorker.js');
	const oppo = new Worker('script/DurakRandomPlayer.js');
    const uicards = [...Array(36).keys()].map(card => new UICard(card));
    
	function logText(txt, noPrefix) {
        turnPrefix = (!state || noPrefix) ? '' : state.turn + '. ';
		log.innerHTML += turnPrefix + txt + '\n';
	}
    
    newGame = function(nPlayer, deckLength) {
        worker.postMessage(new NewGame(nPlayer, deckLength));
		selected = null;
		hovering = null;
		lastDiscard = null;
		savedReverseAction = null;
		savedCoverAction = null;
		reverseButton.visible = false;
    }

	oppo.onmessage = function(e) {
		e = e.data;
		if (e.type == 'Action') {
			worker.postMessage(e);
		} else {
			console.log(e);
		}
	}

	function informOppo() {
		oppo.postMessage(new Event('statechange', state));
	}
    
    worker.onmessage = function(e) {
        e = e.data;
        //console.log(e);
        switch (e.type) {
            case 'NewGame': {
                state = null;
                results = null;
				lastDiscard = null;
                for (let i=0; i<uicards.length; i++) 
                    uicards[i].discarded = false;
				if (!controlOpp.checked) 
					oppo.postMessage(new Event('activate', !controlOpp.checked));
                break;
            }
            case 'Error': {
				if (e.action && e.action.player == 1) {
					oppo.postMessage(e);
				} else if (e.action && e.action.verb == 'reverse' && e.action.query) {
					if (savedCoverAction) {
						worker.postMessage(savedCoverAction);
					}
					savedReverseAction = null;
					savedCoverAction = null;
					reverseButton.visible = false;
					repaint();
				} else {
					console.log('Strange error');
					console.log(e);
				}
                break;
            }
            case 'Reply': {
				if (e.action.player == 1) {
					oppo.postMessage(e);
                    logText(`Confirm ${e.action.verb} from player ${e.action.player}`);
				}
                if (e.code != 'success') {
                    logText(`Unkown Reply code ${e.code}`);
                    break;
                } else {
					if (e.action.verb == 'stack') {
                    	state = e.state;
						const discard = getUIDiscard();
						if (discard.length > 0) {
							lastDiscard = last(discard);
						}
						informOppo();
						recalc();
						repaint();
						return;
					} else if (e.action.verb == 'discardDeck') {
						state = e.state;
						informOppo();
						recalc();
						repaint();
						return;
					} else if (e.action.verb == 'reverse' && e.action.query == true) {
						reverseButton.visible = true;
						repaint();
						return;
					}
                    logText(`Confirm ${e.action.verb} from player ${e.action.player}`);
					if (e.action.verb == 'pickup') {
						passButton.checked = false;
					}
                    state = e.state;
					hovering = null;
					selected = null;
					reverseButton.visible = false;
					informOppo();
                    recalc();
                    repaint();
                    break;
                }
            }
            case 'Event': {
                switch (e.code) {
                    case 'gameover': {
                        logText('Game over!');
                        results = e.details;
                        for (let i=0; i<results.length; i++) {
                            logText(`Player ${i} ${results[i]}`);
                        }
                        break;
                    }
                    case 'newhand': {
                        state = e.details;
						pickUpButton.checked = false;
						passButton.checked = false;
						oppo.postMessage(new Event("newhand", state));
                        recalc();
                        repaint();
                        break;
                    }
                    default: logText(`Bad Event code ${e.code}`);
                }
                break;
            }
			case 'Action': {
				if (!controlOpp.checked) {
					oppo.postMessage(new Error('illegal', 'being controlled', e));
					return;
				}
				worker.postMessage(e);
				break;
			}
            default: logText(`Bad type ${e.type}`);
        }
    }
    
    function getPlayerUIHand(player) {
        return state.players[player].hand.map(card => uicards.find(uicard => card == uicard.idx));
    }

	function recalcPlayerCardPositions(player) {
        const hand = getPlayerUIHand(player);
		const maxw = 400;
		const nc = hand.length;
		const sz = maxw/nc;
		const cw = (sz < 40) ? Math.floor(sz) : 40;
		const rsz = cw*(nc-1);
		for (let i=0; i<nc; i++) {
			let c = hand[i];
			c.w = 100;
			c.h = 150;
			switch (player) {
				case 0: 
					c.p = [250-rsz/2+i*cw, 500-25];
					c.angle = Math.PI; 
					break;
				case 1:
					c.p = [250+rsz/2-i*cw, 25];
					c.angle = 0;
					break;
			}	
		}
	}
    
    function getUIDeck() {
        return state.deck.map(card => uicards.find(uicard => card == uicard.idx));
    }

	function recalcDeckPositions() {
        const deck = getUIDeck();
		for (let i=1; i<deck.length; i++) {
			const c = deck[i];
			c.p = [380+170, 100-Math.floor(i/3)*2]
			c.w = 100;
			c.h = 150;
			c.angle = Math.PI/2;
		}
        if (deck.length > 0) {
            const c = deck[0];
            c.p = [380+170, 150];
            c.w = 100;
            c.h = 150;
            c.angle = Math.PI;
        }
	}
    
    function getUIBoard() {
        return state.board.map(pair => pair.map(card => uicards.find(uicard => card == uicard.idx)));
    }

	function recalcBoardPositions() {
        const board = getUIBoard();
		const maxw = 350;
		const nc = board.length;
		const nSolo = getNumNotCovered(board);
		const nPair = nc - nSolo;
		const center = (nPair === 0) ? [250,250] : [250,260];
		const factors = [1.55, 1.15];
		const blf = factors[0]*nPair + factors[1]*nSolo;
		const cw = (blf*100 > maxw) ? Math.floor(maxw/blf) : 100;
		const bw = blf*cw;
		let x = center[0]-bw/2+cw/2;
		for (let i=0; i<board.length; i++) {
			board[i][0].p = [x,center[1]];
			board[i][0].w = cw;
			board[i][0].h = cw*3/2;
			board[i][0].angle = 0;
			if (board[i][1]) {
				board[i][1].p = [x+25,center[1]-25];
				board[i][1].w = cw;
				board[i][1].h = cw*3/2;
				board[i][1].angle = 0;
			}
			x += (board[i][1]) ? (factors[0]*cw) : (factors[1]*cw);
		}
	}
    
    function getUIDiscard() {
        return state.discard.map(card => uicards.find(uicard => card == uicard.idx));
    }
    
    function recalcDiscardPositions() {
        getUIDiscard().map(card => discardCard(card));
    }
    
    function discardCard(card) {
        if (card.discarded) return;
        const w = 70; const h = 105;
		const sx = 450
        const sy = 300;
		const dx = 35;
		const dy = 35;
        const rx = 15*Math.random();
        const ry = 15*Math.random();
        let p0 = [sx, sy];
        if (lastDiscard) {
            let prev = lastDiscard.p0;
            if (prev[0] > 600) {
                p0 = [sx, prev[1]+dy];
            } else {
                p0 = [prev[0]+dx, prev[1]];
            }
        }
        let p = [p0[0]+rx, p0[1]+ry];
        card.p = p;
        card.p0 = p0;
        card.w = w;
        card.h = h;
        card.angle = 0;
        card.discarded = true;
		lastDiscard = card;
    }

	function displayHands() {
		for (let i=0; i<state.players.length; i++) {
			const show = (i === 0 || viewHands.checked) ? 'yes' : 'no';
            const hand = getPlayerUIHand(i); 
			hand.map(c => c.draw(ctx, show));
		}
	}

	function drawRotated(ctx, img, x, y, w, h, angle, cb) {
		ctx.save();
		ctx.translate(x,y);
		ctx.rotate(angle);
		ctx.translate(-w/2, -h/2);
		if (img) ctx.drawImage(img,0,0,w,h);
		if (cb) cb(ctx, w, h);
		ctx.restore();
	}

	function displayDeck() {
        const deck = getUIDeck();
		if (deck.length > 0) deck[0].draw(ctx, 'yes');
		for (let i=1; i<deck.length; i++) {
			deck[i].draw(ctx, 'blank');
		}
		if (deck.length > 1) {
			let c = deck[deck.length-1];
			c.draw(ctx, 'no');
			ctx.save();
			ctx.font = 'bold 60px sans';
			ctx.fillStyle = 'red';
			ctx.fillText(deck.length.toString(), c.p[0], c.p[1]+20);
			ctx.strokeStyle = 'black';
			ctx.strokeText(deck.length.toString(), c.p[0], c.p[1]+20);
			ctx.restore();
		}
	}

	function displayDeckCheat() {
        const deck = getUIDeck();
		const nc = deck.length;
		ctx.save();
		ctx.font = 'normal 18px Sans';
		ctx.fillText('Top', 665, 20);
		ctx.restore();
		for (let i=nc-1; i>=0; i--) {
			let c = deck[i];
			let y = 60+(nc-1-i)*20;
			c.pushAlias([665, y], 33, 55, 0);
			c.draw(ctx, 'yes');
			c.popAlias();
		}
	}

	function displayDiscard() {
		getUIDiscard().map(c => c.draw(ctx, 'yes'));
	}

	function displayGrid() {
		ctx.beginPath();
		ctx.save();
		ctx.setLineDash([2,2]);
		ctx.strokeStyle = 'gray';
		ctx.translate(-0.5,-0.5);
		const ws = 50;
		const hs = 50;
		const nw = Math.floor(canvas.width/ws);
		const nh = Math.floor(canvas.height/hs);
		for (let i=0; i<nw; i++) {
			ctx.moveTo(ws*i, 0);
			ctx.lineTo(ws*i, canvas.height);
			ctx.stroke();
		}
		for (let i=0; i<nh; i++) {
			ctx.moveTo(0, hs*i);
			ctx.lineTo(canvas.width, hs*i);
			ctx.stroke();
		}
		ctx.restore();
	}

	function displayBoard() {
		getUIBoard().map(pair => {
			pair[0].draw(ctx, 'yes');
			if (pair[1]) pair[1].draw(ctx, 'yes');
		});
	}
    
    function displayAttackerDefender() {
        [_top,bottom] = (state.attacker == 0) ? [shieldImage,swordImage] : [swordImage,shieldImage];
        drawRotated(ctx, _top, 90, 105, _top.width, _top.height, 0, null);
        drawRotated(ctx, bottom, 90, canvas.height-105, bottom.width, bottom.height, 0, null);
    }

	function getBoardTargetIdx(tgt) {
		for (let i=0; i<board.length; i++) {
			if (board[i][0] === tgt) return i;
		}
		return -1;
	}
    
    function recalc() {
        for (let i=0; i<state.players.length; i++) {
            recalcPlayerCardPositions(i);
        }
        recalcDeckPositions();
        recalcBoardPositions();
        recalcDiscardPositions();
    }

	function repaint() {
		if (!loadingComplete()) return;
        if (!state) return;
		ctx.clearRect(0,0, canvas.width, canvas.height);
		displayBoard();
		displayDiscard();
		displayDeck();
		if (viewDeck.checked) displayDeckCheat();
		displayHands();
		buttons.map(b => b.draw(ctx));
        displayAttackerDefender();
		if (viewGrid.checked) displayGrid();
	}
    
	canvas.addEventListener('mousemove', e => {
		if (!loadingComplete()) return;
        if (!state) return;
		const p = getCursorPosition(canvas, e);
		// Hovering over buttons
		let found = false;
		buttons.map(b => {
			if (b.contains(p)) {
				hovering = b;
				b.hover();
				found = true;
			}
		});
		if (found) return;
		// Hovering over elligible cards
		const cand = [];
		for (let i=0; i<state.players.length; i++) {
			if (i === 0 || viewHands.checked) {
				getPlayerUIHand(i).map(c => cand.push(c));
			}
		}
        const board = getUIBoard();
        for (let i=0; i<board.length; i++) {
            if (!board[i][1]) {
                cand.push(board[i][0]);
            }
        }
		hovering = null;
		cand.map(c => {
			if (c.contains(p)) hovering = c;
		});
		repaint();
	}, false);
        
	function getPlayerFromUICard(uicard) {
		for (let i=0; i<state.players.length; i++) {
			for (let j=0; j<state.players[i].hand.length; j++) {
				if (uicard.idx == state.players[i].hand[j])
					return i;
			}
		}
		return null; // If on board, for instance
	}

	canvas.addEventListener('click', e => {
		if (!loadingComplete()) return;
        if (!state) return;
        if (results != null) return;
		const p = getCursorPosition(canvas, e);
		const oldSelected = selected;
		selected = null;
		if (!hovering || !hovering.contains(p)) return;
		if (hovering instanceof UICheckButton) {
			hovering.click();
			// Button above may be reverse button
			reverseButton.visible = false;
			savedReverseAction = null;
			savedCoverAction = null;
		} else if (hovering instanceof UICard) {
			reverseButton.visible = false;
			savedReverseAction = null;
			savedCoverAction = null;
            const player = getPlayerFromUICard(hovering);
			if (player == null) {
                // defender covering
				if (oldSelected) {
                	worker.postMessage(new Action(state.defender,'cover',null,
						oldSelected.idx,hovering.idx,state.defender));
				}

            } else if (player != state.defender) {
                // attacker
                worker.postMessage(new Action(player,'play',null,
					hovering.idx,null,state.defender));
            } else {
				// defender
                // Try reverse 
				if (!savedReverseAction) {
					savedReverseAction = new Action(player,'reverse',null,
						hovering.idx,null,state.defender,true);
					worker.postMessage(savedReverseAction);
					savedReverseAction.query = false;
					reverseButton.visible = true;
					selected = hovering;
				}
				if (getPlayerFromUICard(hovering) == state.defender) {
					// Autoplay
					const board = getUIBoard();
					const cand = [];
					for (let i=0; i<board.length; i++) {
						if (beats(hovering.idx, board[i][0].idx, state.trump) 
								&& !truthy(board[i][1])) 
							cand.push(board[i][0].idx);
					}
					if (cand.length == 1) {
						savedCoverAction = new Action(player,'cover',null,
							hovering.idx,cand[0],state.defender);
						selected = hovering;
					} else if (cand.length > 1) {
						selected = hovering;
					} else {
						logText(`${getName(hovering.idx)} covers nothing`);
					}
				} else {
					if (oldSelected) {
						worker.postMessage(new Action(player,'cover',null,
							oldSelected.idx,hovering.idx,state.defender));
					}
				}
            }  
        }
	});

	canvas.addEventListener('mouseout', e => {
		if (hovering) {
			hovering = null;
			repaint();
		}
	}, false);
}, false);
