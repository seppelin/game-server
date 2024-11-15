#include "search.h"
#include <stdio.h>

//
// Board
//
Board bInit() {
  return (Board){
      {{0, 0, 0}, {0, 0, 0}},
      {{2, 2, 2}, {2, 2, 2}},
      0,
  };
}

void printLayers(int layer, int layerN) {
  for (int i = 0; i < 9; ++i) {
    char c = '-';
    if (layer & 1)
      c = 'O';
    if (layerN & 1)
      c = 'X';
    if (layer & layerN & 1)
      c = '#';
    printf("%c", c);
    layer >>= 1;
    layerN >>= 1;
  }
}

void bPrint(Board *b) {
  printf("Board %d:\n", b->player);
  for (int size = 2; size >= 0; size--) {
    printf("  ");
    printLayers(b->layers[0][size], b->layers[1][size]);
    printf(" | %d-%d\n", b->pieces[0][size], b->pieces[1][size]);
  }
}

int getTopView(Board *b, int sign) {
  int zro = b->layers[sign][0] & ~b->layers[!sign][1];
  int one = (zro | b->layers[sign][1]) & ~b->layers[!sign][2];
  int two = (one | b->layers[sign][2]) & 0b111111111;
  return two;
}

int winningSpots(Board *b) {
  int view = getTopView(b, b->player);
  int spots = 0;
  spots |= (view << 1) & (view << 2) & 0b100100100; // Right spots horizontal
  spots |= (view >> 1) & (view << 1) & 0b010010010; // Mid
  spots |= (view >> 2) & (view >> 1) & 0b001001001; // Left
  spots |= (view << 3) & (view << 6) & 0b111000000; // Top spots vertical
  spots |= (view >> 3) & (view << 3) & 0b000111000; // Mid
  spots |= (view >> 6) & (view >> 3) & 0b000000111; // Bot
  spots |= (view << 2) & (view << 4) & 0b000000100; // Top spots diag
  spots |= (view >> 2) & (view << 2) & 0b000010000; // Mid
  spots |= (view >> 4) & (view >> 2) & 0b001000000; // Bot
  spots |= (view << 4) & (view << 8) & 0b000000001; // Top spots diag
  spots |= (view >> 4) & (view << 4) & 0b000010000; // Mid
  spots |= (view >> 8) & (view >> 4) & 0b100000000; // Bot  
  return spots;
}

int isLine(int layer) {
  int check = layer & (layer << 1) & (layer << 2) & 0b100100100;
  check |= layer & (layer << 2) & (layer << 4) & 0b001000000;
  check |= layer & (layer << 3) & (layer << 6);
  check |= layer & (layer << 4) & (layer << 8);
  return check != 0;
}

int bGetState(Board *b) {
  int win = isLine(getTopView(b, b->player));
  int loss = isLine(getTopView(b, !b->player));
  return win | (loss << 1);
}

int isLeft(Board *b, int size) {
  return b->pieces[b->player][size] != 0;
}

int isFree(Board *b, int size, int pos) {
  int same = b->layers[0][2] | b->layers[1][2];
  switch (size) {
  case 0:
    same |= b->layers[0][0] | b->layers[1][0];
  case 1:
    same |= b->layers[0][1] | b->layers[1][1];
  }
  return ((1 << pos) & same) == 0;
}

int isMovable(Board *b, int size, int pos) {
  int bigger = 0;
  switch (size) {
  case 0:
    bigger |= b->layers[0][1] | b->layers[1][1];
  case 1:
    bigger |= b->layers[0][2] | b->layers[1][2];
  }
  return ((1 << pos) & b->layers[b->player][size] & (~bigger)) != 0;
}

int isCover(Board *b, int size, int pos) {
  int smaller = 0;
  switch (size) {
  case 1:
    smaller = b->layers[0][0] | b->layers[1][0];
    break;
  case 2:
    smaller = b->layers[0][1] | b->layers[1][1];
    break;
  }
  return ((1 << pos) & smaller) != 0;
}

void bDoBoardMove(Board *b, int size, int fromPos, int toPos) {
  b->layers[b->player][size] ^= 1 << fromPos;
  b->layers[b->player][size] |= 1 << toPos;
  b->player = !b->player;
}

