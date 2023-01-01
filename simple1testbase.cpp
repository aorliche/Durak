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
    cout << get<bool>(higher_rank(vector<Property>{c1.id,c2.id})) << endl;
    cout << get<bool>(higher_rank(vector<Property>{c2.id,c1.id})) << endl;
    cout << get<bool>(higher_rank(vector<Property>{c3.id,c4.id})) << endl;
    cout << get<bool>(higher_rank(vector<Property>{c4.id,c3.id})) << endl;
    cout << get<bool>(same_suit(vector<Property>{c4.id,c3.id})) << endl;    
    cout << get<bool>(same_suit(vector<Property>{c2.id,c4.id})) << endl;    
}