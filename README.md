# shellreminders
Shows reminders in my terminal about things that I have to pay soon, like credit cards ...

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
