# go-galton

Uses go concurrency to simulate a [Galton board](https://en.wikipedia.org/wiki/Galton_board).

The marble source is modeled as a bool channel that emits one true value per marble.

Each peg is modeled as a bool channel with two "child" channels, that forwards
values to either its left or right child with 50% probability.

Results are collected into bins from the final row of leaf channels.

Channel closures and WaitGroups are used to gate completion of the simulation run.

```sh
% go run main.go -marbles=2000 -bins=7
 31 ***
198 ************************
456 *******************************************************
660 ********************************************************************************
439 *****************************************************
194 ***********************
 22 **
```