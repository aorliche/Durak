.RECIPEPREFIX = >

simple3: base.hpp cards.hpp search.hpp simple3enumerate.cpp
> g++ -std=c++17 -o simple3 simple3enumerate.cpp

simple2: base.hpp cards.hpp search.hpp simple2composition.cpp
> g++ -std=c++17 -o simple2 simple2composition.cpp

simple1: base.hpp simple1testbase.cpp
> g++ -std=c++17 -o simple1 simple1testbase.cpp

begin2:
> g++ -o begin2beats begin2beats.cpp

begin1:
> g++ -o begin1action begin1action.cpp

begin1dbg:
> g++ -g -o begin1action begin1action.cpp
