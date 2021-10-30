# Why
Constructing objects is very cumbersome - I need to know how the data shape of each one.
Instead, I should make a builder and give it a flat file.
Consumers can use that instead.

# What does this look like
The Builder will become responsible for consuming YAML and JSON definition files.
The power and squaddie repository will call the builder, passing the YAML/JSON in.
The Builder will figure out what format the file is and then return the object (or nil if it can't.)
The repository is responsible for dealing with failure.

# Future Caveats
Squaddies are tightly coupled with Powers - the interface requires you to add a power rather than adding a reference. Squaddies should only have references.
