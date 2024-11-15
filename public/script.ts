class GobblersGame {
    board_cells: Array<HTMLElement>
    selected: string
    onMove: Function
    next: string
    local: boolean
    can_click: boolean

    constructor(onMove, next, local, can_click) {
        this.board_cells = new Array(9)
        this.selected = ""
        this.onMove = onMove
        this.next = next
        this.local = local
        this.can_click = can_click

        for (let p = 0; p < 2; p++) {
            for (let s = 0; s < 3; s++) {
                const id = "g-new-" + (p*3+s)
                const el = document.getElementById(id)
                if (el !== null) {
                    el.onclick = () => this.handle_click(el)
                } else {
                    throw new Error("Can't find new: " + id)
                }
            }
        }
        for (let i = 0; i < 9; i++) {
            const id = "g-board-" + i
            const el = document.getElementById(id)
            if (el !== null) {
                this.board_cells[i] = el
                el.onclick = (ev) => this.handle_click(el)
            } else {
                throw new Error("Can't find cell: " + id)
            }
        }
    }

    handle_click(el: HTMLElement) {
        if (!this.can_click) {
            return
        }
        // Store el act state for move dicision
        const el_act = el.classList.contains("g-cell-act")
        // Remove old ui state
        const old_sel = document.getElementById(this.selected)
        if (old_sel !== null) {
            old_sel.classList.remove("g-cell-sel")
        }
        this.board_cells.forEach(element => {
            element.classList.remove("g-cell-act")
        })
        // Do move if selection and action cell
        if (this.selected != "" && el_act) {
            this.onMove(this.selected, el.id)
            this.doMove(document.getElementById(this.selected), el)
            this.selected = ""
            if (!this.local) {
                this.can_click = false
            }
        } else {
            const piece = el.lastElementChild
            // Select cell if right color
            if (piece == null || el.id == this.selected || !piece.classList.contains(this.next)) {
                this.selected = ""
                return
            }
            this.selected = el.id
            el.classList.add("g-cell-sel")
            this.board_cells.forEach(element => {
                var cell_piece = element.lastElementChild
                var act = false
                if (cell_piece == null) {
                    act = true
                } else if (cell_piece.classList.contains("g-piece-0")
                    && (piece.classList.contains("g-piece-1")
                        || piece.classList.contains("g-piece-2"))) {
                    act = true
                } else if (cell_piece.classList.contains("g-piece-1")
                    && piece.classList.contains("g-piece-2")) {
                    act = true
                } else {
                    act = false
                }
                if (act) {
                    element.classList.add("g-cell-act")
                }
            });
        }
    }

    doMove(from, to) {
        var move_piece = from.lastElementChild
        move_piece.classList.remove("g-piece-double")
        from.removeChild(move_piece)
        to.appendChild(move_piece)
        if (this.next === "g-piece-o") {
            this.next = "g-piece-x"
        } else {
            this.next = "g-piece-o"
        }
    }

    static stateToNext(state: string) {
        return "g-piece-" + state.charAt(1)
    }

    static posToId(pos: string) {
        var prefix = pos.charAt(0)
        var number = pos.charAt(1)
        if (prefix == 'b') {
            return "g-board-" + number
        } else {
            return "g-new-" + number
        }
    }

    static idToPos(id: string) {
        var parts = id.split('-')
        if (parts[1] == "board") {
            return "b" + parts[2]
        } else {
            return "n" + parts[2]
        }
    }
}

function gobblersHandleGame() {
    const game_id = document.getElementById("g-game-id").innerText
    const ws_url = "ws://" + window.location.host + "/gobblers/ws-" + game_id
    var connection = new WebSocket(ws_url)

    const onMove = function(from, to) {
        connection.send("move:" + GobblersGame.idToPos(from) + "-" + GobblersGame.idToPos(to))
    }

    const next = game_id.split(':')[1]
    var game = new GobblersGame(onMove, GobblersGame.stateToNext(next), false, false)

    connection.onmessage = (event) => {
        console.log(event.data)
        const parts = event.data.split(":")
        switch (parts[0]) {
            case "move": // move:n0-b3
                const move_parts = parts[1].split('-')
                const from = GobblersGame.posToId(move_parts[0])
                const to = GobblersGame.posToId(move_parts[1])
                game.doMove(document.getElementById(from), document.getElementById(to))
                break
            case "turn":
                game.can_click = true
                console.log(game)
                break
            case "stop":
                game.can_click = false
        }
    };

    connection.onerror = (error) => {
        console.error("WebSocket error:", error);
    };
}
