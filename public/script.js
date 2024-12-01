var GobblersUi = /** @class */ (function () {
    function GobblersUi(on_click) {
        for (var p = 0; p < 2; p++) {
            var _loop_1 = function (s) {
                var id = "g-new-" + (p * 3 + s);
                var el = document.getElementById(id);
                if (el !== null) {
                    el.onclick = function () { return on_click(el); };
                }
                else {
                    throw new Error("Can't find new: " + id);
                }
            };
            for (var s = 0; s < 3; s++) {
                _loop_1(s);
            }
        }
        this.board_cells = new Array(9);
        var _loop_2 = function (i) {
            var id = "g-board-" + i;
            var el = document.getElementById(id);
            if (el !== null) {
                this_1.board_cells[i] = el;
                el.onclick = function () { return on_click(el); };
            }
            else {
                throw new Error("Can't find board: " + id);
            }
        };
        var this_1 = this;
        for (var i = 0; i < 9; i++) {
            _loop_2(i);
        }
    }
    GobblersUi.prototype.remove_highlights = function (selected) {
        if (selected !== "") {
            var old_sel = document.getElementById(selected);
            old_sel.classList.remove("g-cell-sel");
            this.board_cells.forEach(function (element) {
                element.classList.remove("g-cell-act");
            });
        }
    };
    GobblersUi.prototype.set_highlights = function (selected, selected_piece) {
        selected.classList.add("g-cell-sel");
        this.board_cells.forEach(function (element) {
            var cell_piece = element.lastElementChild;
            var act = false;
            if (cell_piece == null) {
                act = true;
            }
            else if (cell_piece.classList.contains("g-piece-0")
                && (selected_piece.classList.contains("g-piece-1")
                    || selected_piece.classList.contains("g-piece-2"))) {
                act = true;
            }
            else if (cell_piece.classList.contains("g-piece-1")
                && selected_piece.classList.contains("g-piece-2")) {
                act = true;
            }
            else {
                act = false;
            }
            if (act) {
                element.classList.add("g-cell-act");
            }
        });
    };
    GobblersUi.prototype.do_move = function (from, to) {
        var move_piece = from.lastElementChild;
        move_piece.classList.remove("g-piece-double");
        from.removeChild(move_piece);
        to.appendChild(move_piece);
    };
    return GobblersUi;
}());
var GobblersGame = /** @class */ (function () {
    function GobblersGame() {
        var _this = this;
        var meta_el = document.getElementById("g-game-meta");
        this.meta = JSON.parse(meta_el.innerText);
        var ws_url = "ws://" + window.location.host + "/gobblers/play-" + this.meta.id + "-" + this.meta.state;
        this.conn = new WebSocket(ws_url);
        this.conn.onopen = function () { return _this.setup_conn(); };
        this.ui = new GobblersUi(this.on_click.bind(this));
        this.selected = "";
        this.can_click = false;
    }
    GobblersGame.prototype.get_next = function () {
        var next = "g-piece-";
        if (this.meta.state % 2 == 0) {
            next += "o";
        }
        else {
            next += "x";
        }
        return next;
    };
    GobblersGame.prototype.do_move = function (move) {
        var from = '';
        if (move.new) {
            from = 'g-new-' + (move.size + (this.meta.state % 2) * 3);
        }
        else {
            from = 'g-board-' + move.from;
        }
        var to = 'g-board-' + move.to;
        this.ui.do_move(document.getElementById(from), document.getElementById(to));
    };
    GobblersGame.prototype.setup_conn = function () {
        console.log("set up");
        var game = this;
        this.conn.onerror = function (ev) {
            console.log("Connection error:", ev);
        };
        this.conn.onmessage = function (ev) {
            console.log(ev.data);
            var msg = JSON.parse(ev.data);
            switch (msg.type) {
                case "move":
                    game.do_move(msg.data);
                    game.can_click = false;
                    game.meta.state += 1;
                    break;
                case "turn":
                    game.can_click = true;
                    break;
                case "stop":
                    game.can_click = false;
                    break;
                default:
                    console.log("Wrong msg type");
            }
            console.log(game);
        };
    };
    GobblersGame.prototype.get_move = function (from, to) {
        var piece = from.firstElementChild;
        return {
            new: from.id.includes("new"),
            size: Number(piece.classList[1].slice(-1)),
            from: Number(from.id.slice(-1)),
            to: Number(to.id.slice(-1)),
        };
    };
    GobblersGame.prototype.on_click = function (clicked) {
        if (!this.can_click)
            return;
        var clicked_act = clicked.classList.contains("g-cell-act");
        this.ui.remove_highlights(this.selected);
        if (this.selected != "" && clicked_act) {
            var selected_el = document.getElementById(this.selected);
            var move = this.get_move(selected_el, clicked);
            this.conn.send(JSON.stringify({
                type: "move",
                data: move
            }));
            this.selected = "";
            this.meta.state += 1;
            this.can_click = false;
            this.ui.do_move(selected_el, clicked);
        }
        else {
            var clicked_piece = clicked.lastElementChild;
            // Select cell if right color
            if (clicked.id == this.selected || clicked_piece == null || !clicked_piece.classList.contains(this.get_next())) {
                this.selected = "";
                return;
            }
            this.ui.set_highlights(clicked, clicked_piece);
            this.selected = clicked.id;
        }
    };
    return GobblersGame;
}());
// function gobblersHandleGame() {
//     const game_id = document.getElementById("g-game-id").innerText
//     const ws_url = "ws://" + window.location.host + "/gobblers/play-" + game_id
//     var connection = new WebSocket(ws_url)
//     const onMove = function (from: HTMLElement, to: HTMLElement) {
//         const piece = from.firstElementChild!
//         var m: GobblersMove = {
//             new: from.id.includes("new"),
//             size: Number(piece.classList[1].slice(-1)),
//             from: Number(from.id.slice(-1)),
//             to: Number(to.id.slice(-1)),
//         }
//         console.log("My move", m)
//         connection.send(JSON.stringify({
//             type: "move",
//             data: m,
//         }))
//     }
//     const next = game_id.split(':')[1]
//     var game = new GobblersGame(onMove, GobblersGame.stateToNext(next), false)
//     connection.onmessage = (event) => {
//         console.log(event.data)
//         const msg = JSON.parse(event.data)
//         switch (msg.type) {
//             case "move":
//                 const move: GobblersMove = msg.data
//                 console.log("Server move", move)
//                 game.doMoveEl(document.getElementById(from), document.getElementById(to))
//                 break
//             case "turn":
//                 game.can_click = true
//                 console.log(game)
//                 break
//             case "stop":
//                 game.can_click = false
//                 break
//             default:
//                 console.log("didnt resolve msg")
//         }
//     };
//     connection.onerror = (error) => {
//         console.error("WebSocket error:", error);
//     };
// }
