# Exercise 3

## Question

Write a program that builds a `[]Person` with 10,000,000 entries (they could all be the same names and ages). See how long it takes to run. Change the value of `GOGC` and see how that effects the time it takes for the program to complete. Set the environment variable `GODEBUG=gctrace=1` to see when garbage collections happen and see how changing `GOGC` changes the number of garbage collections. What happens if you create the slice with a capacity of 10,000,000?

## Solution
