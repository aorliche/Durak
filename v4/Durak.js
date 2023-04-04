
$ = q => document.querySelector(q);
$$ = q => [...document.querySelectorAll(q)]

const suits = ['hearts', 'diamonds', 'clubs', 'spades'];
const ranks = ['6', '7', '8', '9', '10', 'jack', 'queen', 'king', 'ace'];

const Suits = ['Hearts', 'Diamonds', 'Clubs', 'Spades'];
const Ranks = ['6', '7', '8', '9', '10', 'Jack', 'Queen', 'King', 'Ace'];

const scale = 0.45;

function getCoords(e) {
    const box = e.target.getBoundingClientRect();
    return [e.clientX - box.left, e.clientY - box.top];
}

function sub(a, b) {
    return [a[0]-b[0], a[1]-b[1]];
}

function suit(i) {
    return suits[Math.floor(i/ranks.length)];
}

function rank(i) {
    return ranks[i%ranks.length];
}

function cardIndexFromObj(obj) {
    const ri = Ranks.indexOf(obj.Rank)
    const si = Suits.indexOf(obj.Suit)
    return si*ranks.length + ri;
}
	
const cardImages = {};
const cardBackImage = new Image;

function loadImages(cb) {
	// Load images
	const numImagesToLoad = 36+1;
	let numImagesLoaded = 0;

	function loadingComplete() {
		return numImagesLoaded === numImagesToLoad;
	} 
    
    function onLoadFn() {
        numImagesLoaded++;	
		if (loadingComplete()) cb();
    }
    
    for (let i=0; i<36; i++) {
        const s = suit(i);
        const r = rank(i);
        cardImages[i] = new Image;
        cardImages[i].addEventListener('load', onLoadFn);
        cardImages[i].src = `cards/fronts/${s}_${r}.png`;
    }
    
	cardBackImage.addEventListener('load', onLoadFn);
	cardBackImage.src = 'cards/backs/astronaut.png';
}

let game;

function newGame() {
    game = new Game();
}

class Board {
    constructor(plays, covers) {
        this.plays = plays ?? [];
        this.covers = covers ?? [];
    }

    draw(ctx) {
        this.layout();
        this.plays.forEach(c => c.draw(ctx, true)); 
        this.covers.forEach(c => c ? c.draw(ctx, true) : null); 
    }

    init(obj) {
        this.plays = obj.Plays.map(c => {
            new Card(cardIndexFromObj(c));
        });
        this.covers = obj.Covers.map(c => {
            new Card(cardIndexFromObj(c));
        });
    }

    layout() {
        if (this.plays.length == 0) {
            return;
        }
        const cx = 400;
        const cy = 250; 
        const cw = scale*cardImages[0].naturalWidth;
        const cd = 10;
        const co = 20;
        const nc = this.covers.reduce((acc, cur) => cur ? acc+1 : acc, 0);
        let len = (cw+cd)*(this.plays.length-1)+co*nc;
        let px = cx-len/2;
        for (let i=0; i<this.plays.length; i++) {
            this.plays[i].x = px;
            this.plays[i].y = cy;
            this.plays[i].theta = 0;
            if (this.covers[i]) {
                this.covers[i].x = px+co;
                this.covers[i].y = cy-10;
                this.covers[i].theta = 0;
                px += cw+co+cd;
            } else {
                px += cw+cd;
            }
        }
    }
}

class Game {
    constructor() {
        this.players = [new Player(0, true), new Player(1, true)]; 
        this.board = new Board();
        fetch('http://10.100.205.6:8080/game')
        .then(resp => resp.json())
        .then(json => {
            console.log(json);
            this.init(json)
        })
        .catch(err => console.log(err));
        this.dragging = null;
        this.pending = false;
        this.initButtons();
    }

    down(e) {
        const [card, area, player] = this.over(e);
        if (card && area == 'hand') {
            card.dragging = true;
            this.dragging = card;
            player.hand.splice(player.hand.indexOf(card), 1);
            this.dragging.offset = sub(getCoords(e), [card.x, card.y]);
            this.dragging.theta = 0;
            this.dragging.player = player;
        }
    }

