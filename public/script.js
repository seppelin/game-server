var GobblersGame = /** @class */ (function () {
    function GobblersGame(onMove, next, local, can_click) {
        var _this = this;
        this.board_cells = new Array(9);
        this.selected = "";
        this.onMove = onMove;
        this.next = next;
        this.local = local;
        this.can_click = can_click;
        for (var p = 0; p < 2; p++) {
            var _loop_1 = function (s) {
                var id = "g-new-" + (p * 3 + s);
                var el = document.getElementById(id);
                if (el !== null) {
                    el.onclick = function () { return _this.handle_click(el); };
                }
                else {
                    throw new Error("Can't find new: " + id);
                }
            };
            for (var s = 0; s < 3; s++) {
                _loop_1(s);
            }
        }
        var _loop_2 = function (i) {
            var id = "g-board-" + i;
            var el = document.getElementById(id);
            if (el !== null) {
                this_1.board_cells[i] = el;
                el.onclick = function (ev) { return _this.handle_click(el); };
            }
            else {
                throw new Error("Can't find cell: " + id);
            }
        };
        var this_1 = this;
        for (var i = 0; i < 9; i++) {
            _loop_2(i);
        }
    }
    GobblersGame.prototype.handle_click = function (el) {
        if (!this.can_click) {
            return;
        }
        // Store el act state for move dicision
        var el_act = el.classList.contains("g-cell-act");
        // Remove old ui state
        var old_sel = document.getElementById(this.selected);
        if (old_sel !== null) {
            old_sel.classList.remove("g-cell-sel");
        }
        this.board_cells.forEach(function (element) {
            element.classList.remove("g-cell-act");
        });
        // Do move if selection and action cell
        if (this.selected != "" && el_act) {
            this.onMove(this.selected, el.id);
            this.doMove(document.getElementById(this.selected), el);
            this.selected = "";
            if (!this.local) {
                this.can_click = false;
            }
        }
        else {
            var piece_1 = el.lastElementChild;
            // Select cell if right color
            if (piece_1 == null || el.id == this.selected || !piece_1.classList.contains(this.next)) {
                this.selected = "";
                return;
            }
            this.selected = el.id;
            el.classList.add("g-cell-sel");
            this.board_cells.forEach(function (element) {
                var cell_piece = element.lastElementChild;
                var act = false;
                if (cell_piece == null) {
                    act = true;
                }
                else if (cell_piece.classList.contains("g-piece-0")
                    && (piece_1.classList.contains("g-piece-1")
                        || piece_1.classList.contains("g-piece-2"))) {
                    act = true;
                }
                else if (cell_piece.classList.contains("g-piece-1")
                    && piece_1.classList.contains("g-piece-2")) {
                    act = true;
                }
                else {
                    act = false;
                }
                if (act) {
                    element.classList.add("g-cell-act");
                }
            });
        }
    };
    GobblersGame.prototype.doMove = function (from, to) {
        var move_piece = from.lastElementChild;
        move_piece.classList.remove("g-piece-double");
        from.removeChild(move_piece);
        to.appendChild(move_piece);
        if (this.next === "g-piece-o") {
            this.next = "g-piece-x";
        }
        else {
            this.next = "g-piece-o";
        }
    };
    GobblersGame.stateToNext = function (state) {
        return "g-piece-" + state.charAt(1);
    };
    GobblersGame.posToId = function (pos) {
        var prefix = pos.charAt(0);
        var number = pos.charAt(1);
        if (prefix == 'b') {
            return "g-board-" + number;
        }
        else {
            return "g-new-" + number;
        }
    };
    GobblersGame.idToPos = function (id) {
        var parts = id.split('-');
        if (parts[1] == "board") {
            return "b" + parts[2];
        }
        else {
            return "n" + parts[2];
        }
    };
    return GobblersGame;
}());
function gobblersHandleGame() {
    var game_id = document.getElementById("g-game-id").innerText;
    var ws_url = "ws://" + window.location.host + "/gobblers/ws-" + game_id;
    var connection = new WebSocket(ws_url);
    var onMove = function (from, to) {
        connection.send("move:" + GobblersGame.idToPos(from) + "-" + GobblersGame.idToPos(to));
    };
    var next = game_id.split(':')[1];
    var game = new GobblersGame(onMove, GobblersGame.stateToNext(next), false, false);
    connection.onmessage = function (event) {
        console.log(event.data);
        var parts = event.data.split(":");
        switch (parts[0]) {
            case "move": // move:n0-b3
                var move_parts = parts[1].split('-');
                var from = GobblersGame.posToId(move_parts[0]);
                var to = GobblersGame.posToId(move_parts[1]);
                game.doMove(document.getElementById(from), document.getElementById(to));
                break;
            case "turn":
                game.can_click = true;
                console.log(game);
                break;
            case "stop":
                game.can_click = false;
        }
    };
    connection.onerror = function (error) {
        console.error("WebSocket error:", error);
    };
}
