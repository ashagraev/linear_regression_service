# linear_regression_service

## 1. Simple linear regression problem

A simple linear regression problem states as follows: given two-dimensional sample points having one independent variable and one target value, build a linear model to minimize the residual sum of squared errors. See [1]   for further details.

![](https://user-images.githubusercontent.com/6789687/93579011-a5d04a80-f9a6-11ea-975c-1f69443bcf0c.png)

This problem is relatively easy in terms of computation costs. However, numerical errors could potentially lead to unstable and improper results. To deal with that problem, we use Welford's method [2] for calculating means and covariations as well as Kahan's summation algorithm [3].

Links:
1. https://en.wikipedia.org/wiki/Simple_linear_regression
2. https://en.wikipedia.org/wiki/Algorithms_for_calculating_variance#Welford's_online_algorithm
3. https://en.wikipedia.org/wiki/Kahan_summation_algorithm