    draw(ctx) {
        ctx.clearRect(0, 0, 800, 500);
        this.players.forEach(p => {
            p.draw(ctx);
        });
        this.board.draw(ctx);
        this.layout();
        if (this.trump) this.trump.draw(ctx, true);
        if (this.deck) this.deck.draw(ctx, false);
        if (this.dragging) {
            this.dragging.draw(ctx, true);
        }
    }

    hover(e) {
        this.players.forEach(p => {
            p.hand.forEach(c => {
                c.hovering = false;
            });
        });
        const [card, area, player] = this.over(e);
        if (card && area == 'hand') {
            card.hovering = true;
        }
    }
    
    init(json) {
        this.deck = new Card(0);
        this.trump = new Card(cardIndexFromObj(json.Trump));
        this.board.init(json.Board);
        this.players[0].hand = json.Players[0].Hand.map(c => new Card(cardIndexFromObj(c)));
        this.players[1].hand = json.Players[1].Hand.map(c => new Card(cardIndexFromObj(c)));
        this.players[0].actions = json.Actions[0].map(a => new Action(a));
        this.players[1].actions = json.Actions[1].map(a => new Action(a));
    }

    initButtons() {
        for (let i=0; i<2; i++) {
            const passb = $(`#p${i}pass`);
            const pickupb = $(`#p${i}pickup`);

            this.players[i].passb = passb;
            this.players[i].pickupb = pickupb;

            passb.disabled = true;
            pickupb.disabled = true;

            bpass.addEventListener('click', e => {
                e.preventDefault();
                //if (e.target.disabled) return;
                console.log('passed');
            }
        }
    }

    layout() {
        if (this.deck) {
            this.deck.x = 700;
            this.deck.y = 40;
            this.deck.theta = 3.14/2;
        }
        if (this.trump) {
            this.trump.x = 700;
            this.trump.y = 80;
            this.trump.theta = 0;
        }
    }

    move(e) {
        if (!this.dragging) {
            this.hover(e);
        } else {
            const coords = sub(getCoords(e), this.dragging.offset);
            this.dragging.x = coords[0];
            this.dragging.y = coords[1];
        }
    }

    over(e) {
        let card, area, player;
        const [x,y] = getCoords(e);
        if (y > 150 && y < 350) {
            area = 'board';
        }
        this.players.forEach(p => {
            p.hand.forEach(c => {
                if (c.contains(e)) {
                    card = c;
                    area = 'hand';
                    player = p;
                }
            });
        });
        this.board.plays.forEach((c, i) => {
            if (this.board.covers[i]) {
                return;
            }
            if (c.contains(e)) {
                card = c;
                area = 'board';
            }
        });
        return [card, area, player];
    }

    out(e) {
        if (this.dragging) {
            this.dragging.player.hand.push(this.dragging);
            this.dragging.dragging = false;
            this.dragging.hovering = false;
            this.dragging = null;
        }
    }

    up(e) {
        if (!this.dragging) {
            return;
        }
        let taken = false;
        const [card, area, player] = this.over(e);
        const actions = this.dragging.player.actions;
        if (game.pending) {
            // Wait for response for last action
        } else if (card && area == 'board') {
            // this.over checks whether covers is empty
            for (let i=0; i<actions.length; i++) {
                const cur = actions[i];
                if (cur.verb == 'Defend' && cur.card.i == this.dragging.i && cur.cover.i == card.i) {
                    const i = this.board.plays.indexOf(card);
                    this.board.covers[i] = this.dragging;
                    taken = true;
                    cur.take();
                    break;
                }
            }
        } else if (area == 'board') {
            for (let i=0; i<actions.length; i++) {
                const cur = actions[i];
                if (cur.verb == 'Attack' && cur.card.i == this.dragging.i) {
                    this.board.plays.push(this.dragging);
                    this.board.covers.push(null);
                    taken = true;
                    cur.take();
                    break;
                }
            }
        }
        if (!taken) {
            this.dragging.player.hand.push(this.dragging);
        }
        this.dragging.dragging = false;
        this.dragging.hovering = false;
        this.dragging = null;
    }
}

class Action {
    constructor(obj) {
        this.orig = obj;
        this.pidx = obj.PlayerIdx;
        this.verb = obj.Verb;
        this.card = null;
        if (["Attack", "Defend", "Reverse"].includes(obj.Verb)) {
            this.card = new Card(cardIndexFromObj(obj.Card));
        } 
        this.cover = obj.Cover ? new Card(cardIndexFromObj(obj.Cover)) : null;
    }

