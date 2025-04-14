cmake_minimum_required(VERSION 3.0)

project(hash-benchmarks)

set(CMAKE_CXX_STANDARD 23)

set(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} /EHsc -march=native")

# = Libraries

# == RollingHashCPP
# Note: made a subdirectory, so use like `#include <rollinghashcpp/rabinkarphash.h>`
include_directories(SYSTEM lib/rollinghashcpp/include)

# == HighwayHash
# TODO: import via CMake properly

# TODO: add_executable
