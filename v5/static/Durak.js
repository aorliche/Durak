
$ = q => document.querySelector(q);
$$ = q => [...document.querySelectorAll(q)]

const suits = ['hearts', 'diamonds', 'clubs', 'spades'];
const ranks = ['6', '7', '8', '9', '10', 'jack', 'queen', 'king', 'ace'];

const Suits = ['Hearts', 'Diamonds', 'Clubs', 'Spades'];
const Ranks = ['6', '7', '8', '9', '10', 'Jack', 'Queen', 'King', 'Ace'];

const scale = 0.45;

let ip;
let game;

const cardImages = {};
const cardBackImage = new Image;

function drawText(ctx, text, p, color, font, stroke) {
    ctx.save();
    if (font) ctx.font = font;
    const tm = ctx.measureText(text);
    ctx.fillStyle = color;
    if (p.ljust)
        ctx.fillText(text, p.x, p.y);
    else if (p.rjust)
        ctx.fillText(text, p.x-tm.width, p.y);
    else
        ctx.fillText(text, p.x-tm.width/2, p.y);
    if (stroke) {
        ctx.strokeStyle = stroke;
        ctx.lineWidth = 1;
        ctx.strokeText(text, p.x-tm.width/2, p.y);
    }
    ctx.restore();
    return tm;
}

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

