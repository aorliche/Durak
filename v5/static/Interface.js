window.addEventListener("load", function(){
    const players = ['Human'];
    const flavor = ['Poppins', 'Wiggins', 'Emote', 'Worts', 'Biggly', 'Smyte'];

    $('#name').value = flavor[Math.floor(Math.random() * flavor.length)];
    $('#qs').addEventListener('click', e => {
        if ($('#qs-div').style.display == 'block') {
            $('#qs-div').style.display = 'none';
        } else {
            $('#qs-div').style.display = 'block';
        }
    });

    function updateNumber() {
        $('#number').innerText = `${players.length} Players`;
    }

    function addPlayer(p, P) {
        if (players.length == 4) {
            alert('Only up to 4 players allowed');
            return;
        }
        const player = document.createElement('div');
        const type = document.createElement('div');
        const button = document.createElement('button');
        player.classList.add('player');
        type.classList.add('type');
        type.classList.add(p);
        type.innerText = P;
        button.innerText = 'Delete';
        player.appendChild(type);
        player.appendChild(button);
        $('#players').appendChild(player);
        players.push(P);
        updateNumber();
        button.addEventListener('click', e => {
            $$('#players .player').forEach((p,idx) => {
                if (p == player) {
                    players.splice(idx,1);
                }
            });
            player.remove();
            updateNumber();
        });
    }

    $('#human').addEventListener('click', e => {
        e.preventDefault();
        addPlayer('human', 'Human');
    });
    
    $('#easy').addEventListener('click', e => {
        e.preventDefault();
        addPlayer('easy', 'Easy');
    });
    
    $('#medium').addEventListener('click', e => {
        e.preventDefault();
        addPlayer('medium', 'Medium');
    });

    $('#start').addEventListener('click', e => {
        e.preventDefault();
        if (players.length < 2) {
            alert('Must have at least 2 players');
            return;
        }
        const name = $('#name').value;
        if (!name) {
            alert('Must have a name');
            return;
        }
        if (name == 'Easy' || name == 'Medium') {
            alert('Name must not be a computer difficulty');
            return;
        }
        newGame(players, name);
    })
});