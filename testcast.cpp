struct MyType {
    int id;
    MyType(int i) : id(i) {}
};

struct MyDerivedType : public MyType {
    long poop;
    MyDerivedType(const MyType &t) : MyType(t), poop(3) {}
    MyDerivedType() : MyType(0), poop(2) {}
};

#include <iostream>

using namespace std;

int main(void) {
    MyType t(1);
    MyDerivedType dt = (MyDerivedType)t;
    MyType back = (MyType)dt;
    cout << dt.id << endl;
    cout << back.id << endl;
}