# Exercise 1

## Question

The simple calculator program doesn't handle one error case: division by zero. Change the function signature for the math operations to return both an `int` and an `error`. In the `div` function, if the divisor is `0`, return `0, errors.New("division by zero")`. In all other cases, return the calculated value and `nil`. Adjust the `main` function to check for this error.

## Solution

