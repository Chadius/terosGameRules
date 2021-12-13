# Why do this
Logic can be compartmentalized by a certain field.
It's difficult to inject healing amounts (partial vs full vs none) because the logic is all over the place in various modules.

# What will this be
It would be easier to build a class that answers questions based on the type.
For example, make a Healing Bonus class that will multiply Full, Half, or None to healing amounts.
Tests will be used on the logical objects.
Use cases will inject healing objects as needed to work on their tests.
Builders will create logic objects as needed

# What's next
Use case tests will work on interfacing with logic objects rather than focus on how they work.

# Yikes! Circular import! Make more interfaces!
I want the Healing Logic object to know about Power and Squaddie, so I can keep it general.
BUT I want the power to hold the Healing Logic.

  Power depends on Healing Logic
  Healing Logic depends on Power

Create an Interface object for Power:
  Healing Logic uses power.interface
  Power uses healing.interface

I also had to make squaddie.interface. There is promise here to avoid all dependencies.

# Now it's time for Factories
How do I easily copy logic objects?
- I can give a string to a logic factory and that will return an object 

# Logic objects progress
## In progress
Healing multipliers as objects
## Converted
