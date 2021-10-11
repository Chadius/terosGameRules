# What is it?
- The main program now takes a script file as an argument.
- The script file has enough information to load and recreate an encounter.
- The script has actions that contain the user, the power, and the targets
- Each action also has a random number seed.
- Each action has a version, so we know when it was made.
- The script will stop if there is an invalid action (the target dies and then tries to do something.)

## Side effects
- We weren't checking to see if the user was alive.

# Why is it so important?
- We'll be able to replay multiple actions.
- We can remove the hardcoded encounter from main.go
- We can manipulate the random seed to recreate specific scenarios

# Now that we can do this...
- We can work on recording scripts.
- We can troubleshoot events with a recording.
- We can replace the random seed with a fixed die roll to generate randomization (instead of a random seed)
- We can verify versions, too.