    rollback() {
        switch (this.verb) {
            case "Attack": {
                game.players[this.pidx].hand.push(this.card);
                game.board.plays.splice(game.board.plays.indexOf(this.card), 1);
                break;
            }
        }
    }

    take() {
        game.pending = true;
        fetch('http://10.100.205.6:8080/action', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(this.orig)
        })
        .then(resp => resp.json())
        .then(json => {
            console.log(json);
            game.pending = false;
            if (!json.Success) {
                this.rollback();
            } else {
                // Update actions for other player
                game.players[1-this.orig.PlayerIdx].updateActions();
            }
            game.players[this.orig.PlayerIdx].actions = json.Actions.map(act => new Action(act));
        })
        .catch(err => console.log(err));
    }
}

class Player {
    constructor(n, show) {
        this.n = n;
        this.show = show;
        this.hand = [];
        this.actions = [];
    }

    draw(ctx) {
        this.layout();
        this.hand.forEach(c => {
            c.draw(ctx, this.show);
        });
    }

    layout() {
        const n = this.hand.length;
        const cx = 400;
        const cy = this.n == 0 ? 500 : 0;
        const tmult = this.n == 0 ? 1 : -1;
        for (let i=0; i<n; i++) {
            const px = cx+(i-(n-1)/2)*40;
            const theta = (i-(n-1)/2)*0.1;
            this.hand[i].x = px;
            this.hand[i].y = cy;
            this.hand[i].theta = theta*tmult;
        }
    }

    updateActions() {
        fetch(`http://10.100.205.6:8080/actions?p=${this.n}`)
        .then(resp => resp.json())
        .then(json => {
            this.actions = json.map(act => new Action(act));
        })
        .catch(err => console.log(err));

    }
}

class Card {
    constructor(i, x, y, theta) {
        this.x = x;
        this.y = y;
        this.theta = theta;
        this.i = i;
        this.hovering = false;
        this.dragging = false;
    }

    contains(e) {
        let inside;
        this.xform((x,y,w,h,m) => {
            const [ex, ey] = getCoords(e);
            const [cx, cy] = [m.a*ex + m.c*ey + m.e, m.b*ex + m.d*ey + m.f];
            inside = cx > 0 && cx < w && cy > 0 && cy < h;
        });
        return inside;
    }

    draw(ctx, show) {
        let img;
        switch (show) {
            case true: img = cardImages[this.i]; break;
            case false: img = cardBackImage; break;
        }
        this.xform((x,y,w,h,m) => {
            ctx.drawImage(img,0,0,w,h);
        });
    }

    xform(cb) {
        const img = cardImages[this.i];
        const w = scale*img.naturalWidth;
        const h = scale*img.naturalHeight;
        const hovoff = this.y < 250 ? 20 : -20;
        const y = this.hovering ? this.y+hovoff : this.y;
		ctx.save();
		ctx.translate(this.x, y);
		ctx.rotate(this.theta);
		ctx.translate(-w/2, -h/2);
        const m = ctx.getTransform().invertSelf(); 
        if (cb) {
            cb(this.x,y,w,h,m);
        }
        ctx.restore();
    }
}

window.addEventListener('load', e => {
    loadImages(newGame);

    canvas = $('#durak-canvas');
    ctx = canvas.getContext('2d');

    canvas.addEventListener('mousemove', e => {
        if (!game) return;
        game.move(e);
        game.draw(ctx);
    });

    canvas.addEventListener('mousedown', e => {
        if (!game) return;
        game.down(e);
    });

    canvas.addEventListener('mouseup', e => {
        if (!game) return;
        game.up(e);
    });

    canvas.addEventListener('mouseout', e => {
        if (!game) return;
        game.out(e);
    });

    /*function anim(ts) {
        if (game) {
            game.draw(ctx);
        }
        setTimeout(requestAnimationFrame(anim), 100);
    }

    anim();*/
});