function loadImages(cb) {
	// Load images
	const numImagesToLoad = 36+1;
	let numImagesLoaded = 0;

	function loadingComplete() {
		return numImagesLoaded === numImagesToLoad;
	} 
    
    function onLoadFn() {
        numImagesLoaded++;	
		if (loadingComplete() && cb) cb();
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

function newGame(id) {
    if (game) {
        return;
    }
    const numPlayers = $('select[name="durak-num-players"]').selectedIndex+2;
    const numComputers = $('select[name="durak-num-computers"]').selectedIndex;
    const difficulty = $('select[name="difficulty"]').selectedIndex == 0 ? 'Easy' : 'Medium';
    game = new Game(id, difficulty, numPlayers, numComputers);
}

class Board {
    constructor(plays, covers) {
        this.plays = plays ?? [];
        this.covers = covers ?? [];
    }

    draw(ctx) {
        this.layout();
        this.plays.forEach(c => c.draw(ctx)); 
        this.covers.forEach(c => c ? c.draw(ctx) : null); 
    }

    hasUncovered() {
        let n = this.plays.length;
        this.covers.forEach(c => c ? n-- : null);
        return n > 0;
    }

    init(state) {
        this.plays = state.Plays.map(c => {
            return new Card(c);
        });
        this.covers = state.Covers.map(c => {
            return c == -1 ? null : new Card(c);
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

// TODO don't always be player zero 
// TODO Redo this constructor for players and computers
class Game {
    constructor(id, difficulty, numPlayers, numComputers) {
        // You are always player 0 in the client
        // Must remap if necessary when talking to the server
        // join is useful for when id is not -1 for both players
        this.id = id;
        this.join = id == -1 ? false : true;
        this.players = [];
        this.player = -1;
        for (let i=0; i<numPlayers; i++) {
            this.players.push(new Player(i));
        }
        this.board = new Board();
        // Connect to socket
        this.conn = new WebSocket(`ws://${location.host}/ws`);
        this.conn.onopen = () => {
            const msg = {};
            if (this.join) {
                msg.Type = 'Join';
                msg.Game = id;
            } else {
                msg.Type = 'New';
                msg.Players = [];
                for (let i=0; i<numPlayers; i++) {
                    if (i >= numPlayers-numComputers) {
                        msg.Players.push(difficulty);
                    } else {
                        msg.Players.push('Human');
                    }
                }
            }
            this.conn.send(JSON.stringify(msg));
        }
        // Get messages
        this.conn.onmessage = e => {
            const json = JSON.parse(e.data);
            if (this.id == -1) {
                this.id = json.Key;
            }
            if (this.player == -1) {
                for (let i=0; i<numPlayers; i++) {
                    if (json.Actions[i]) {
                        this.player = i;
                        break;
                    }
                }
            }
            console.log(json);
            this.init(json);
            updateKnowledge(json.Memory);
            this.pending = false;
            if (json.Winner != -1) {
                this.winner = json.Winner;
                this.conn.close();
            }
        }
        this.dragging = null;
        this.pending = false;
        this.winner = -1;
        this.initButtons();
    }

    down(e) {
        const [card, area, player] = this.over(e);
        if (card && area == 'hand') {
            if (this.versus != "self" && player != this.players[0]) {
                return;
            }
            card.dragging = true;
            this.dragging = card;
            //player.hand.splice(player.hand.indexOf(card), 1);
            this.dragging.offset = sub(getCoords(e), [card.x, card.y]);
            this.dragging.theta = 0;
            this.dragging.player = player;
        }
    }

    draw(ctx) {
        ctx.clearRect(0, 0, 800, 500);
        this.layout();
        if (this.trump) this.trump.draw(ctx, true);
        if (this.deck) {
            this.deck.draw(ctx);
            drawText(ctx, `${this.decksize}`, {x: 700, y: 60}, 'red', 'bold 48px sans', 'navy');
        }
        this.players.forEach(p => {
            p.draw(ctx);
        });
        this.board.draw(ctx);
        if (this.winner !== -1) {
            let text = "You lose...";
            if ((this.join && this.winner == 1) || (!this.join && this.winner == 0)) {
                text = "You win!";
            }
            drawText(ctx, `${text}`, {x: 400, y: 275}, 'red', 'bold 64px sans', 'navy');
        }
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
            if (this.versus != "self" && player != this.players[0]) {
                return;
            }
            card.hovering = true;
            this.lastHover = e;
        }
    }
    
    init(json) {
        this.deck = new Card(-1);
        this.trump = new Card(json.State.Trump);
        this.update(json);
    }

    initButtons() {
        this.passb = $('#pass');
        this.pickupb = $('#pickup');

        this.passb.disabled = true;
        this.pickupb.disabled = true;

        this.passb.addEventListener('click', e => {
            e.preventDefault();
            try {
                this.players[0].actions.forEach(act => {
                    if (act.verb == PassVerb) {
                        act.take();
                        throw 0;
                    }
                });
            } catch {}
        });

        this.pickupb.addEventListener('click', e => {
            e.preventDefault();
            try {
                this.players[0].actions.forEach(act => {
                    if (act.verb == PickupVerb) {
                        act.take();
                        throw 0;
                    }
                });
            } catch {}
        });
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
            this.lastHover = e;
        }
    }

    over(e) {
        let card, area, player;
        const [x,y] = getCoords(e);
        if (y > 150 && y < 350) {
            area = 'board';
        }
        // Keep both players for debug
        if (!this.dragging) {
            this.players.forEach(p => {
                p.hand.forEach(c => {
                    if (c.contains(e)) {
                        card = c;
                        area = 'hand';
                        player = p;
                    }
                });
            });
        }
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
            //this.dragging.player.hand.push(this.dragging);
            this.dragging.dragging = false;
            this.dragging.hovering = false;
            this.dragging = null;
        }
    }

    up(e) {
        if (!this.dragging) {
            return;
        }
        //let taken = false;
        const [card, area, player] = this.over(e);
        const actions = this.dragging.player.actions;
        if (game.pending || game.winner !== -1) {
            // ... Do nothing
            // Wait for response for last action
        } else if (card && area == 'board') {
            // this.over checks whether covers is empty
            for (let i=0; i<actions.length; i++) {
                const cur = actions[i];
                if (cur.verb == CoverVerb && cur.card.i == this.dragging.i && cur.cover.i == card.i) {
                    const i = this.board.plays.indexOf(card);
                    this.board.covers[i] = this.dragging;
                    //taken = true;
                    cur.take();
                    break;
                }
            }
        } else if (area == 'board') {
            for (let i=0; i<actions.length; i++) {
                const cur = actions[i];
                if (cur.verb == PlayVerb && cur.card.i == this.dragging.i) {
                    this.board.plays.push(this.dragging);
                    this.board.covers.push(null);
                    //taken = true;
                    cur.take();
                    break;
                } else if (cur.verb == ReverseVerb && cur.card.i == this.dragging.i) {
                    this.board.plays.push(this.dragging);
                    this.board.covers.push(null);
                    //taken = true;
                    cur.take();
                    break;
                }
            }
        }
        /*if (!taken) {
            this.dragging.player.hand.push(this.dragging);
        }*/
        this.dragging.dragging = false;
        this.dragging.hovering = false;
        this.dragging = null;
    }

    // Same object as passed to init
    update(info) {
        this.decksize = info.DeckSize;
        if (info.DeckSize <= 1) {
            this.deck = null;
        }
        if (info.DeckSize == 0) {
            this.trump = null;
        }
        this.board.init(info.State);
        // Human is player zero in the game, but some other number on the server
        // delta is this.player
        // Only client player is sent actions
        const covi = this.players[0].getHovering();
        for (let i=0; i<this.players.length; i++) {
            const j = (i+this.player) % this.players.length;
            this.players[i].hand = info.State.Hands[j].map(c => new Card(c));
            if (i == 0) {
                this.players[0].actions = info.Actions[this.player].map(a => new Action(a));
            }
        }
        this.updateButtons();
        this.players[0].setHovering(covi);
    }

    updateButtons() {
        let [pass, pickup] = [false, false];
        this.players[0].actions.forEach(a => {
            if (a.verb == PassVerb) {
                pass = true;
            }
            if (a.verb == PickupVerb) {
                pickup = true;
            }
        });
        this.passb.disabled = !pass;
        this.pickupb.disabled = !pickup;
        if (pass && this.board.hasUncovered()) {
            $('#pickingup').style.display = 'block';
        } else {
            $('#pickingup').style.display = 'none';
        }
    }
}

const [PlayVerb, CoverVerb, ReverseVerb, PassVerb, PickupVerb, DeferVerb] = [0,1,2,3,4,5,6];

class Action {
    constructor(act) {
        this.orig = act;
        this.pidx = game.join ? 1-act.Player : act.Player;
        this.verb = act.Verb;
        this.card = null;
        if ([PlayVerb, CoverVerb, ReverseVerb].includes(act.Verb)) {
            this.card = new Card(act.Card);
        } 
        this.cover = act.Covering != -1 ? new Card(act.Covering) : null;
    }

    take() {
        const msg = {Type: 'Action', Game: game.id, Action: this.orig};
        game.pending = true;
        game.conn.send(JSON.stringify(msg));
    }
}

class Player {
    constructor(n) {
        this.n = n;
        this.hand = [];
        this.actions = [];
    }

    draw(ctx) {
        this.layout();
        this.hand.forEach(c => {
            for (let i=0; i<game.board.plays.length; i++) {
                const c1 = game.board.plays[i];
                const c2 = game.board.covers[i];
                if (c.i == c1.i || (c2 && c2.i == c.i)) {
                    //console.log('skipped');
                    return;
                }
            }
            if (game.dragging && game.dragging.i == c.i) {
                return;
            }
            c.draw(ctx);
        });
    }

    getHovering() {
        let i = -1;
        this.hand.forEach(c => {
            if (c.hovering) {
                i = c.i;
            }
        });
        return i;
    }

    layout() {
        const num = game.players.length;
        let cx, scaling;
        switch (num) {
            case 2: {
                cx = 400;
                scaling = 1;
                break;
            }
            case 3: {
                switch (this.n) {
                    case 0: cx = 400; break;
                    case 1: cx = 267; break;
                    case 2: cx = 534; break;
                }
                scaling = (this.n == 0) ? 1 : 0.5;
                break;
            }
            case 4: {
                switch (this.n) {
                    case 0: case 2: cx = 400; break;
                    case 1: cx = 200; break;
                    case 3: cx = 600; break;
                }
                scaling = (this.n == 0) ? 1 : 0.33;
                break;
            }
        }
        const n = this.hand.length;
        const cy = this.n == 0 ? 500 : 0;
        const tmult = this.n == 0 ? 1 : -1;
        for (let i=0; i<n; i++) {
            if (this.hand[i].dragging) {
                continue;
            }
            const px = cx+(i-(n-1)/2)*40;
            const theta = (i-(n-1)/2)*0.1;
            this.hand[i].x = px;
            this.hand[i].y = cy;
            this.hand[i].theta = theta*tmult;
            this.hand[i].scaling = scaling;
        }
    }

    setHovering(i) {
        this.layout();
        this.hand.forEach(c => {
            if (c.i == i && c.contains(game.lastHover)) {
                c.hovering = true;
            }
        });
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

    draw(ctx) {
        const img = this.i == -1 ? cardBackImage : cardImages[this.i];
        const scaling = this.scaling ?? 1;
        this.xform((x,y,w,h,m) => {
            ctx.drawImage(img,0,0,w*scaling,h*scaling);
        });
    }

    xform(cb) {
        const img = cardBackImage;
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

function updateKnowledge(json) {
    const p0 = $('#p0-hand');
    const p1 = $('#p1-hand');
    const discard = $('#discard');
    p0.innerHTML = '';
    p1.innerHTML = '';
    discard.innerHTML = '';
    for (let i=0; i<json.Sizes[0]; i++) {
        if (i < json.Hands[0].length) {
            const c = json.Hands[0][i];
            p0.appendChild(cardImages[c].cloneNode());
        } else {
            p0.appendChild(cardBackImage.cloneNode());
        }
    }
    for (let i=0; i<json.Sizes[1]; i++) {
        if (i < json.Hands[1].length) {
            const c = json.Hands[1][i];
            p1.appendChild(cardImages[c].cloneNode());
        } else {
            p1.appendChild(cardBackImage.cloneNode());
        }
    }
    const ds = json.Discard;
    for (let i=0; i<36; i++) {
        if (ds.includes(i)) {
            const elt = cardImages[i].cloneNode();
            elt.style.opacity = 1.0;
            discard.append(elt);
        } else {
            discard.append(cardBackImage.cloneNode());
        }
    }
}

window.addEventListener('load', e => {
    loadImages();

    // Separate connection for calling list
    const conn = new WebSocket(`ws://${location.host}/ws`);

    conn.onmessage = e => {
        // List of integer game ids
        const json = JSON.parse(e.data);
        json.sort((a,b) => a-b);

        const select = $('select[name="durak-list-select"]');
        const toAdd = [];
        const games = [...select.options].map(opt => parseInt(opt.value));
        if (game && games.includes(game.id)) {
            for (let i=0; i<select.options.length; i++) {
                const opt = select.options[i];
                if (parseInt(opt.value) == game.id) {
                    select.remove(i);
                    break;
                }
            }
        } 
        for (let i=0; i<select.options.length; i++) {
            const opt = select.options[i];
            if (!json.includes(parseInt(opt.value))) {
                console.log(opt.value);
                select.remove(i--);
            }
        }
        json.forEach(key => {
            if (!games.includes(key) && !(game && game.id == key)) {
                const opt = document.createElement('option');
                opt.value = key;
                opt.innerHTML = `Game ${key}`;
                select.appendChild(opt);
            }
        });
    }

    setInterval(e => {
        if (!conn.readyState == 1) return;
        conn.send(JSON.stringify({Type: 'List'}));
    }, 1000);

    $('#new').addEventListener('click', e => {
        newGame(-1);
    });

    $('#join').addEventListener('click', e => {
        const select = $('select[name="durak-list-select"]');
        const opt = select.options[select.selectedIndex];
        if (!opt) {
            return;
        }
        newGame(parseInt(opt.value));
    });

    $('select[name="durak-num-players"]').addEventListener('change', e => {
        let numPlayers = $('select[name="durak-num-players"]').selectedIndex;
        numPlayers = parseInt(numPlayers)+2;
        let numComputers = $('select[name="durak-num-computers"]').selectedIndex;
        numComputers = parseInt(numComputers);

        if (numComputers >= numPlayers) {
            $('select[name="durak-num-computers"]').selectedIndex = numPlayers-1;
        }
    });

    /*$('#quit').addEventListener('click', e => {
        console.log('Quit');
    });*/

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

    setInterval(() => {
        if (game) {
            game.draw(ctx);
        }
    }, 100);
});
