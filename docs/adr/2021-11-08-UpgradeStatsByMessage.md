# Why
UseCases are directly manipulating (what should be) private variables
This means it's difficult to test.
It also steals responsibility from the Squaddie - the Squaddie should deal with level ups.

# How to solve this problem
Any time a stat is changed, add a Squaddie public function that can alter it instead
Use the public function to change the stat
Delete the direct manipulation