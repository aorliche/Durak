import math
from functools import reduce
from inspect import getmembers, isfunction
import random

class Card:
	def __init__(self, rank=None, suit=None):
		self.rank = Rank(rank)
		self.suit = Suit(suit)

	def __repr__(self):
		if self != noCard:
			return str(self.rank.value + 9*self.suit.value)
		else:
			return 'NoCard'

	def __eq__(self, other):
		return type(other) == Card and self.rank == other.rank and self.suit == other.suit

	@staticmethod
	def random():
		return Card(random.randint(0,9), random.randint(0,3))

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

def allSameRank(pile):
    return reduce(lambda x,y: getRank(x) == getRank(y), pile, True)

def beats(over: Card, under: Card, trump: Card) -> bool:
    if over.suit == under.suit:
        return over.rank > under.rank
    if over.suit == trump.suit:
        return True
    return False

def contains(lst, elt):
	if type(elt) == type(lst[0]) and type(elt) != list:
		return lst.index(elt) != -1

def count(pile):
    return len(pile)

def getUncovered(board):
    return [pair[0] for pair in board if pair[1] == None] 

# def countBoard(game):
#     return len([card for pair in game.board for card in pair])

# def getUncoveredOnBoard(game):
#     return [pair[0] for pair in game.board if pair[1] == None]

def getAttacker(game):
    return game.attacker

def getBoard(game):
    return game.board

def getCardA(action):
    return action.card

def getDefender(game):
    return game.defender

def getDefenderA(action):
    return action.defender

def getDiscard(game):
    return game.discard

def getIndex(lst, elt):
	if type(elt) == type(lst[0]):
		return lst.index(elt)

def getHand(player):
    return player.hand

def getItem(lst: list, idx: int):
	return lst[idx]

# Special code in driver
def getItemAny(lst: list):
	return lst

noCard = Card()

def getNoCard():
	return noCard

def getPlayer(game, playerIdx):
    return game.players[playerIdx]

def getPlayerA(action):
    return action.player

def getRank(card):
	return card.rank

def getSuit(card):
	return card.suit

def getTargetA(action):
    return action.target

def getTrump(game):
    return game.trump

def getVerbA(action):
    return action.verb

def hasRank(pile, card):
    return not reduce(lambda x: getRank(x) != getRank(card), True)

#def apply(helper, lst):
#    return map(helper, lst)

#def contract(helper, lst, base):
#    return reduce(helper, lst, base)

def isPositive(num):
    return num > 0 and type(num) != type(True)

def isZero(num):
    return eq(num, 0)

def lessThan(less, greater):
    return less < greater and type(less) != type(True) and type(greater) != type(True)

def makeVerbCheck(verb):
    return lambda v: v == verb

verbChecks = [makeVerbCheck(verb) for verb in ['cover', 'play', 'reverse', 'pickup', 'pass']]

def eq(a, b):
	return a == b and type(a) == type(b)

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
	blacklist = [eq, makeVerbCheck, getFunctions, reduce, getmembers, isfunction]
	res = [b for a,b in getmembers(module, isfunction) if b not in blacklist]
	return res + verbChecks
