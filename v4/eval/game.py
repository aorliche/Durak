'''
Read JSON file containing game data
And recreate game
'''

import json

game = json.load(open('../games/1686772004.durak', 'r'))
print(game[0])
