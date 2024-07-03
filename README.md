# figuring

Figuring is a geometry library for math on simple 2D shapes and for Bézier curves. It also has some simple polynomial solvers for functions on Bézier curves.

```
one = figuring.ConstantA(1)
linear = figuring.LinearAb(2, 1)
quadratic = figuring.Quadratic(3, 2, 1)
cubic = figuring.CubicAbcd(4, 3, 2, 1)
quartic = figuring.QuarticAbcde(5, 4, 3, 2, 1)
```

For Bézier curves, there is the general `ParamCurve`, in addition to the `Bezier` type for Cubic Bézier curves. Many of the figuring functions were informed by [Pomax’s *A Primer on Bézier Curves*](https://pomax.github.io/bezierinfo/).

A bunch of other 2D shapes have been added to the library to round out the functionality, and make it useful for general geometry.