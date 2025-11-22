# Exercise 3

## Question

Define an interface called `Ranker` that has a single method called `Ranking` that returns a slice of trings. Write a function called `RankPrinter` with two parameters, the first of type `Ranker` and the second of type `io.Writer`. Use the `io.WriteString` function to write the values returned by `Ranker` to the `io.Writer`, with a newline separating each result. Call this function from main.

## Solution
