# TODO (heh)

A simple yet specific todo app catered to pretty much my own workflow

For now only local storage is supported - implemented with `nanobox-io/golang-scribble` as an extremely simple flat file json storage. I might get to making it work cross device one day.

# Why / Use case / Value Prop

- No web UI that takes some precious screen real estate
- Doesn't deal in deadlines by design
- Fuzzy search

This todo app is quite dedicated to an engineer that gets trapped in 1500 context switches a day - with alerts bubbling up left and right, and million things to "get back to", "when I have the spare cycles".

Now, instead of claiming "yeah I'll do that later" and forget about it into oblivion, you can punt the item into a weight ranked todo list, and churn down on that todo list when you do have the spare cycle.

# Install

Install from source -

```shell
go install github.com/2785/todo@latest
```

# Use

```shell
# add item
todo add 5 here is an item having weight of 5

> added todo with id 467fa0aa-519f-4770-a4f9-0e38ec0483b4

todo add 1 here is some other low value items

> added todo with id b61cbe45-e133-4191-a88c-0adfa85aaa32

todo add 10 here is an item that has something todo with a rabbit

> added todo with id ba777036-59ee-4954-b357-d7e8d9b9fe49

# list items
todo ls

# by default rank asc by weight - top priority item shows up on bottom
>   W: 1 - added Oct 02 11:07 Sat
    ID: b61cbe45-e133-4191-a88c-0adfa85aaa32
    Desc: here is some other low value items
    ----------
    W: 5 - added Oct 02 11:07 Sat
    ID: 467fa0aa-519f-4770-a4f9-0e38ec0483b4
    Desc: here is an item having weight of 5
    ----------
    W: 10 - added Oct 02 11:08 Sat
    ID: ba777036-59ee-4954-b357-d7e8d9b9fe49
    Desc: here is an item that has something todo with a rabbit

# mark things as done
todo done

>   ? Which TODOs to close?  [Use arrows to move, space to select, <right> to all, <left> to none, type to filter]
    > [ ]  here is an item that has something todo with a rabbit
    [ ]  here is an item having weight of 5
    [ ]  here is some other low value items

# ergonomic fuzzy search - granted the selection screen itself lets you do this as well
todo done rabbit

>   ? Which TODOs to close?  [Use arrows to move, space to select, <right> to all, <left> to none, type to filter]
    > [ ]  here is an item that has something todo with a rabbit
```
