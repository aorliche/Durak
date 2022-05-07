import math
from functools import reduce
from inspect import getmembers, isfunction
import random
from typing import Any

class IntWrapper:
	def __init__(self, value=None):
		self.value = value if type(value) == int or value is None else value.value

	def __gt__(self, other):
		return self.value > other.value

	def __eq__(self, other):
		return self.value == other.value

	def __lt__(self, other):
		return not (self > other and self == other)

class Rank(IntWrapper):
	def __init__(self, value=None):
		super(Rank, self).__init__(value)

class Suit(IntWrapper):
	def __init__(self, value=None):
		super(Suit, self).__init__(value)

class Card:
	def __init__(self, rank=None, suit=None):
		self.rank = Rank(rank)
		self.suit = Suit(suit)

	def __repr__(self):
		if self != NoCard:
			return str(self.rank.value + 9*self.suit.value)
		else:
			return 'NoCard'

	def __eq__(self, other):
		return type(other) == Card and self.rank == other.rank and self.suit == other.suit

	@staticmethod
	def random():
		return Card(random.randint(0,9), random.randint(0,3))

NoCard = Card()

class Board(list):
	pass

class Discard(list):
	pass

class Pair(list):
	def __init__(self, under, over=NoCard):
		super(Pair, self).__init__([under, over])

class Hand(list):
	pass

class Deck(list):
	pass

class BaseGame:
	pass

class BaseAction:
	pass

class BasePlayer:
	pass

def allSameRank(pile: Hand|Board|Discard|Deck) -> bool:
	if isinstance(pile, Board):
		pile = flatten(pile)
	return all(card.rank == pile[0].rank for card in pile)

def beats(over: Card, under: Card, trump: Card) -> bool:
    if over.suit == under.suit:
        return over.rank > under.rank
    if over.suit == trump.suit:
        return True
    return False

def contains(lst: list, elt: Any) -> bool:
	if type(elt) == type(lst[0]):
		return lst.index(elt) != -1

def count(pile: list) -> int:
    return len(pile)

def getAttacker(game: BaseGame) -> BasePlayer:
    return game.attacker

def getBoard(game: BaseGame) -> Board:
    return game.board

def getCardA(action: BaseAction) -> Card:
    return action.card

def getDefender(game: BaseGame) -> BasePlayer:
    return game.defender

def getDefenderA(action: BaseAction) -> BasePlayer:
    return action.defender

def getDiscard(game: BaseGame) -> Discard:
    return game.discard

def getIndex(lst: list, elt: Any) -> int:
	if type(elt) == type(lst[0]):
		return lst.index(elt)

def getHand(player: BasePlayer) -> Hand:
    return player.hand

def getItem(lst: list, idx: int):
	return lst[idx]

def getNoCard() -> Card:
	return NoCard

def getPlayerA(action: BaseAction) -> BasePlayer:
    return action.player

def getRank(card: Card) -> Rank:
	return card.rank

def getSuit(card: Suit) -> Suit:
	return card.suit

def getTargetA(action: BaseAction) -> Card:
    return action.target

def getTrump(game: BaseGame) -> Card:
    return game.trump

def getVerbA(action: BaseAction) -> str:
    return action.verb

def getUncovered(board: Board) -> list[Card]:
    return [pair[0] for pair in board if pair[1] == NoCard] 

def hasRank(pile: Hand|Board|Discard|Deck, rank: Rank) -> bool:
	if isinstance(pile, Board):
		pile = flatten(pile)
	return any(card.rank == rank for card in pile)

def isPositive(num: int|float) -> bool:
    return num > 0

def isZero(num: int|float) -> bool:
    return eq(num, 0)

def lessThan(less: int|float, greater: int|float) -> bool:
    return less < greater

def makeVerbCheck(verb):
	def verbCheckFn(cand: str) -> bool:
		return cand == verb
	return verbCheckFn

verbChecks = [makeVerbCheck(verb) for verb in ['cover', 'play', 'reverse', 'pickup', 'pass']]

def eq(a, b):
	if type(a) == bool or type(b) == bool:
		return a == b and type(a) == type(b)
	else:
		return a == b

def allowType(obj, typ):
	if typ == Any:
		return True
	elif type(obj) == bool:
		return typ == bool
	else:
		return isinstance(obj, typ)

def flatten(board):
	return [card for pair in board for card in pair if card is not NoCard]

# Rules

class Rule:
	def __init__(self):
		self.expected = []
		self.seqs = []

	def __call__(self, query):
		results = []
		for s in self.seqs:
			s.clean()
			s.bind(query.params)
			results.append(s())
		return results

	def filter(seq):
		pass

class TakeAction(Rule):
	def filter(seq):
		return type(seq.res) == bool

rules = [TakeAction()]

def getFunctions(module, rules=True):
	blacklist = [flatten, allowType, eq, makeVerbCheck, getFunctions, reduce, getmembers, isfunction]
	res = [b for a,b in getmembers(module, isfunction) if b not in blacklist]
	return res + verbChecks
