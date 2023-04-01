
let reset;

// Debug
let testWin;
let forceOnBoard;
let drainDeck;

window.addEventListener("load", e => {
	const cardWidth = 234;
	const cardHeight = 333;
	const cardW = 40;
	const cardH = 57;
	const canvas = document.querySelector("#durak-canvas");
	const viewHands = document.querySelector('#viewHandsCB');
	const viewGrid = document.querySelector('#viewGridCB');
	const viewDeck = document.querySelector('#viewDeckCB');
    const controlOpp = document.querySelector('#controlOpponentsCB');
	const log = document.querySelector('#log');

	// Debug
	testWin = function() {
		players[0].hand.splice(1,5);
		deck = [];
		repaint();
	}
    
    forceOnBoard = function(cardNumber, cardSuit) {
        let s,n;
        for (let i=0; i<suits.length; i++) {
            if (suits[i] == cardSuit.toLowerCase()) {
                s = i;
                break;
            }
        }
        for (let i=0; i<numbers.length; i++) {
            if (numbers[i] == cardNumber.toLowerCase()) {
                n = i;
                break;
            }
        }
        const idx = s*9+n;
        // Check player hands
        for (let i=0; i<players.length; i++) {
            for (let j=0; j<players[i].hand.length; j++) {
                if (players[i].hand[j].idx == idx) {
                    const c = players[i].hand[j];
                    players[i].hand.splice(j,1);
                    board.push([c]);
                    logText(`Forced ${getCardName(c)} from player ${i}'s hand to board`);
                    recalcPlayerCardPositions(players[i]);
                    recalcBoardPositions();
                    repaint();
                    return;
                }
            }
        }
        // Check deck
        for (let i=0; i<deck.length; i++) {
            if (deck[i].idx == idx) {
                const c = deck[i];
                deck.splice(i,1);
                board.push([c]);
                logText(`Forced ${getCardName(c)} from deck to board`);
                recalcDeckPositions();
                recalcBoardPositions();
                repaint();
                return;
            }
        }
        // Check discard
        for (let i=0; i<discard.length; i++) {
            if (discard[i].idx == idx) {
                const c = discard[i];
                discard.splice(i,1);
                board.push([c]);
                logText(`Forced ${getCardName(c)} from discard to board`);
                recalcDeckPositions();
                recalcBoardPositions();
                repaint();
                return;
            }
        }
    }
    
    drainDeck = function(numToDrain) {
        let initLength = deck.length;
        while (numToDrain-- > 0 && deck.length > 0) {
            const c = deck[deck.length-1];
            deck.splice(deck.length-1,1);
            discardCard(c);
        }
        logText(`Drained ${initLength-deck.length} cards from deck`);
        recalcDeckPositions();
        repaint();
    }
	
	let ctx = canvas.getContext('2d');
    
	ctx.textAlign = 'center';

	// Card constructor
	class Card {
		constructor(idx, p) {
			this.idx = idx;
			this.aliases = [];
		}

		contains(p) {
			let tp = [p[0]-this.p[0], p[1]-this.p[1]];
			let rp = rotate(tp, this.angle);
			return inBoundingBox(rp[0],rp[1],
				-this.w/2,-this.h/2,this.w,this.h);
		}

		pushAlias(p, w, h, angle) {
			this.aliases.push([this.p, this.w, this.h, this.angle]);
			this.p = p;
			this.w = w;
			this.h = h;
			this.angle = angle;
		}

		popAlias() {
			let alias = this.aliases.pop();
			this.p = alias[0];
			this.w = alias[1];
			this.h = alias[2];
			this.angle = alias[3];
		}
	
		draw(ctx, show) {
			let img;
			switch (show)  {
				case 'yes': img = idxToImg(this.idx, cardImages); break;
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

		get player() {
			for (let i=0; i<players.length; i++) {
				let found = false;
				players[i].hand.map(c => {
					if (c === this) found = true;
				});
				if (found) return players[i];
			}
		}

		get handIdx() {
			const player = this.player;
			for (let i=0; i<player.hand.length; i++) {
				if (player.hand[i] === this) {
					return i;
				}
			}
		}

		get isUnmetAttackOnBoard() {
			for (let i=0; i<board.length; i++) {
				if (board[i][0] === this && !board[i][1]) return true;
			}
			return false;
		}
	}

	// Button constructor
	class CheckButton {
		constructor(ctx, text, params, cbs, checked) {
			this.text = text;
			this.params = params;
			this.cbs = cbs;
			this.checked = checked || false;

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

	// Game state
	let players;
	let deck;
	let board;
	let discard;
	let attacker;
	let defender;
	let direction;
	let trump;

	// UI state
	let selected;
	let hovering;
    let turn;

	// UI elements
	let passButton = new CheckButton(ctx, 'Pass', 
		{p: [180, 375], font: 'bold 18px sans'},
		{
			click: b => {
				//if (players[0] === defender) return;
				b.checked = !b.checked;
				b.draw(ctx);
			}, 
			hover: b => {
				//if (players[0] === defender) return;
				const checkSav = b.checked;
				b.checked = true;
				b.draw(ctx)
				b.checked = checkSav;
			},
		}, false);

	let pickUpButton = new CheckButton(ctx, 'Pick up',
		{p: [260, 375], font: 'bold 18px sans'},
		{
			click: b => {
				//if (players[0] !== defender) return;
				b.checked = true;
				b.draw(ctx);
			},
			hover: b => {
				//if (players[0] !== defender) return;
				const checkSav = b.checked;
				b.checked = true;
				b.draw(ctx)
				b.checked = checkSav;
			},
		}, false);

	let buttons = [passButton, pickUpButton];

	// Deck contents
	const suits = ['hearts', 'diamonds', 'clubs', 'spades'];
	const numbers = ['6', '7', '8', '9', '10', 
		'jack', 'queen', 'king', 'ace'];
	const cardImages = {};

	// Load images
	const numImagesToLoad = 36+3;
	let numImagesLoaded = 0;

	function loadingComplete() {
		return numImagesLoaded === numImagesToLoad;
	} 
    
    function onLoadFn() {
        numImagesLoaded++;	
		if (loadingComplete()) reset();
    }
    
	for (let i=0; i<suits.length; i++) {
		let s = cardImages[suits[i]] = {};
		for (let j=0; j<numbers.length; j++) {
			let n = numbers[j];
			s[n] = new Image;
			s[n].addEventListener('load', onLoadFn, false);
			s[n].src = `cards/fronts/${suits[i]}_${n}.png`;
		}
	}
    
    const cardBackImage = new Image;
	cardBackImage.addEventListener('load', onLoadFn, false);
	cardBackImage.src = 'cards/backs/astronaut.png';
    
    const swordImage = new Image;
    swordImage.addEventListener('load', onLoadFn, false);
    swordImage.src = 'images/sword.png';
    
    const shieldImage = new Image;
    shieldImage.addEventListener('load', onLoadFn, false);
    shieldImage.src = 'images/shield.png';

	function escapeHTML(text) {
		let map = {
			'&': '&amp;',
			'<': '&lt;',
			'>': '&gt;',
			'"': '&quot;',
			"'": '&#039;'
		};
		return text.replace(/[&<>"']/g, m => map[m]);
	}

	function logText(text, omitPrefix) {
        turnPrefix = (omitPrefix) ? '' : turn + '. ';
		log.innerHTML += turnPrefix + text + '\n';
	}

	viewHands.addEventListener('change', e => repaint());
	viewGrid.addEventListener('change', e => repaint());
	viewDeck.addEventListener('change', e => repaint());
	
	reset = function() {
		logText('New game', true);
		if (!loadingComplete()) return;
		players = [
			{hand:[], position:'bottom'}, 
			{hand:[], position:'top'}
		];
		deck = shuffleDeck();
		board = [];
		attacker = players[0];
		defender = players[1];
		direction = 1;
		discard = [];
		trump = deck[0];
		hovering = null;
		selected = null;
        turn = 0;
		passButton.checked = false;
		pickUpButton.checked = false;
		deal();
		recalcDeckPositions();
		repaint();
	}

	ctx.save()
	ctx.font = 'bold 30px sans';
	ctx.fillText('Loading...', canvas.width/2, canvas.height/2-20);
	ctx.restore();

	function shuffleDeck() {
		let deck = [...Array(36).keys()];
		for (let i=0; i<36; i++) {
			let j = Math.floor(Math.random()*36);
			[deck[i], deck[j]] = [deck[j], deck[i]];
		}
		return deck.map(idx => new Card(idx));
	}

	function dealCard(player) {
		let card = deck.splice(deck.length-1, 1)[0];
		card.player = player;
		player.hand.push(card);
	}

	function deal() {
		const dealOrder = [];
		for (let i=0, j=players.length+getPlayerIdx(attacker); 
				i<players.length; 
				i++, j-=direction) {
			dealOrder.push(players[j%players.length]);
		}
		for (let i=0; i<dealOrder.length; i++) {
			if (deck.length === 0) {
				break;
			}
			let hand = dealOrder[i].hand;
			if (hand.length < 6) {
				let td = 6 - hand.length;
				if (td > deck.length) {
					td = deck.length;
				}
				logText(`Dealt ${td} cards to player ${getPlayerIdx(dealOrder[i])}`);
				for (let j=0; j<td; j++) {
					dealCard(dealOrder[i]);
				}
				recalcPlayerCardPositions(dealOrder[i]);
			}
		}
		if (deck.length === 0) {
			logText('Deck is empty');
		}
	}

	function recalcPlayerCardPositions(player) {
		const maxw = 400;
		const nc = player.hand.length;
		const sz = maxw/nc;
		const cw = (sz < 40) ? Math.floor(sz) : 40;
		const rsz = cw*(nc-1);
		for (let i=0; i<nc; i++) {
			let c = player.hand[i];
			c.w = 100;
			c.h = 150;
			switch (player.position) {
				case 'top': 
					c.p = [250+rsz/2-i*cw, 25];
					c.angle = Math.PI; 
					break;
				case 'bottom':
					c.p = [250-rsz/2+i*cw, 500-25];
					c.angle = 0;
					break;
			}	
		}
	}

	function recalcDeckPositions() {
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

	function recalcBoardPositions() {
		const maxw = 350;
		const nc = board.length;
		const nPair = getNumReplies(board);
		const nSolo = nc - nPair;
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

	function getSuit(i) {
		if (i instanceof Card) i = i.idx;
		return Math.floor(i/9);
	}

	function getNumber(i) {
		if (i instanceof Card) i = i.idx;
		return i%9;
	}

	function idxToImg(i) {
		const s = getSuit(i);
		const n = getNumber(i);
		return cardImages[suits[s]][numbers[n]];
	}
    
    function getCardName(i, writeArticle) {
		const s = suits[getSuit(i)];
		const n = numbers[getNumber(i)];
        let article = '';
        if (writeArticle) {
            switch (n) {
                case '8':
                case 'ace':
                    article = 'an ';
                    break;
                default:
                    article = 'a ';
            }
        }
        return `${article}${n} of ${s}`;
    }

	function displayHands() {
		for (let i=0; i<players.length; i++) {
			let show = (i === 0 || viewHands.checked) ? 'yes' : 'no';
			players[i].hand.map(c => c.draw(ctx, show));
		}
	}

	function rotate(p, angle) {
		let ca = Math.cos(angle);
		let sa = Math.sin(angle);
		return [p[0]*ca - p[1]*sa, p[0]*sa + p[1]*ca];
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
		discard.map(c => c.draw(ctx, 'yes'));
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

	function getNumReplies() {
		return board.reduce((n, pair) => n + (pair[1] ? 1 : 0), 0);
	}

	function displayBoard() {
		board.map(pair => {
			pair[0].draw(ctx, 'yes');
			if (pair[1]) pair[1].draw(ctx, 'yes');
		});
	}
    
    function displayAttackerDefender() {
        [_top,bottom] = (attacker == players[0]) ? [shieldImage,swordImage] : [swordImage,shieldImage];
        drawRotated(ctx, _top, 90, 105, _top.width, _top.height, 0, null);
        drawRotated(ctx, bottom, 90, canvas.height-105, bottom.width, bottom.height, 0, null);
    }

	function getBoardTargetIdx(tgt) {
		for (let i=0; i<board.length; i++) {
			if (board[i][0] === tgt) return i;
		}
		return -1;
	}

	function playCard(card, tgt) {
		const playerSav = card.player;
        if (!tgt) {
            logText(`Player ${getPlayerIdx(playerSav)} attacks player ${getPlayerIdx(defender)} with ${getCardName(card)}`);
        } else {
            logText(`Player ${getPlayerIdx(playerSav)} covers ${getCardName(tgt)} with ${getCardName(card)}`);
        }
		card.player.hand.splice(card.handIdx, 1);
		if (playerSav.hand.length === 0 && deck.length === 0) {
			logText(`Player ${getPlayerIdx(playerSav)} wins!`);
			let loserIdx = null;
			for (let i=0; i<players.length; i++) {
				if (players[i].hand.length > 0) {
					if (loserIdx === null) loserIdx = i;
					else {
						loserIdx = null;
						break;
					}
				}
			}
			if (loserIdx !== null) logText(`Player ${loserIdx} loses!`);
		}
		if (!tgt) board.push([card]);
		else {
			const idx = getBoardTargetIdx(tgt);
			board[idx][1] = card;
		}
		recalcBoardPositions();
		recalcPlayerCardPositions(playerSav);
	}

	function repaint() {
		if (!loadingComplete()) return;
		ctx.clearRect(0,0, canvas.width, canvas.height);
		displayHands();
		displayDeck();
		if (viewDeck.checked) displayDeckCheat();
		displayBoard();
		displayDiscard();
		buttons.map(b => b.draw(ctx));
        displayAttackerDefender();
		if (viewGrid.checked) displayGrid();
	}

	function getCursorPosition(canvas, e) {
		const r = canvas.getBoundingClientRect();
		const x = e.clientX - r.left;
		const y = e.clientY - r.top;
		return [x,y];
	}

	function inBoundingBox(mx,my,x,y,w,h) {
		return mx > x && mx < x+w && my > y && my < y+h;
	}

	function canAttackerPlay(card) {
		if (board.length === 0) return true;
        if (board.length-getNumReplies() >= defender.hand.length) {
            logText(`Attackers cannot play more cards than defender has in hand!`);
        }
		let n = getNumber(card);
		for (let i=0; i<board.length; i++) {
			if (n === getNumber(board[i][0])) return true;
			if (board[i][1] && n === getNumber(board[i][1])) return true;
		}
	}

	function isTrump(card) {
		return getSuit(card) === getSuit(trump);
	}

	function getReverseNumber() {
		let cardNum = null;
		for (let i=0; i<board.length; i++) {
			if (board[i][1]) return null;
			if (!cardNum && cardNum !== 0) 
				cardNum = getNumber(board[i][0]);
			else if (cardNum !== getNumber(board[i][0])) return null;
		}
		return cardNum;
	}

	function canDefenderPlay(defense, attack) {
		if (!attack.isUnmetAttackOnBoard) return false;
		if (isTrump(attack)) {
			if (!isTrump(defense)) return false;
			return getNumber(defense) > getNumber(attack);
		} else if (isTrump(defense)) {
			return true;
		} else {
			if (getSuit(attack) !== getSuit(defense)) return false;
			return getNumber(defense) > getNumber(attack);
		}
	}

	function getPotentialDefenderTargets(defense) {
		const cand = [];
		for (let i=0; i<board.length; i++) {
			if (canDefenderPlay(defense, board[i][0])) {
				cand.push(board[i][0]);
			}
		}
		return cand;
	}

	function handFinishedByDefense() {
		if (board.length === 0 || !passButton.checked) return;
		let finished = true;
		board.map(pair => {
			if (!pair[1]) finished = false;
		});
		return finished; 
	}

	function handFinishedBySurrender() {
		return passButton.checked && pickUpButton.checked;
	}

	function getPlayerIdx(player) {
		for (let i=0; i<players.length; i++) {
			if (players[i] === player) return i;
		}
	}

	function getNumCardsOnBoard() {
		return board.reduce((n, pair) => pair[1] ? n+2 : n+1, 0); 
	}

	function reverse() {
		const p1 = getPlayerIdx(defender);
		const p2 = getPlayerIdx(attacker);
		[attacker, defender] = [defender, attacker];
		direction = -direction;
		logText(`Player ${p1} has reversed the attack of ${p2}!`);
	}
    
    function discardCard(card) {
        const w = 70; const h = 105;
		const sx = 450
        const sy = 300;
		const dx = 35;
		const dy = 35;
        const rx = 15*Math.random();
        const ry = 15*Math.random();
        let p0 = [sx, sy];
        if (discard.length > 0) {
            let prev = discard[discard.length-1].p0;
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
        discard.push(card);
    }
    
    function getNumHolding() {
        return players.reduce((numHolding, player) => numHolding + (player.hand.length > 0 ? 1 : 0), 0);
    }
    
    function gameFinished() {
        return deck.length == 0 && getNumHolding() <= 1;
    }
    
    function announceWinner() {
        if (getNumHolding() == 0) {
            logText('Game is drawn!');
        } else {
            for (let i=0; i<players.length; i++) {
                if (player.hand > 0) {
                    logText(`Player ${i} has lost!`);
                } else {
                    logText(`Player ${i} has won!`);
                }
            }
        }
    }
    
	canvas.addEventListener('mousemove', e => {
		if (!loadingComplete()) return;
		const p = getCursorPosition(canvas, e);
		// if (hovering && hovering.contains(p)) return;
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
		for (let i=0; i<players.length; i++) {
			if (i === 0 || viewHands.checked) {
				players[i].hand.map(c => cand.push(c));
			}
		}
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

	canvas.addEventListener('click', e => {
		if (!loadingComplete()) return;
        if (gameFinished()) return;
		const p = getCursorPosition(canvas, e);
		const oldSelected = selected;
		selected = null;
		if (!hovering || !hovering.contains(p)) return;
		if (hovering instanceof Card) {
			if (hovering.player === defender || oldSelected) {
                // Defender
				if (!pickUpButton.checked) {
                    const rCardNum = getReverseNumber();
                    if (rCardNum !== null && hovering.player === defender) {
                        if (rCardNum === getNumber(hovering)) {
                            if (getSuit(hovering) === getSuit(trump)) {
                                selected = hovering;
                            } else {
                                if (board.length+1 > attacker.hand.length) {
                                    logText(`Player ${getPlayerIdx(defender)} cannot play more cards than player ${getPlayerIdx(attacker)} has in hand`);
                                } else {
                                    reverse();
                                    playCard(hovering);
                                }
                            }
                            repaint();
                        }
                    }
                    if (oldSelected) {
						if (canDefenderPlay(oldSelected, hovering)) {
							playCard(oldSelected, hovering);
                            repaint();
						}
					} else {
						const tgts = getPotentialDefenderTargets(hovering);
						if (tgts.length === 1) {
							playCard(hovering, tgts[0]);
							hovering = null;
						} else {
							selected = hovering;
						}
                        repaint();
					}
				}
			} else {
                // Attacker
				if (canAttackerPlay(hovering)) {
					playCard(hovering);
					hovering = null;
					repaint();
				}
			}
		} else if (hovering instanceof CheckButton) {
			hovering.click();
		}
		if (handFinishedByDefense()) {
            if (gameFinished()) {
                announceWinner();
                repaint();
                return;
            }
			logText(`Player ${getPlayerIdx(defender)} successfuly defends`);
			board.map(pair => {
				discardCard(pair[0]);
                discardCard(pair[1]);
			});
			board = [];
            turn += 1;
			deal();
			passButton.checked = false;
			pickUpButton.checked = false;
			[defender, attacker] = [attacker, defender];
			repaint();
		} else if (handFinishedBySurrender()) {
            if (gameFinished()) {
                announceWinner();
                repaint();
                return;
            }
			logText(`Player ${getPlayerIdx(defender)} picks up ${getNumCardsOnBoard()} cards`);
			board.map(pair => {
				defender.hand.push(pair[0]);
                logText(`${getCardName(pair[0])}`, true);
				if (pair[1]) {
                    defender.hand.push(pair[1]);
                    logText(`${getCardName(pair[1])}`, true);
                }
			});
			board = [];
            turn += 1;
			deal();
			recalcPlayerCardPositions(defender);
			passButton.checked = false;
			pickUpButton.checked = false;
			repaint();
		}
	});

	canvas.addEventListener('mouseout', e => {
		if (hovering) {
			hovering = null;
			repaint();
		}
	}, false);
}, false);
