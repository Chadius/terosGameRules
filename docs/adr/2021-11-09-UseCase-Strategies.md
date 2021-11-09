# Why
Tests require a lot of setup to make sure UseCases produce the expected behavior.
Instead, let's make UseCases that can be abstracted.
This is the Strategy pattern.

# What does this look like
Each existing Use Case will become a Strategy.
Strategy will implement an interface
Tests will use that interface.