interface GobblersMove {
    new: boolean
    from: number
    to: number
    size: number
}

class GobblersUi {
    board_cells: Array<HTMLElement>

    constructor(on_click: Function) {
        for (let p = 0; p < 2; p++) {
            for (let s = 0; s < 3; s++) {
                const id = "g-new-" + (p * 3 + s)
                const el = document.getElementById(id)
                if (el !== null) {
                    el.onclick = () => on_click(el)
                } else {
                    throw new Error("Can't find new: " + id)
                }
            }
        }
        this.board_cells = new Array(9)
        for (let i = 0; i < 9; i++) {
            const id = "g-board-" + i
            const el = document.getElementById(id)
            if (el !== null) {
                this.board_cells[i] = el
                el.onclick = () => on_click(el)
            } else {
                throw new Error("Can't find board: " + id)
            }
        }
    }

    remove_highlights(selected: string) {
        if (selected !== "") {
            const old_sel = document.getElementById(selected)
            old_sel.classList.remove("g-cell-sel")
            this.board_cells.forEach(element => {
                element.classList.remove("g-cell-act")
            })
        }
    }

    set_highlights(selected: HTMLElement, selected_piece: Element) {
        selected.classList.add("g-cell-sel")
        this.board_cells.forEach(element => {
            var cell_piece = element.lastElementChild
            var act = false
            if (cell_piece == null) {
                act = true
            } else if (cell_piece.classList.contains("g-piece-0")
                && (selected_piece.classList.contains("g-piece-1")
                    || selected_piece.classList.contains("g-piece-2"))) {
                act = true
            } else if (cell_piece.classList.contains("g-piece-1")
                && selected_piece.classList.contains("g-piece-2")) {
                act = true
            } else {
                act = false
            }
            if (act) {
                element.classList.add("g-cell-act")
            }
        });
    }

    do_move(from: HTMLElement, to: HTMLElement) {
        var move_piece = from.lastElementChild
        move_piece.classList.remove("g-piece-double")
        from.removeChild(move_piece)
        to.appendChild(move_piece)
    }
}

interface GobblersGameMeta {
    id: string
    state: number
}

class GobblersGame {
    ui: GobblersUi
    selected: string
    conn: WebSocket
    meta: GobblersGameMeta
    can_click: boolean

    constructor() {
        const meta_el = document.getElementById("g-game-meta")
        this.meta = JSON.parse(meta_el.innerText)
        const ws_url = "ws://" + window.location.host + "/gobblers/play-" + this.meta.id + "-" + this.meta.state
        this.conn = new WebSocket(ws_url)
        this.conn.onopen = () => this.setup_conn()
        this.ui = new GobblersUi(this.on_click.bind(this))
        this.selected = ""
        this.can_click = false
    }

    get_next(): string {
        var next = "g-piece-"
        if (this.meta.state % 2 == 0) {
            next += "o"
        } else {
            next += "x"
        }
        return next
    }

    do_move(move: GobblersMove) {
        var from = ''
        if (move.new) {
            from = 'g-new-' + (move.size + (this.meta.state%2)*3)
        } else {
            from = 'g-board-' + move.from
        }
        var to = 'g-board-' + move.to
        this.ui.do_move(document.getElementById(from), document.getElementById(to))
    }

    setup_conn() {
        console.log("set up")
        const game = this
        this.conn.onerror = function (ev) {
            console.log("Connection error:", ev)
        }
        this.conn.onmessage = function (ev) {
            console.log(ev.data)
            const msg = JSON.parse(ev.data)
            switch (msg.type) {
                case "move":
                    game.do_move(msg.data)
                    game.can_click = false
                    game.meta.state += 1
                    break
                case "turn":
                    game.can_click = true
                    break
                case "stop":
                    game.can_click = false
                    break
                default:
                    console.log("Wrong msg type")
            }
            console.log(game)
        }
    }

    get_move(from: HTMLElement, to: HTMLElement): GobblersMove {
        const piece = from.firstElementChild!
        return {
            new: from.id.includes("new"),
            size: Number(piece.classList[1].slice(-1)),
            from: Number(from.id.slice(-1)),
            to: Number(to.id.slice(-1)),
        }
    }

    on_click(clicked: HTMLElement) {
        if (!this.can_click) return

        const clicked_act = clicked.classList.contains("g-cell-act")
        this.ui.remove_highlights(this.selected)
        if (this.selected != "" && clicked_act) {
            const selected_el = document.getElementById(this.selected)

            const move = this.get_move(selected_el, clicked)
            this.conn.send(JSON.stringify({
                type: "move",
                data: move
            }))
            this.selected = ""
            this.meta.state += 1
            this.can_click = false

            this.ui.do_move(selected_el, clicked)
        } else {
            const clicked_piece = clicked.lastElementChild
            // Select cell if right color
            if (clicked.id == this.selected || clicked_piece == null || !clicked_piece.classList.contains(this.get_next())) {
                this.selected = ""
                return
            }
            this.ui.set_highlights(clicked, clicked_piece)
            this.selected = clicked.id
        }
    }
}
