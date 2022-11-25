#include <memory>
#include <vector>
#include <iostream>

using namespace std;

struct A {
    int a =1;
};

struct B : public A {
    int b = 2;
};

template<typename T, typename enable_if<is_base_of<A, T>::value>::type* = nullptr>
shared_ptr<T> make(const T& t) {
    return make_shared<T>(t);
}

int main(void) {
    vector<shared_ptr<int>> nums;
    nums.push_back(make_shared<int>(1));
    nums.push_back(shared_ptr<int>(new int(3)));
    for (auto n : nums) {
        cout << *n << endl;
    }
    cout << make(B())->b << endl;
}