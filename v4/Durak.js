
let ip;

$ = q => document.querySelector(q);
$$ = q => [...document.querySelectorAll(q)]

const suits = ['hearts', 'diamonds', 'clubs', 'spades'];
const ranks = ['6', '7', '8', '9', '10', 'jack', 'queen', 'king', 'ace'];

const Suits = ['Hearts', 'Diamonds', 'Clubs', 'Spades'];
const Ranks = ['6', '7', '8', '9', '10', 'Jack', 'Queen', 'King', 'Ace'];

const scale = 0.45;

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

function cardIndexFromObj(obj) {
    const ri = Ranks.indexOf(obj.Rank);
    const si = Suits.indexOf(obj.Suit);
    if (ri == -1) {
        return -1;
    }
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

let game;

function newGame(id, computer) {
    if (game) {
        return;
    }
    game = new Game(id, computer);
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

    init(obj) {
        this.plays = obj.Plays.map(c => {
            return new Card(cardIndexFromObj(c));
        });
        this.covers = obj.Covers.map(c => {
            return c == null ? null : new Card(cardIndexFromObj(c));
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
    constructor(id, computer) {
        // You are always player 0 in the client
        // Must remap if necessary when talking to the server
        this.id = id;
        this.join = id == -1 ? false : true;
        this.players = [new Player(0, true), new Player(1, true)]; 
        this.board = new Board();
        fetch(this.join ? `http://${ip}:8080/join?game=${this.id}&p=1` : `http://${ip}:8080/new?computer=${!!computer}`)
        .then(resp => resp.json())
        .then(json => {
            this.id = json.Key;
            this.init(json)
        })
        .catch(err => console.log(err));
        this.dragging = null;
        this.pending = false;
        this.winner = null;
        this.initButtons();
        this.startPoll();
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
        if (this.winner !== null) {
            let text = "You lose...";
            if ((this.join && this.winner == 1) || this.winner == 0) {
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
        this.trump = new Card(cardIndexFromObj(json.Trump));
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
                    if (act.verb == 'Pass') {
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
                    if (act.verb == 'Pickup') {
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

    startPoll() {
        const p = this.join ? 1 : 0;
        this.poll = setInterval(() => {
            fetch(`http://${ip}:8080/info?game=${this.id}&p=${p}`)
            .then(resp => resp.json())
            .then(json => {
                this.update(json)
            })
            .catch(err => console.log(err));
        }, 500);
    }

    stopPoll() {
        if (this.poll) {
            clearInterval(this.poll);
            this.poll = null;
        }
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
        if (game.pending || game.winner !== null) {
            // ... Do nothing
            // Wait for response for last action
        } else if (card && area == 'board') {
            // this.over checks whether covers is empty
            for (let i=0; i<actions.length; i++) {
                const cur = actions[i];
                if (cur.verb == 'Defend' && cur.card.i == this.dragging.i && cur.cover.i == card.i) {
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
                if (cur.verb == 'Attack' && cur.card.i == this.dragging.i) {
                    this.board.plays.push(this.dragging);
                    this.board.covers.push(null);
                    //taken = true;
                    cur.take();
                    break;
                } else if (cur.verb == 'Reverse' && cur.card.i == this.dragging.i) {
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
    update(json) {
        this.decksize = json.Deck;
        if (json.Deck < 2) {
            this.deck = null;
        }
        if (json.Deck < 1) {
            this.trump = null;
        }
        this.board.init(json.Board);
        const [p0, p1] = this.join ? [1, 0] : [0, 1];
        const i = this.players[0].getHovering();
        this.players[0].hand = json.Players[p0].Hand.map(c => new Card(cardIndexFromObj(c)));
        this.players[1].hand = json.Players[p1].Hand.map(c => new Card(cardIndexFromObj(c)));
        this.players[0].actions = json.Actions[p0].map(a => new Action(a));
        this.players[1].actions = json.Actions[p1].map(a => new Action(a));
        /*this.players[0].updateButtons();
        this.players[1].updateButtons();*/
        this.updateButtons();
        if (parseInt(json.Winner) != -1) {
            this.winner = parseInt(json.Winner);
            this.stopPoll();
        }
        this.players[0].setHovering(i);
    }

    updateButtons() {
        let [pass, pickup] = [false, false];
        this.players[0].actions.forEach(a => {
            if (a.verb == 'Pass') {
                pass = true;
            }
            if (a.verb == 'Pickup') {
                pickup = true;
            }
        });
        this.passb.disabled = !pass;
        this.pickupb.disabled = !pickup;
    }
}

class Action {
    constructor(obj) {
        this.orig = obj;
        this.pidx = game.join ? 1-obj.PlayerIdx : obj.PlayerIdx;
        this.verb = obj.Verb;
        this.card = null;
        if (["Attack", "Defend", "Reverse"].includes(obj.Verb)) {
            this.card = new Card(cardIndexFromObj(obj.Card));
        } 
        this.cover = obj.Cover ? new Card(cardIndexFromObj(obj.Cover)) : null;
    }

    take() {
        game.pending = true;
        fetch(`http://${ip}:8080/action?game=${game.id}`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(this.orig)
        })
        .then(resp => resp.json())
        .then(json => {
            game.pending = false;
        })
        .catch(err => console.log(err));
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
                    console.log('skipped');
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
        const n = this.hand.length;
        const cx = 400;
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
        this.xform((x,y,w,h,m) => {
            ctx.drawImage(img,0,0,w,h);
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

window.addEventListener('load', e => {
    loadImages();

    fetch('/Durak.ip')
    .then(resp => resp.json())
    .then(json => {
        ip = json;
    })
    .catch(err => console.log(err));

    setInterval(e => {
        if (!ip) return;
        fetch(`http://${ip}:8080/list`)
        .then(resp => resp.json())
        .then(json => {
            json = json.sort();
            const select = $('#durak-list select');
            const toAdd = [];
            const games = [...select.options].map(opt => parseInt(opt.value));
            if (game && games.includes(game.id)) {
                for (let i=0; i<select.options.length; i++) {
                    const opt = select.options[i];
                    if (parseInt(opt.value) == game.id) {
                        console.log('a');
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
        })
        .catch(err => console.log(err));
    }, 1000);

    $('#new').addEventListener('click', e => {
        newGame(-1);
    });

    $('#join').addEventListener('click', e => {
        const select = $('#durak-list select');
        const opt = select.options[select.selectedIndex];
        if (!opt) {
            return;
        }
        newGame(parseInt(opt.value));
    });

    $('#computer').addEventListener('click', e => {
        newGame(-1, true);
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

    $('#memory').addEventListener('click', e => {
        e.preventDefault();
        fetch(`http://${ip}:8080/memory?game=${game.id}`)
        .then(resp => resp.json())
        .then(json => {
            const text = $('#text');
            text.value = JSON.stringify(json);
        })
        .catch(err => console.log(err));
    });
});
