// Utility

function last(list) {
    return list[list.length-1];
}

function truthy(item) {
    return item || item === 0;
}

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

function rotate(p, angle) {
    let ca = Math.cos(angle);
    let sa = Math.sin(angle);
    return [p[0]*ca - p[1]*sa, p[0]*sa + p[1]*ca];
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

function getRandomId() {
	Math.random().toString(36).substr(2, 9);
}

function shuffleArray(arr) {
	for (let i=0; i<arr.length; i++) {
		const j = Math.floor(Math.random()*arr.length);
		const k = Math.floor(Math.random()*arr.length);
		if (j != k) [arr[j],arr[k]] = [arr[k],arr[j]];
	}
}
