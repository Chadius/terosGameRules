# Why
In converting this to a library, we need to stop using stdout for error pipes.
We also need to stop using files as data input.

# What is it
We should be able to use ANY stream for input or output.

# What can we do now?
The library can output to file, or stdout, or a network stream.

# Caveats that will trigger future change
Assumed that text will always be good enough; what about different languages?