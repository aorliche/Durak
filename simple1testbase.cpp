#include "base.hpp"

int main(void) {
    Card c("6", "Clubs");
    Card c1("6", "Clubs");
    Card c2("7", "Diamonds");
    Card c3("Ace", "Spades");
    Card c4("King", "Diamonds");
    cout << c << endl;
    cout << (c == c1) << endl;
    cout << (c == c2) << endl;
    cout << get<bool>(higher_rank(vector<Object>{c1,c2})) << endl;
    cout << get<bool>(higher_rank(vector<Object>{c2,c1})) << endl;
    cout << get<bool>(higher_rank(vector<Object>{c3,c4})) << endl;
    cout << get<bool>(higher_rank(vector<Object>{c4,c3})) << endl;
    cout << get<bool>(same_suit(vector<Object>{c4,c3})) << endl;    
    cout << get<bool>(same_suit(vector<Object>{c2,c4})) << endl;    
}