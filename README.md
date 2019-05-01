# shellreminders
Shows reminders in my terminal about things that I have to pay soon, like credit cards ...

## Why am I doing this?
Because I spend a lot of time in the terminal, I can use my phone but I think that seeing this everything I open a terminal
and my ~/.bashrc is read will make me to remember.

```
./shellreminders 
┌────────────────────────────────────────────┐
│ Remaining days for 'Santander Platino' : 6 │
└────────────────────────────────────────────┘
┌─────────────────────────────────────┐
│ Remaining days for 'Promotions' : 1 │
└─────────────────────────────────────┘
```

Configuration (input) file:

```
$ cat ~/.shellreminder/reminders 
Santander Platino;18
Promotions;13;counter

```