void bDoNewMove(Board *b, int size, int toPos) {
  b->pieces[b->player][size] -= 1;
  b->layers[b->player][size] |= 1 << toPos;
  b->player = !b->player;
}

void bUndoBoardMove(Board *b, int size, int fromPos, int toPos) {
  b->player = !b->player;
  b->layers[b->player][size] ^= 1 << toPos;
  b->layers[b->player][size] |= 1 << fromPos;
}

void bUndoNewMove(Board *b, int size, int toPos) {
  b->player = !b->player;
  b->layers[b->player][size] ^= 1 << toPos;
  b->pieces[b->player][size] += 1;
}

const int maxScore = 10000;
const int minScore = -maxScore;
const int winScore = 1000; // +depth
const int lossScore = -winScore; // -depth
const int tooFarScore = 100; // == depth
const int drawScore = 0; // -depth
int negamax(Board *b, int alpha, int beta, int depth, long long *nodes) {
  *nodes += 1;
  switch (bGetState(b)) {
  case 1:
    return winScore + depth;
  case 2:
    return lossScore - depth;
  case 3:
    return drawScore + depth;
  }

  if (depth == 0)
    return tooFarScore;

  int bestPossible = winScore + depth - 1;
  if (beta > bestPossible) {
    beta = bestPossible;
    if (alpha >= beta)
      return alpha;
  }

  // Winning
  int spots = winningSpots(b);
  for (int to = 0; to < 9; to++) {
    if ((spots & (1<<to)) == 0)
      continue;

    for (int size = 0; size < 3; size++) {
      if (!isFree(b, size, to))
        continue;

      if (isLeft(b, size)) {
        bDoNewMove(b, size, to);
        int score = negamax(b, -beta, -alpha, depth - 1, nodes);
        bUndoNewMove(b, size, to);
        if (score != tooFarScore)
          score = -score;
        if (score >= beta)
          return score;
        if (score > alpha)
          alpha = score;
      }

      for (int from = 0; from < 9; from++) {
        if (!isMovable(b, size, from))
          continue;
        
        bDoBoardMove(b, size, from, to);
        int score = negamax(b, -beta, -alpha, depth - 1, nodes);
        bUndoBoardMove(b, size, from, to);
        if (score != tooFarScore)
          score = -score;
        if (score >= beta)
          return score;
        if (score > alpha)
          alpha = score;
      }
    }
  }

  // New cover
  for (int size = 2; size > 0; size--) { // size 0 can't cover
    if (!isLeft(b, size))
      continue;
    for (int toPos = 0; toPos < 9; toPos++) {
      if (!isFree(b, size, toPos) || !isCover(b, size, toPos))
        continue;
      bDoNewMove(b, size, toPos);
      int score = negamax(b, -beta, -alpha, depth - 1, nodes);
      bUndoNewMove(b, size, toPos);

      if (score != tooFarScore)
        score = -score;
      if (score >= beta)
        return score;
      if (score > alpha)
        alpha = score;
    }
  }

  // New !cover
  for (int size = 2; size > -1; size--) {
    if (!isLeft(b, size))
      continue;
    for (int toPos = 0; toPos < 9; toPos++) {
      if (!isFree(b, size, toPos) || isCover(b, size, toPos))
        continue;
      bDoNewMove(b, size, toPos);
      int score = negamax(b, -beta, -alpha, depth - 1, nodes);
      bUndoNewMove(b, size, toPos);

      if (score != tooFarScore)
        score = -score;
      if (score >= beta)
        return score;
      if (score > alpha)
        alpha = score;
    }
  }

  // Board
  for (int size = 2; size > -1; size--) {
    for (int toPos = 0; toPos < 9; toPos++) {
      if (!isFree(b, size, toPos))
        continue;
      for (int fromPos = 0; fromPos < 9; fromPos++) {
        if (!isMovable(b, size, fromPos))
          continue;
        bDoBoardMove(b, size, fromPos, toPos);
        int score = negamax(b, -beta, -alpha, depth - 1, nodes);
        bUndoBoardMove(b, size, fromPos, toPos);

        if (score != tooFarScore)
          score = -score;
        if (score >= beta)
          return score;
        if (score > alpha)
          alpha = score;
      }
    }
  }

  return alpha;
}
