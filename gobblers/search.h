#pragma once

typedef struct {
  int layers[2][3];
  int pieces[2][3];
  int player;
} Board;

Board bInit();
void bPrint(Board *b);
// 0:none 1:win 2:loss 3:draw
int bGetState(Board *b);
void bDoNewMove(Board *b, int size, int toPos);
void bDoBoardMove(Board *b, int size, int fromPos, int toPos);
void bUndoNewMove(Board *b, int size, int toPos);
void bUndoBoardMove(Board *b, int size, int fromPos, int toPos);

// look into function on what the score means
int negamax(Board *b, int alpha, int beta, int depth, long long *nodes);
