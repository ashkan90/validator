# About Package

## Usage (Tested only with 'url.Values')
Using validator is easy as much as needed.
Some Usage rules:
```
"name": "required|max:100|min:20|equal:Blabla"
"age": "required|digit|between:0-255"
```

There are three types of token<br> 
<hr>
First one (Separator): ``|`` <br>
Second one (Value Input): ``:`` <br>
Third one (Value Separator): ``-`` <br>
<hr>

```
fmt.Println(validator.Load(map[string]string{
		"name": "required|equal:Emirhan",
		"surname": "max:10",
		"age": "required|digit|between:0-255",
		"password": "required|confirmation:password_confirmation",
		"password_confirmation": "required",
		"today": "required|date",
	}, r.Form).Run())
```

You don't need to initialize bunch of unused structs. There's only simplicity

## About Rules
### Date: 
date rule is just using 'Y-m-D' format currently.
```
"today": "required|date",
```

### Max-Min:
As you can understand, max/min is checking for string's length.
Strings are represented as character-bytes of array in memory like this:
``
['H', 'e', 'l', 'l', 'o', ' ', 'W'...]
``
and we're counting its length then comparing with your rule input.



## Features
* [x] Is array rule ('array')
* [x] Is string rule ('string')
* [x] Is integer and digit rules ('integer'|'digit')
* [x] Is between rule('between:x-y')
* [x] Is confirmed rule ('confirmed:Hello world') // case sensitive
* [x] Is date rule ('date') // uses Y-m-D format
* [] date rule format and localization can be changed by user
* [x] Is file rule ('file')
* [] Is file's max size grater than i decided to ?
* [] Is file's min size lower than i decided to ?
* [] Is field email ?
* [x] Errors messages can be internalization
* [] Agnostic rule checkers
* [x] Any type form fields.
* [] Tests