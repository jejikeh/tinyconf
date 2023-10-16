# ambient


change project name somewhere in the future.

# NOTEs:

1. Maybe create Tree as based for virtual machine 

the stack based virtual machine looks like this



# aasm

```asm
main:
    psh 1 2 2
```

# future syntax?

## loops

So, in c# we have many different ways to iterate on things. But, the main problem with it what is requires a lot of refactoring when it comes to changes from one type of loop to another. 

```c#
while (something.IsTrue) { ... }

for (var i = 0; i < 12; i++) { ... }

foreach (var x in container) { ... }

```

So, the first we can do, it remove foreach and while loops. The main loop is `for`, so 

```go

for (i := something.is_true) { }

// translates to (var i = 0; i < 12; i++)
for (i := 0; i < 12) { }

// translates to (var i = 0; i < 12 && something.is_true; i += 2)
for (i := 0; i > 12 && something.is_true; i += 2)

// translates to foreach(var in arr_from_0_12) {}
for (i := *[0..12])


for (i := *[container]; i += 2)
```

## arrays