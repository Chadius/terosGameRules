# Why
By separating creation and access from the model objects, I no longer need to make them public.
This protects them from getting set arbitrarily.

# How do I accomplish this?
1. Let the builders use YAML and JSON to make objects (so the API does not rely on the implementation.)
2. Call NewClassName() to create the objects (so other libraries can't build it piece meal and error-prone)
3. Privatize fields and remove YAML/JSON annotations (so they cannot be set publicly)

# Caveats
- Squaddies are heavily intertwined with other systems. I need to separate them first before I can privatize them.
- Many helper objects need to be privatized first, so I can separate them from Squaddie.
